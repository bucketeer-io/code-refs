package search

import (
	"context"
	"sort"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/bucketeer-io/code-refs/internal/bucketeer"
	"github.com/bucketeer-io/code-refs/internal/helpers"
	"github.com/go-git/go-git/v5/plumbing"
)

const (
	// These are defensive limits intended to prevent corner cases stemming from
	// large repos, false positives, etc. The goal is a) to prevent the program
	// from taking a very long time to run and b) to prevent the program from
	// PUTing a massive json payload. These limits will likely be tweaked over
	// time. The LaunchDarkly backend will also apply limits.
	maxFileCount     = 10000 // Maximum number of files containing code references
	maxHunkCount     = 25000 // Maximum number of total code references
	maxLineCharCount = 500   // Maximum number of characters per line
)

// Truncate lines to prevent sending over massive hunks, e.g. a minified file.
// NOTE: We may end up truncating a valid flag key reference. We accept this risk
// and will handle hunks missing flag key references on the frontend.
func truncateLine(line string, maxCharCount int) string {
	if utf8.RuneCountInString(line) <= maxCharCount {
		return line
	}
	// convert to rune slice so that we don't truncate multibyte unicode characters
	runes := []rune(line)
	return string(runes[0:maxCharCount]) + "…"
}

type file struct {
	path    string
	lines   []string
	fileExt string
}

// hunkForLine returns a matching code reference for a given flag key on a line
func (f file) hunkForLine(flagKey string, lineNum int, matcher Matcher) *bucketeer.HunkRep {
	line := f.lines[lineNum]
	ctxLines := matcher.ctxLines

	aliasMatches := matcher.FindAliases(line, flagKey)
	if len(aliasMatches) == 0 && !matcher.MatchElement(line, flagKey) {
		return nil
	}

	startingLineNum := lineNum
	var hunkLines []string
	if ctxLines >= 0 {
		startingLineNum -= ctxLines
		if startingLineNum < 0 {
			startingLineNum = 0
		}
		endingLineNum := lineNum + ctxLines + 1
		if endingLineNum >= len(f.lines) {
			hunkLines = f.lines[startingLineNum:]
		} else {
			hunkLines = f.lines[startingLineNum:endingLineNum]
		}
	}

	for i, line := range hunkLines {
		hunkLines[i] = truncateLine(line, maxLineCharCount)
	}

	lines := strings.Join(hunkLines, "\n")
	contentHash := getContentHash(lines)

	ret := bucketeer.HunkRep{
		FlagKey:            flagKey,
		StartingLineNumber: startingLineNum + 1,
		Lines:              lines,
		Aliases:            aliasMatches,
		ContentHash:        contentHash,
		FileExt:            f.fileExt,
	}
	return &ret
}

// aggregateHunksForFlag finds all references in a file, and combines matches if their context lines overlap
func (f file) aggregateHunksForFlag(flagKey string, matcher Matcher, lineNumbers []int) []bucketeer.HunkRep {
	var hunksForFlag []bucketeer.HunkRep
	for _, lineNumber := range lineNumbers {
		match := f.hunkForLine(flagKey, lineNumber, matcher)
		if match != nil {
			lastHunkIdx := len(hunksForFlag) - 1
			// If the previous hunk overlaps or is adjacent to the current hunk, merge them together
			if lastHunkIdx >= 0 && hunksForFlag[lastHunkIdx].Overlap(*match) >= 0 {
				hunksForFlag = append(hunksForFlag[:lastHunkIdx], mergeHunks(hunksForFlag[lastHunkIdx], *match)...)
			} else {
				hunksForFlag = append(hunksForFlag, *match)
			}
		}
	}
	return hunksForFlag
}

func (f file) toHunks(matcher Matcher) *bucketeer.ReferenceHunksRep {
	hunks := make([]bucketeer.HunkRep, 0)
	elementSearch := matcher.Element
	if elementSearch.Dir != "" {
		matchDir := strings.HasPrefix(f.path, elementSearch.Dir)
		if !matchDir {
			return nil
		}
	}
	lineNumbersByElement := f.findMatchingLineNumbersByElement(elementSearch)
	for element, lineNumbers := range lineNumbersByElement {
		hunks = append(hunks, f.aggregateHunksForFlag(element, matcher, lineNumbers)...)
	}
	if len(hunks) == 0 {
		return nil
	}
	// Set the file extension for all hunks
	for i := range hunks {
		hunks[i].FileExt = f.fileExt
	}
	return &bucketeer.ReferenceHunksRep{Path: f.path, Hunks: hunks}
}

func (f file) findMatchingLineNumbersByElement(matcher ElementMatcher) map[string][]int {
	lineNumbersByElement := make(map[string][]int)
	for lineNum, line := range f.lines {
		for _, element := range matcher.FindMatches(line) {
			lineNumbersByElement[element] = append(lineNumbersByElement[element], lineNum)
		}
	}
	return lineNumbersByElement
}

// mergeHunks combines the lines and aliases of two hunks together for a given file
// if the hunks do not overlap, returns each hunk separately
// assumes the startingLineNumber of a is less than b and there is some overlap between the two
func mergeHunks(a, b bucketeer.HunkRep) []bucketeer.HunkRep {
	if a.StartingLineNumber > b.StartingLineNumber {
		a, b = b, a
	}

	aLines := strings.Split(a.Lines, "\n")
	bLines := strings.Split(b.Lines, "\n")

	overlap := a.Overlap(b)
	// no overlap
	if overlap < 0 || len(a.Lines) == 0 && len(b.Lines) == 0 {
		return []bucketeer.HunkRep{a, b}
	} else if overlap >= len(bLines) {
		// subset hunk
		return []bucketeer.HunkRep{a}
	}

	combinedLines := append(aLines, bLines[overlap:]...)
	lines := strings.Join(combinedLines, "\n")
	contentHash := getContentHash(lines)

	return []bucketeer.HunkRep{
		{
			StartingLineNumber: a.StartingLineNumber,
			Lines:              lines,
			FlagKey:            a.FlagKey,
			Aliases:            helpers.Dedupe(append(a.Aliases, b.Aliases...)),
			ContentHash:        contentHash,
		},
	}
}

// processFiles starts goroutines to process files individually. When all files have completed processing, the references channel is closed to signal completion.
func processFiles(ctx context.Context, files <-chan file, references chan<- bucketeer.ReferenceHunksRep, matcher Matcher) {
	defer close(references)
	w := sync.WaitGroup{}
	for f := range files {
		if ctx.Err() != nil {
			// context cancelled, stop processing files, but let the waitgroup finish organically
			continue
		}
		w.Add(1)
		go func(f file) {
			reference := f.toHunks(matcher)
			if reference != nil {
				references <- *reference
			}
			w.Done()
		}(f)
	}
	w.Wait()
}

func SearchForRefs(directory string, matcher Matcher) ([]bucketeer.ReferenceHunksRep, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	files := make(chan file)
	references := make(chan bucketeer.ReferenceHunksRep)
	// Start workers to process files asynchronously as they are written to the files channel
	go processFiles(ctx, files, references, matcher)

	err := readFiles(ctx, files, directory)
	if err != nil {
		return nil, err
	}

	ret := make([]bucketeer.ReferenceHunksRep, 0, len(references))

	defer sort.SliceStable(ret, func(i, j int) bool {
		return ret[i].Path < ret[j].Path
	})

	totalHunks := 0
	for reference := range references {
		ret = append(ret, reference)

		// Reached maximum number of files with code references
		if len(ret) >= maxFileCount {
			return ret, nil
		}
		totalHunks += len(reference.Hunks)
		// Reached maximum number of hunks across all files
		if totalHunks > maxHunkCount {
			return ret, nil
		}
	}
	return ret, nil
}

func getContentHash(lines string) string {
	return plumbing.ComputeHash(plumbing.BlobObject, []byte(lines)).String()
}

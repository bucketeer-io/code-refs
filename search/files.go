package search

import (
	"bufio"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/monochromegane/go-gitignore"

	"github.com/bucketeer-io/code-refs/internal/validation"
)

type ignore struct {
	path    string
	ignores []gitignore.IgnoreMatcher
}

func newIgnore(path string, ignoreFiles []string) ignore {
	ignores := make([]gitignore.IgnoreMatcher, 0, len(ignoreFiles))
	for _, ignoreFile := range ignoreFiles {
		i, err := gitignore.NewGitIgnore(filepath.Join(path, ignoreFile))
		if err != nil {
			continue
		}
		ignores = append(ignores, i)
	}
	return ignore{path: path, ignores: ignores}
}

func (m ignore) Match(path string, isDir bool) bool {
	for _, i := range m.ignores {
		if i.Match(path, isDir) {
			return true
		}
	}

	return false
}

func readFileLines(path string) ([]string, error) {
	if !validation.FileExists(path) {
		return nil, errors.New("file does not exist")
	}

	/* #nosec */
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

// isText reports whether the first kilobyte of the file (its lines joined by
// newlines) looks like correct UTF-8 without control characters; that is, if
// it is likely human-readable text. Taking the lines directly avoids copying
// the whole file just to inspect its prefix. Adapted from
// golang.org/x/tools/godoc/util (BSD-3-Clause), whose module is deprecated
// and frozen.
func isText(lines []string) bool {
	const maxCheck = 1024 // at least utf8.UTFMax
	s := make([]byte, 0, maxCheck)
	for i, line := range lines {
		if i > 0 {
			s = append(s, '\n')
		}
		if len(line) > maxCheck-len(s) {
			line = line[:maxCheck-len(s)]
		}
		s = append(s, line...)
		if len(s) >= maxCheck {
			break
		}
	}
	for i, c := range string(s) {
		if i+utf8.UTFMax > len(s) {
			// last char may be incomplete - ignore
			break
		}
		if c == 0xFFFD || c < ' ' && c != '\n' && c != '\t' && c != '\f' {
			// decoding error or control character - not a text file
			return false
		}
	}
	return true
}

func readFiles(ctx context.Context, files chan<- file, workspace, subdirectory string) error {
	defer close(files)
	ignoreFiles := []string{".gitignore", ".ignore", ".ldignore"}
	allIgnores := newIgnore(workspace, ignoreFiles)
	workspace = filepath.ToSlash(workspace)

	readFile := func(path string, info os.FileInfo, err error) error {
		if err != nil || ctx.Err() != nil {
			// global context cancelled, don't read any more files
			return nil
		}

		isDir := info.IsDir()
		path = filepath.ToSlash(path)

		// Skip directories, hidden files, and ignored files
		if allIgnores.Match(path, isDir) {
			if isDir {
				return filepath.SkipDir
			}
			return nil
		} else if strings.HasPrefix(info.Name(), ".") {
			if isDir {
				// don't skip github dir
				if strings.HasPrefix(info.Name(), ".github") {
					return nil
				}
				return filepath.SkipDir
			}
			return nil
		} else if !info.Mode().IsRegular() {
			return nil
		}

		lines, err := readFileLines(path)
		if err != nil {
			return err
		}

		// only read text files
		if !isText(lines) {
			return nil
		}

		relativePath := resolvePath(path, workspace, subdirectory)
		files <- file{
			path:    relativePath,
			lines:   lines,
			fileExt: strings.TrimPrefix(filepath.Ext(relativePath), "."),
		}
		return nil
	}

	return filepath.Walk(workspace, readFile)
}

// resolvePath makes path relative to the repo root rather than the searched
// workspace, so paths stay correct when a subdirectory is configured (the
// workspace is then <root>/<subdirectory>).
func resolvePath(path, workspace, subdirectory string) string {
	dir := workspace
	if subdirectory != "" {
		// Normalize the same way filepath.Join normalized subdirectory when
		// building workspace, so a trailing slash or "./" prefix doesn't
		// prevent the suffix match below (which would silently drop the
		// subdirectory from the resolved path).
		cleanSubdirectory := filepath.ToSlash(filepath.Clean(subdirectory))
		dir = strings.TrimSuffix(workspace, "/"+cleanSubdirectory)
	}

	return strings.TrimPrefix(path, dir+"/")
}

package coderefs

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bucketeer-io/code-refs/internal/bucketeer"
	"github.com/bucketeer-io/code-refs/internal/git"
	"github.com/bucketeer-io/code-refs/internal/helpers"
	"github.com/bucketeer-io/code-refs/internal/log"
	"github.com/bucketeer-io/code-refs/internal/validation"
	"github.com/bucketeer-io/code-refs/options"
	"github.com/bucketeer-io/code-refs/search"
)

func Run(opts options.Options, output bool) {
	absPath, err := validation.NormalizeAndValidatePath(opts.Dir)
	if err != nil {
		log.Error.Fatalf("could not validate directory option: %s", err)
	}

	log.Info.Printf("absolute directory path: %s", absPath)

	totalEnvs := len(opts.ApiKey)
	log.Info.Printf("Processing %d api key(s)", totalEnvs)

	// Process each API key
	for i, apiKey := range opts.ApiKey {
		log.Info.Printf("Processing api key %d/%d", i+1, totalEnvs)

		// Create a copy of options with current API key
		currentOpts := opts
		currentOpts.ApiKey = []string{apiKey}

		bucketeerApi := initializeAPI(currentOpts)
		branchName, revision := setupGitInfo(currentOpts, absPath)
		repoType := determineRepoType(currentOpts.RepoType)

		matcher, refs := search.Scan(currentOpts, absPath)
		if output {
			generateOutput(currentOpts, matcher, refs)
		}

		if !currentOpts.DryRun {
			processCodeReferences(currentOpts, bucketeerApi, refs, branchName, revision, repoType)
		}
	}
}

//nolint:ireturn // This function returns an interface for testing/mocking purposes
func initializeAPI(opts options.Options) bucketeer.ApiClient {
	return bucketeer.InitApiClient(bucketeer.ApiOptions{
		ApiKey:    opts.ApiKey[0],
		BaseUri:   opts.BaseUri,
		UserAgent: helpers.GetUserAgent(opts.UserAgent),
	})
}

func setupGitInfo(opts options.Options, absPath string) (branchName, revision string) {
	branchName = opts.Branch
	revision = opts.Revision
	if revision == "" {
		gitClient, err := git.NewClient(absPath, branchName, opts.AllowTags)
		if err != nil {
			log.Error.Fatalf("%s", err)
		}
		branchName = gitClient.GitBranch
		revision = gitClient.GitSha
	}
	return branchName, revision
}

func determineRepoType(repoType string) string {
	repoType = strings.ToUpper(repoType)
	if repoType != "GITHUB" && repoType != "GITLAB" && repoType != "BITBUCKET" {
		repoType = "CUSTOM"
	}
	return repoType
}

func processCodeReferences(
	opts options.Options,
	bucketeerApi bucketeer.ApiClient,
	refs []bucketeer.ReferenceHunksRep,
	branchName, revision, repoType string,
) {
	existingRefs := fetchExistingReferences(opts, bucketeerApi, refs)
	processNewReferences(opts, bucketeerApi, refs, existingRefs, branchName, revision, repoType)
	deleteStaleReferences(opts, bucketeerApi, existingRefs)
}

func fetchExistingReferences(
	opts options.Options,
	bucketeerApi bucketeer.ApiClient,
	refs []bucketeer.ReferenceHunksRep,
) map[string]bucketeer.CodeReference {
	existingRefs := make(map[string]bucketeer.CodeReference)
	flagCounts := aggregateFeatureFlags(refs)

	for flag := range flagCounts {
		codeRefs, _, _, err := bucketeerApi.ListCodeReferences(context.Background(), opts, flag, bucketeer.DefaultPageSize)
		if err != nil {
			helpers.FatalServiceError(
				fmt.Errorf("error getting existing code references from Bucketeer for flag %s: %w", flag, err),
				opts.IgnoreServiceErrors,
			)
		}

		for _, ref := range codeRefs {
			existingRefs[ref.ContentHash] = ref
		}
	}
	return existingRefs
}

func processNewReferences(
	opts options.Options,
	bucketeerApi bucketeer.ApiClient,
	refs []bucketeer.ReferenceHunksRep,
	existingRefs map[string]bucketeer.CodeReference,
	branchName, revision, repoType string,
) {
	for _, ref := range refs {
		for _, hunk := range ref.Hunks {
			codeRef := createCodeReference(opts, hunk, ref, branchName, revision, repoType)
			updateOrCreateReference(opts, bucketeerApi, codeRef, existingRefs)
		}
	}
}

func createCodeReference(opts options.Options,
	hunk bucketeer.HunkRep,
	ref bucketeer.ReferenceHunksRep,
	branchName, revision, repoType string,
) bucketeer.CodeReference {
	return bucketeer.CodeReference{
		FeatureID:        hunk.FlagKey,
		FilePath:         ref.Path,
		FileExtension:    hunk.FileExt,
		LineNumber:       hunk.StartingLineNumber,
		CodeSnippet:      hunk.Lines,
		ContentHash:      hunk.ContentHash,
		Aliases:          hunk.Aliases,
		RepositoryName:   opts.RepoName,
		RepositoryOwner:  opts.RepoOwner,
		RepositoryType:   repoType,
		RepositoryBranch: strings.TrimPrefix(branchName, "refs/heads/"),
		CommitHash:       revision,
	}
}

func updateOrCreateReference(opts options.Options,
	bucketeerApi bucketeer.ApiClient,
	codeRef bucketeer.CodeReference,
	existingRefs map[string]bucketeer.CodeReference,
) {
	if existing, exists := existingRefs[codeRef.ContentHash]; exists {
		log.Info.Printf("updating code reference in Bucketeer: id: %s, content hash: %s", existing.ID, codeRef.ContentHash)
		err := bucketeerApi.UpdateCodeReference(context.Background(), opts, existing.ID, codeRef)
		if err != nil {
			helpers.FatalServiceError(fmt.Errorf("error updating code reference in Bucketeer: %w", err), opts.IgnoreServiceErrors)
		}
		delete(existingRefs, codeRef.ContentHash)
	} else {
		err := bucketeerApi.CreateCodeReference(context.Background(), opts, codeRef)
		if err != nil {
			helpers.FatalServiceError(fmt.Errorf("error sending code reference to Bucketeer: %w", err), opts.IgnoreServiceErrors)
		}
	}
}

func deleteStaleReferences(opts options.Options, bucketeerApi bucketeer.ApiClient, existingRefs map[string]bucketeer.CodeReference) {
	for _, ref := range existingRefs {
		if ref.RepositoryOwner == opts.RepoOwner && ref.RepositoryName == opts.RepoName {
			err := bucketeerApi.DeleteCodeReference(context.Background(), opts, ref.ID)
			if err != nil {
				helpers.FatalServiceError(fmt.Errorf("error deleting code reference from Bucketeer: %w", err), opts.IgnoreServiceErrors)
			}
			log.Info.Printf("deleted code reference from Bucketeer: %+v", ref)
		}
	}
}

func generateOutput(opts options.Options, matcher search.Matcher, refs []bucketeer.ReferenceHunksRep) {
	outDir := opts.OutDir
	if outDir != "" {
		outPath, err := writeToCSV(outDir, opts.RepoName, opts.Revision, refs)
		if err != nil {
			log.Error.Fatalf("error writing code references to csv: %s", err)
		}
		log.Info.Printf("wrote code references to %s", outPath)
	}

	if opts.Debug {
		printReferenceCountTable(refs)
	}

	if opts.DryRun {
		totalHunks := 0
		for _, ref := range refs {
			totalHunks += len(ref.Hunks)
		}
		log.Info.Printf(
			"dry run found %d code references across %d flags and %d files",
			totalHunks,
			len(matcher.Element.Elements),
			len(refs),
		)
		return
	}

	log.Info.Printf(
		"sending %d code references across %d flags and %d files to Bucketeer",
		getTotalHunkCount(refs),
		len(matcher.Element.Elements),
		len(refs),
	)
}

func getTotalHunkCount(refs []bucketeer.ReferenceHunksRep) int {
	total := 0
	for _, ref := range refs {
		total += len(ref.Hunks)
	}
	return total
}

func writeToCSV(outDir, repoName, revision string, refs []bucketeer.ReferenceHunksRep) (string, error) {
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("code-references-%s-%s-%d.csv", repoName, revision, timestamp)
	outPath := filepath.Join(outDir, filename)

	file, err := os.Create(outPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	err = writer.Write([]string{
		"Flag Key",
		"File Path",
		"File Extension",
		"Line Number",
		"Code Snippet",
		"Content Hash",
		"Aliases",
	})
	if err != nil {
		return "", err
	}

	// Write data
	for _, ref := range refs {
		for _, hunk := range ref.Hunks {
			err := writer.Write([]string{
				hunk.FlagKey,
				ref.Path,
				hunk.FileExt,
				strconv.Itoa(hunk.StartingLineNumber),
				hunk.Lines,
				hunk.ContentHash,
				strings.Join(hunk.Aliases, "|"),
			})
			if err != nil {
				return "", err
			}
		}
	}

	return outPath, nil
}

// aggregateFeatureFlags returns a map of feature flags and their counts from the references
func aggregateFeatureFlags(refs []bucketeer.ReferenceHunksRep) map[string]int {
	flagCounts := make(map[string]int)
	for _, ref := range refs {
		for _, hunk := range ref.Hunks {
			flagCounts[hunk.FlagKey]++
		}
	}
	return flagCounts
}

func printReferenceCountTable(refs []bucketeer.ReferenceHunksRep) {
	flagCounts := aggregateFeatureFlags(refs)
	log.Info.Printf("Flag Reference Counts:")
	for flag, count := range flagCounts {
		log.Info.Printf("  %s: %d", flag, count)
	}
}

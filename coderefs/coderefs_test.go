package coderefs

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bucketeer-io/code-refs/internal/bucketeer"
	"github.com/bucketeer-io/code-refs/internal/log"
	"github.com/bucketeer-io/code-refs/options"
)

func init() {
	log.Init(true)
}
func TestMain(m *testing.M) {
	log.Init(true)
	os.Exit(m.Run())
}

type fakeApiClient struct {
	existing []bucketeer.CodeReference
	created  []bucketeer.CodeReference
	updated  map[string]bucketeer.CodeReference
	deleted  []string
}

func (f *fakeApiClient) GetFlagKeyList(_ context.Context, _ options.Options) ([]string, error) {
	return nil, nil
}

func (f *fakeApiClient) CreateCodeReference(_ context.Context, _ options.Options, ref bucketeer.CodeReference) error {
	f.created = append(f.created, ref)
	return nil
}

func (f *fakeApiClient) UpdateCodeReference(_ context.Context, _ options.Options, id string, ref bucketeer.CodeReference) error {
	f.updated[id] = ref
	return nil
}

func (f *fakeApiClient) DeleteCodeReference(_ context.Context, _ options.Options, id string) error {
	f.deleted = append(f.deleted, id)
	return nil
}

func (f *fakeApiClient) ListCodeReferences(_ context.Context, _ options.Options, _ string, _ int64) ([]bucketeer.CodeReference, string, string, error) {
	return f.existing, "", "", nil
}

func Test_processCodeReferences_identicalContentHashes(t *testing.T) {
	// two refs for the same flag whose snippets are textually identical (e.g.
	// different secrets both redacted to the same text) share a ContentHash;
	// each must still reconcile with its own scanned hunk instead of one
	// overwriting the other and the second being recreated
	api := &fakeApiClient{
		existing: []bucketeer.CodeReference{
			{ID: "ref-a", FeatureID: "my-flag", FilePath: "a.go", ContentHash: "same-hash", RepositoryOwner: "owner", RepositoryName: "repo"},
			{ID: "ref-b", FeatureID: "my-flag", FilePath: "b.go", ContentHash: "same-hash", RepositoryOwner: "owner", RepositoryName: "repo"},
		},
		updated: map[string]bucketeer.CodeReference{},
	}
	opts := options.Options{RepoOwner: "owner", RepoName: "repo"}
	refs := []bucketeer.ReferenceHunksRep{
		{Path: "a.go", Hunks: []bucketeer.HunkRep{{FlagKey: "my-flag", ContentHash: "same-hash", Lines: `key = "[REDACTED]"`}}},
		{Path: "b.go", Hunks: []bucketeer.HunkRep{{FlagKey: "my-flag", ContentHash: "same-hash", Lines: `key = "[REDACTED]"`}}},
	}

	processCodeReferences(opts, api, refs, "main", "abc123", "GITHUB")

	require.Empty(t, api.created)
	require.Empty(t, api.deleted)
	require.Len(t, api.updated, 2)
	require.Equal(t, "a.go", api.updated["ref-a"].FilePath)
	require.Equal(t, "b.go", api.updated["ref-b"].FilePath)
}

func Test_processCodeReferences_duplicateReferenceKeys(t *testing.T) {
	// references that share even the composite key (identical snippets of one
	// flag in one file, or every hunk when contextLines < 0 hashes the empty
	// string) must reconcile one-to-one instead of the duplicate row being
	// dropped from the map and accumulating on the server run after run
	api := &fakeApiClient{
		existing: []bucketeer.CodeReference{
			{ID: "ref-a", FeatureID: "my-flag", FilePath: "a.go", ContentHash: "h", RepositoryOwner: "owner", RepositoryName: "repo"},
			{ID: "ref-b", FeatureID: "my-flag", FilePath: "a.go", ContentHash: "h", RepositoryOwner: "owner", RepositoryName: "repo"},
		},
		updated: map[string]bucketeer.CodeReference{},
	}
	opts := options.Options{RepoOwner: "owner", RepoName: "repo"}
	refs := []bucketeer.ReferenceHunksRep{
		{Path: "a.go", Hunks: []bucketeer.HunkRep{
			{FlagKey: "my-flag", ContentHash: "h"},
			{FlagKey: "my-flag", ContentHash: "h"},
		}},
	}

	processCodeReferences(opts, api, refs, "main", "abc123", "GITHUB")

	require.Empty(t, api.created)
	require.Empty(t, api.deleted)
	require.Len(t, api.updated, 2)
	require.Contains(t, api.updated, "ref-a")
	require.Contains(t, api.updated, "ref-b")
}

func Test_processCodeReferences_backfillsMissingFeatureID(t *testing.T) {
	// the list endpoint is queried by featureId; if the server omits the
	// field in response items, the queried flag must be backfilled so the
	// reconciliation key still matches instead of every ref being recreated
	// and the old rows deleted
	api := &fakeApiClient{
		existing: []bucketeer.CodeReference{
			{ID: "ref-a", FilePath: "a.go", ContentHash: "h", RepositoryOwner: "owner", RepositoryName: "repo"},
		},
		updated: map[string]bucketeer.CodeReference{},
	}
	opts := options.Options{RepoOwner: "owner", RepoName: "repo"}
	refs := []bucketeer.ReferenceHunksRep{
		{Path: "a.go", Hunks: []bucketeer.HunkRep{{FlagKey: "my-flag", ContentHash: "h"}}},
	}

	processCodeReferences(opts, api, refs, "main", "abc123", "GITHUB")

	require.Empty(t, api.created)
	require.Empty(t, api.deleted)
	require.Len(t, api.updated, 1)
	require.Contains(t, api.updated, "ref-a")
}

// func Test_calculateStaleBranches(t *testing.T) {
// 	specs := []struct {
// 		name           string
// 		branches       []string
// 		remoteBranches []string
// 		expected       []string
// 	}{
// 		{
// 			name:           "stale branch",
// 			branches:       []string{"main", "another-branch"},
// 			remoteBranches: []string{"main"},
// 			expected:       []string{"another-branch"},
// 		},
// 		{
// 			name:           "no stale branches",
// 			branches:       []string{"main"},
// 			remoteBranches: []string{"main"},
// 			expected:       []string{},
// 		},
// 	}

// 	for _, tt := range specs {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// transform test args into the format expected by calculateStaleBranches
// 			branchReps := make([]ld.BranchRep, 0, len(tt.branches))
// 			for _, b := range tt.branches {
// 				branchReps = append(branchReps, ld.BranchRep{Name: b})
// 			}
// 			remoteBranchMap := map[string]bool{}
// 			for _, b := range tt.remoteBranches {
// 				remoteBranchMap[b] = true
// 			}

// 			assert.ElementsMatch(t, tt.expected, calculateStaleBranches(branchReps, remoteBranchMap))
// 		})
// 	}
// }

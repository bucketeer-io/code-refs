package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bucketeer-io/code-refs/internal/log"
	o "github.com/bucketeer-io/code-refs/options"
)

func TestMain(m *testing.M) {
	log.Init(true)
	os.Exit(m.Run())
}

func TestMergeBitbucketOptions_withCliRepoName(t *testing.T) {
	os.Setenv("BITBUCKET_GIT_HTTP_ORIGIN", "https://bitbucket.com/yus")
	os.Setenv("BITBUCKET_BUILD_NUMBER", "100")
	var options o.Options = o.Options{
		ApiKey:   []string{"deaf-beef"},
		RepoName: "myapp-react",
	}

	result, _ := mergeBitbucketOptions(options)

	assert.Equal(t, "myapp-react", result.RepoName)
	assert.Equal(t, "bitbucket", result.RepoType)
	assert.Equal(t, "https://bitbucket.com/yus", result.RepoUrl)
	assert.Equal(t, 100, result.UpdateSequenceId)
}

func TestMergeBitbucketOptions_withBitbucketRepoName(t *testing.T) {
	os.Setenv("BITBUCKET_GIT_HTTP_ORIGIN", "https://bitbucket.com/yus")
	os.Setenv("BITBUCKET_BUILD_NUMBER", "200")
	os.Setenv("BITBUCKET_REPO_SLUG", "myapp-vue")
	var options o.Options = o.Options{
		ApiKey: []string{"deaf-beef"},
	}

	result, _ := mergeBitbucketOptions(options)

	assert.Equal(t, "myapp-vue", result.RepoName)
	assert.Equal(t, "bitbucket", result.RepoType)
	assert.Equal(t, "https://bitbucket.com/yus", result.RepoUrl)
	assert.Equal(t, 200, result.UpdateSequenceId)
}

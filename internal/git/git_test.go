package git

import (
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/stretchr/testify/require"
)

func Test_remoteSupportsHTTPAuth(t *testing.T) {
	newRemote := func(url string) *git.Remote {
		return git.NewRemote(nil, &config.RemoteConfig{Name: "origin", URLs: []string{url}})
	}

	t.Run("http remote supports auth", func(t *testing.T) {
		require.True(t, remoteSupportsHTTPAuth(newRemote("http://github.com/example/repo.git")))
	})

	t.Run("https remote supports auth", func(t *testing.T) {
		require.True(t, remoteSupportsHTTPAuth(newRemote("https://github.com/example/repo.git")))
	})

	t.Run("ssh remote does not support http auth", func(t *testing.T) {
		// go-git's ssh transport type-asserts Auth to ssh.AuthMethod and
		// rejects anything else (transport.ErrInvalidAuthMethod), so handing
		// it an http.BasicAuth would break branch listing outright.
		require.False(t, remoteSupportsHTTPAuth(newRemote("git@github.com:example/repo.git")))
	})

	t.Run("scp-like ssh remote does not support http auth", func(t *testing.T) {
		require.False(t, remoteSupportsHTTPAuth(newRemote("ssh://git@github.com/example/repo.git")))
	})

	t.Run("remote with no configured URLs does not support auth", func(t *testing.T) {
		require.False(t, remoteSupportsHTTPAuth(git.NewRemote(nil, &config.RemoteConfig{Name: "origin", URLs: []string{}})))
	})
}

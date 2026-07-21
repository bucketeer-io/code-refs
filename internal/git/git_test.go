package git

import (
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/stretchr/testify/require"
)

func Test_remoteAuth(t *testing.T) {
	basicAuth := &http.BasicAuth{Username: "x-access-token", Password: "token"}

	newRemote := func(url string) *git.Remote {
		return git.NewRemote(nil, &config.RemoteConfig{Name: "origin", URLs: []string{url}})
	}

	t.Run("nil auth stays nil regardless of remote scheme", func(t *testing.T) {
		require.Nil(t, remoteAuth(newRemote("https://github.com/example/repo.git"), nil))
	})

	t.Run("http remote gets the auth", func(t *testing.T) {
		require.Equal(t, basicAuth, remoteAuth(newRemote("http://github.com/example/repo.git"), basicAuth))
	})

	t.Run("https remote gets the auth", func(t *testing.T) {
		require.Equal(t, basicAuth, remoteAuth(newRemote("https://github.com/example/repo.git"), basicAuth))
	})

	t.Run("ssh remote does not get http auth", func(t *testing.T) {
		// go-git's ssh transport type-asserts Auth to ssh.AuthMethod and
		// rejects anything else (transport.ErrInvalidAuthMethod), so handing
		// it an http.BasicAuth would break branch listing outright.
		require.Nil(t, remoteAuth(newRemote("git@github.com:example/repo.git"), basicAuth))
	})

	t.Run("scp-like ssh remote does not get http auth", func(t *testing.T) {
		require.Nil(t, remoteAuth(newRemote("ssh://git@github.com/example/repo.git"), basicAuth))
	})

	t.Run("remote with no configured URLs does not get auth", func(t *testing.T) {
		require.Nil(t, remoteAuth(git.NewRemote(nil, &config.RemoteConfig{Name: "origin", URLs: []string{}}), basicAuth))
	})
}

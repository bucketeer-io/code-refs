package search

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_redactSecrets(t *testing.T) {
	tests := []struct {
		name string
		line string
		want string
	}{
		{
			name: "line without secrets is unchanged",
			line: `if client.BoolVariation("someFlag", user, false) {`,
			want: `if client.BoolVariation("someFlag", user, false) {`,
		},
		{
			// caught by the generic keyword rule ("key" in the variable name);
			// a bare access key ID is only caught by the AWS rule when the
			// secret access key is nearby, see Test_redactSecrets_awsKeyPair
			name: "aws access key id assigned to key-named variable",
			line: `awsKey := "AKIAZ9X24KQ7NW3JT6BP"`,
			want: `awsKey := "[REDACTED]"`,
		},
		{
			name: "github personal access token",
			line: `client := github.NewClient("ghp_Wq7xK2mV9tRz3bJ8pNn4Y6dL1gF5hC0Tka9s")`,
			want: `client := github.NewClient("[REDACTED]")`,
		},
		{
			name: "github fine-grained token",
			line: `token := "github_pat_11ABCDEFG0abcdefghijkl_mnopqrstuvwxyz"`,
			want: `token := "[REDACTED]"`,
		},
		{
			name: "gitlab personal access token",
			line: `curl --header "PRIVATE-TOKEN: glpat-aBcDeFgHiJkLmNoPqRsT"`,
			want: `curl --header "PRIVATE-TOKEN: [REDACTED]"`,
		},
		{
			// assembled at runtime so secret scanners don't flag the fake token
			name: "slack token",
			line: `slack.New("xox` + `b-593716823045-4fRk2LpQv7WmZx9TbYd3")`,
			want: `slack.New("[REDACTED]")`,
		},
		{
			name: "stripe live key",
			line: `stripe.Key = "sk_live` + `_aBcDeFgHiJkLmNoPqRsTuVwX"`,
			want: `stripe.Key = "[REDACTED]"`,
		},
		{
			name: "google api key",
			line: `maps.NewClient(maps.WithAPIKey("AIzaSyA1bC2dE3fG4hI5jK6lM7nO8pQ9rS0tU1v"))`,
			want: `maps.NewClient(maps.WithAPIKey("[REDACTED]"))`,
		},
		{
			name: "jwt",
			line: `Authorization: eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U`,
			want: `Authorization: [REDACTED]`,
		},
		{
			name: "bearer authorization header",
			line: `req.Header.Set("Authorization", "Bearer AbCdEfGhIjKlMnOpQrStUvWxYz")`,
			want: `req.Header.Set("Authorization", "Bearer [REDACTED]")`,
		},
		{
			name: "private key header",
			line: `-----BEGIN RSA PRIVATE KEY-----`,
			want: `[REDACTED]`,
		},
		{
			name: "generic api key assignment",
			line: `api_key = "0123456789abcdef"`,
			want: `api_key = "[REDACTED]"`,
		},
		{
			name: "generic password assignment with single quotes",
			line: `password: 'hunter2hunter2'`,
			want: `password: '[REDACTED]'`,
		},
		{
			name: "generic secret assignment with prefixed name",
			line: `const dbSecretValue = "super-secret-value"`,
			want: `const dbSecretValue = "[REDACTED]"`,
		},
		{
			name: "short quoted values are not redacted",
			line: `token = "abc"`,
			want: `token = "abc"`,
		},
		{
			name: "flag key named like a secret keyword is not redacted when unquoted",
			line: `enableTokenRefresh := true`,
			want: `enableTokenRefresh := true`,
		},
		{
			name: "redaction is idempotent",
			line: `api_key = "[REDACTED]"`,
			want: `api_key = "[REDACTED]"`,
		},
	}
	r, err := newRedactor(nil, nil)
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, []string{tt.want}, r.redactHunk([]string{tt.line}))
		})
	}
}

func Test_redactSecrets_awsKeyPair(t *testing.T) {
	r, err := newRedactor(nil, nil)
	require.NoError(t, err)

	// the aws-access-token rule is composite: the key ID is only reported
	// when the paired secret access key appears within 5 lines. The fake
	// credentials are assembled at runtime so secret scanners don't flag them.
	awsKeyID := "AKIA" + "W3JT6BPZ2XKMQ7NC"
	awsSecretKey := "q7PkzOxMjTn2" + "FvBw5RcY8dHl3sGe1uAiK9pXm4tE"
	got := r.redactHunk([]string{
		`[default]`,
		"aws_access_key_id = " + awsKeyID,
		"aws_secret_access_key = " + awsSecretKey,
	})
	require.Equal(t, "[default]", got[0])
	require.NotContains(t, got[1], awsKeyID)
	require.NotContains(t, got[2], awsSecretKey)
}

func Test_redactSecrets_multilinePEM(t *testing.T) {
	r, err := newRedactor(nil, nil)
	require.NoError(t, err)

	t.Run("complete block is fully redacted", func(t *testing.T) {
		got := r.redactHunk([]string{
			`key := ` + "`" + `-----BEGIN RSA PRIVATE KEY-----`,
			`MIIEowIBAAKCAQEA7bq0`,
			`-----END RSA PRIVATE KEY-----` + "`",
			`return key`,
		})
		require.Equal(t, "key := `[REDACTED]", got[0])
		require.Equal(t, "[REDACTED]", got[1])
		require.Equal(t, "[REDACTED]", got[2])
		require.Equal(t, "return key", got[3])
	})

	t.Run("block cut off by the hunk boundary is redacted to the end", func(t *testing.T) {
		got := r.redactHunk([]string{
			`-----BEGIN OPENSSH PRIVATE KEY-----`,
			`b3BlbnNzaC1rZXktdjEAAAAABG5vbmUA`,
			`AAAEC5BJHRnfmVLMdSe1BleTLLmRD3wD`,
		})
		require.Equal(t, []string{"[REDACTED]", "[REDACTED]", "[REDACTED]"}, got)
	})
}

func Test_newRedactor_customPatternsAndKeywords(t *testing.T) {
	r, err := newRedactor(
		[]string{`\bmyco_[A-Za-z0-9]{20}\b`},
		[]string{"connStr", "会員キー"},
	)
	require.NoError(t, err)

	tests := []struct {
		name string
		line string
		want string
	}{
		{
			name: "custom pattern redacts whole match",
			line: `client.Auth("myco_aBcDeFgHiJkLmNoPqRsT")`,
			want: `client.Auth("[REDACTED]")`,
		},
		{
			name: "custom keyword redacts quoted assignment",
			line: `dbConnStr = "postgres://user:pass@host/db"`,
			want: `dbConnStr = "[REDACTED]"`,
		},
		{
			name: "custom multibyte keyword redacts quoted assignment",
			line: `会員キー = "0123456789abcdef"`,
			want: `会員キー = "[REDACTED]"`,
		},
		{
			name: "default keywords still apply",
			line: `api_key = "0123456789abcdef"`,
			want: `api_key = "[REDACTED]"`,
		},
		{
			name: "unrelated line is unchanged",
			line: `connect(host, port)`,
			want: `connect(host, port)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, []string{tt.want}, r.redactHunk([]string{tt.line}))
		})
	}
}

func Test_newRedactor_invalidPattern(t *testing.T) {
	_, err := newRedactor([]string{`(unclosed`}, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid redaction pattern")
}

func Test_hunkForLine_redactsSecrets(t *testing.T) {
	secretLine := `apiKey := "0123456789abcdef"`
	lines := []string{secretLine, delimitedTestFlagKey, secretLine}

	r, err := newRedactor(nil, nil)
	require.NoError(t, err)
	matcher := Matcher{
		ctxLines: 1,
		redactor: r,
		Element:  NewElementMatcher("my-project", ``, `"`, []string{testFlagKey}, nil),
	}
	f := file{lines: lines}
	got := f.hunkForLine(testFlagKey, 1, matcher)
	require.NotNil(t, got)
	require.Equal(t, `apiKey := "[REDACTED]"`+"\n"+delimitedTestFlagKey+"\n"+`apiKey := "[REDACTED]"`, got.Lines)
	// the file's cached lines must not be modified by redaction
	require.Equal(t, secretLine, f.lines[0])

	// with redaction disabled, lines are sent as-is
	matcher.redactor = nil
	got = f.hunkForLine(testFlagKey, 1, matcher)
	require.NotNil(t, got)
	require.Equal(t, secretLine+"\n"+delimitedTestFlagKey+"\n"+secretLine, got.Lines)
}

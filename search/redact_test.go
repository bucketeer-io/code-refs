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
			name: "aws access key id",
			line: `awsKey := "AKIAIOSFODNN7EXAMPLE"`,
			want: `awsKey := "[REDACTED]"`,
		},
		{
			name: "github personal access token",
			line: `client := github.NewClient("ghp_aBcDeFgHiJkLmNoPqRsTuVwXyZ0123456789")`,
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
			line: `slack.New("xox` + `b-123456789012-abcdefghijklmnop")`,
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
			require.Equal(t, tt.want, r.redact(tt.line))
		})
	}
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
			require.Equal(t, tt.want, r.redact(tt.line))
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

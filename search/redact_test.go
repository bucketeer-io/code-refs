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
			// secret access key is nearby, see Test_redactSecrets_awsKeyPair.
			// fake tokens are assembled at runtime so secret scanners don't
			// flag them; the same applies to the other concatenated cases below
			name: "aws access key id assigned to key-named variable",
			line: `awsKey := "AKIA` + `Z9X24KQ7NW3JT6BP"`,
			want: `awsKey := "[REDACTED]"`,
		},
		{
			name: "github personal access token",
			line: `client := github.NewClient("ghp_` + `Wq7xK2mV9tRz3bJ8pNn4Y6dL1gF5hC0Tka9s")`,
			want: `client := github.NewClient("[REDACTED]")`,
		},
		{
			name: "github fine-grained token",
			line: `token := "github_pat_` + `11ABCDEFG0abcdefghijkl_mnopqrstuvwxyz"`,
			want: `token := "[REDACTED]"`,
		},
		{
			name: "gitlab personal access token",
			line: `curl --header "PRIVATE-TOKEN: glpat-` + `aBcDeFgHiJkLmNoPqRsT"`,
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
			line: `maps.NewClient(maps.WithAPIKey("AIza` + `SyA1bC2dE3fG4hI5jK6lM7nO8pQ9rS0tU1v"))`,
			want: `maps.NewClient(maps.WithAPIKey("[REDACTED]"))`,
		},
		{
			name: "jwt",
			line: `Authorization: eyJ` + `hbGciOiJIUzI1NiJ9.eyJ` + `zdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U`,
			want: `Authorization: [REDACTED]`,
		},
		{
			name: "bearer authorization header",
			line: `req.Header.Set("Authorization", "Bearer AbCdEfGhIjKlMnOpQrStUvWxYz")`,
			want: `req.Header.Set("Authorization", "Bearer [REDACTED]")`,
		},
		{
			// marker split across literals so scanners don't flag it
			name: "private key header",
			line: `-----BEGIN RSA PRIVATE ` + `KEY-----`,
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
			name: "password in connection url",
			line: `db = "postgres://admin:S3cr3tPw9xKq2m@db.internal:5432/prod"`,
			want: `db = "postgres://admin:[REDACTED]@db.internal:5432/prod"`,
		},
		{
			// userinfo ends at the last @ before the path, so none of the
			// password may leak
			name: "password containing @ in connection url",
			line: `db = "postgres://admin:p@ssw0rd9x@db.internal:5432/prod"`,
			want: `db = "postgres://admin:[REDACTED]@db.internal:5432/prod"`,
		},
		{
			name: "url without credentials is unchanged",
			line: `endpoint = "https://api.example.com:8443/v1/health"`,
			want: `endpoint = "https://api.example.com:8443/v1/health"`,
		},
		{
			// assembled at runtime so secret scanners don't flag the fake secret
			name: "unquoted assignment with slash-containing value is redacted",
			line: `aws_secret_access_key = ` + "q7PkzOxMjTn2" + "FvBw5RcY8dHl3sGe1uAiK9pXm4tE",
			want: `aws_secret_access_key = [REDACTED]`,
		},
		{
			name: "unquoted assignment with no spaces around equals is redacted",
			line: `DB_PASSWORD=S3cr3tPw9xKq2m`,
			want: `DB_PASSWORD=[REDACTED]`,
		},
		{
			name: "unquoted assignment to another identifier is not redacted",
			line: `apiKey := opts.APIKey`,
			want: `apiKey := opts.APIKey`,
		},
		{
			name: "unquoted assignment to a function call is not redacted",
			line: `tokenCount := len(tokens)`,
			want: `tokenCount := len(tokens)`,
		},
		{
			name: "unquoted all-alpha value is not redacted despite length",
			line: `secretMode := productionenvironment`,
			want: `secretMode := productionenvironment`,
		},
		{
			name: "unquoted value followed by a dotted selector is not redacted",
			line: `token := someFactory123456.Build()`,
			want: `token := someFactory123456.Build()`,
		},
		{
			name: "unquoted short value under the length floor is not redacted",
			line: `token = abc`,
			want: `token = abc`,
		},
		{
			name: "unquoted comparison is not redacted",
			line: `if token == expectedToken123 {`,
			want: `if token == expectedToken123 {`,
		},
		{
			name: "quoted comparison is not redacted",
			line: `if authType == "OAuth2.0-flow" {`,
			want: `if authType == "OAuth2.0-flow" {`,
		},
		{
			name: "strict-equality comparison is not redacted",
			line: `password === candidatePassword42`,
			want: `password === candidatePassword42`,
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

func Test_redactHunk_emptyInput(t *testing.T) {
	r, err := newRedactor(nil, nil)
	require.NoError(t, err)

	// the length invariant must hold for empty hunks too (contextLines < 0
	// leaves hunkLines nil); a nil input must not come back as [""]
	require.Empty(t, r.redactHunk(nil))
	require.Empty(t, r.redactHunk([]string{}))
}

func Test_redactSecrets_awsKeyPair(t *testing.T) {
	r, err := newRedactor(nil, nil)
	require.NoError(t, err)

	// gitleaks' aws-access-token rule matches the key ID standalone; the
	// paired secret access key is caught by the unquoted-assignment pattern
	// since it isn't quoted and contains characters gitleaks' own generic
	// rule excludes. The fake credentials are assembled at runtime so secret
	// scanners don't flag them.
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
		// markers split across literals so scanners don't flag them
		got := r.redactHunk([]string{
			`key := ` + "`" + `-----BEGIN RSA PRIVATE ` + `KEY-----`,
			`MIIEowIBAAKCAQEA7bq0`,
			`-----END RSA PRIVATE ` + `KEY-----` + "`",
			`return key`,
		})
		require.Equal(t, "key := `[REDACTED]", got[0])
		require.Equal(t, "[REDACTED]", got[1])
		require.Equal(t, "[REDACTED]", got[2])
		require.Equal(t, "return key", got[3])
	})

	t.Run("block cut off by the hunk boundary is redacted to the end", func(t *testing.T) {
		// marker split across literals so scanners don't flag it
		got := r.redactHunk([]string{
			`-----BEGIN OPENSSH PRIVATE ` + `KEY-----`,
			`b3BlbnNzaC1rZXktdjEAAAAABG5vbmUA`,
			`AAAEC5BJHRnfmVLMdSe1BleTLLmRD3wD`,
		})
		require.Equal(t, []string{"[REDACTED]", "[REDACTED]", "[REDACTED]"}, got)
	})

	t.Run("block cut off at the top of the hunk is redacted from the start", func(t *testing.T) {
		// marker split across literals so scanners don't flag it
		got := r.redactHunk([]string{
			`MIIEowIBAAKCAQEA7bq0qzO5s7fVXygbYtNZ`,
			`dGVzdGtleWZha2VkYXRhMTIzNDU2Nzg5MGFi`,
			`-----END RSA PRIVATE ` + `KEY-----` + "`",
			`return key`,
		})
		require.Equal(t, []string{"[REDACTED]", "[REDACTED]", "[REDACTED]", "return key"}, got)
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

	// patterns matching the empty string would expand every line
	for _, p := range []string{``, `a*`} {
		_, err = newRedactor([]string{p}, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "must not match the empty string")
	}
}

func Test_aggregateHunksForFlag_redactsSecrets(t *testing.T) {
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
	got := f.aggregateHunksForFlag(testFlagKey, matcher, []int{1})
	require.Len(t, got, 1)
	require.Equal(t, `apiKey := "[REDACTED]"`+"\n"+delimitedTestFlagKey+"\n"+`apiKey := "[REDACTED]"`, got[0].Lines)
	// the file's cached lines must not be modified by redaction
	require.Equal(t, secretLine, f.lines[0])

	// with redaction disabled, lines are sent as-is
	matcher.redactor = nil
	got = f.aggregateHunksForFlag(testFlagKey, matcher, []int{1})
	require.Len(t, got, 1)
	require.Equal(t, secretLine+"\n"+delimitedTestFlagKey+"\n"+secretLine, got[0].Lines)
}

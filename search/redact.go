package search

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/zricethezav/gitleaks/v8/detect"
)

const redactedPlaceholder = "[REDACTED]"

// secretPattern pairs a compiled regex with the replacement applied to each
// match. Patterns that capture groups keep the non-sensitive portion (e.g. the
// variable name or the "Bearer" keyword) so the hunk stays readable.
type secretPattern struct {
	regex       *regexp.Regexp
	replacement string
}

// Vendor-specific credential formats are detected by the gitleaks ruleset.
// The patterns below cover generic material that a leak scanner deliberately
// leaves to context-aware tools: authorization header values, passwords in
// URLs, and quoted assignments to secret-looking variable names.
var authHeaderPattern = secretPattern{
	regex:       regexp.MustCompile(`(?i)\b(bearer|basic)\s+[A-Za-z0-9\-._~+/]{16,}=*`),
	replacement: "$1 " + redactedPlaceholder,
}

// Passwords in URL userinfo, e.g. postgres://admin:hunter2@db:5432/prod.
// The password class allows '@' (userinfo ends at the last '@' before the
// path) so a password containing a raw '@' is not partially leaked; quotes
// stop the match at the end of a string literal.
var urlCredentialsPattern = secretPattern{
	regex:       regexp.MustCompile(`(://[^:/?#@\s]+:)[^/?#\s"'` + "`" + `]+@`),
	replacement: "$1" + redactedPlaceholder + "@",
}

var (
	pemHeaderRegex = regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY( BLOCK)?-----`)
	pemFooterRegex = regexp.MustCompile(`-----END [A-Z ]*PRIVATE KEY( BLOCK)?-----`)
)

// defaultKeywords are matched as substrings of variable names in generic
// secret-looking assignments, e.g. `apiKey = "..."`, `password: '...'`
var defaultKeywords = []string{"api[_-]?key", "secret", "token", "passwd", "password", "credential", "auth"}

type redactor struct {
	detector                  *detect.Detector
	patterns                  []secretPattern
	unquotedAssignmentPattern *regexp.Regexp
}

// newRedactor builds a redactor from the gitleaks default ruleset,
// user-provided regexes (whole match is replaced), and user-provided keywords
// merged with the default keywords for generic assignment matching.
func newRedactor(customPatterns, customKeywords []string) (*redactor, error) {
	detector, err := detect.NewDetectorDefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load gitleaks ruleset: %w", err)
	}

	patterns := make([]secretPattern, 0, len(customPatterns)+5) //nolint:mnd
	patterns = append(patterns, authHeaderPattern, urlCredentialsPattern)

	for _, p := range customPatterns {
		regex, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("invalid redaction pattern %q: %w", p, err)
		}
		// a pattern matching the empty string (e.g. "" or "a*") would insert
		// the placeholder between every character of every scanned line
		if regex.MatchString("") {
			return nil, fmt.Errorf("invalid redaction pattern %q: pattern must not match the empty string", p)
		}
		patterns = append(patterns, secretPattern{regex, redactedPlaceholder})
	}

	keywords := make([]string, 0, len(defaultKeywords)+len(customKeywords))
	keywords = append(keywords, defaultKeywords...)
	for _, k := range customKeywords {
		if k == "" {
			continue
		}
		keywords = append(keywords, regexp.QuoteMeta(k))
	}
	alternation := strings.Join(keywords, "|")
	// Assignment separators: `:=`, `=`, or `:`. Spelled out as an alternation
	// (rather than [:=]+) so comparison operators like `==` don't match.
	const assignOp = `(?::=|[:=])`
	// One pattern per quote character since RE2 does not support backreferences.
	for _, q := range []string{`"`, `'`, "`"} {
		patterns = append(patterns, secretPattern{
			regex:       regexp.MustCompile(`(?i)([\w-]*(` + alternation + `)[\w-]*\s*` + assignOp + `\s*` + q + `)[^` + q + `]{6,}` + q),
			replacement: "$1" + redactedPlaceholder + q,
		})
	}

	// Catches unquoted assignments (e.g. .env/.ini/credentials-file style)
	// that gitleaks' own rules and the quoted patterns above miss. The value
	// may not start with `=` so the tail of a `==` comparison never matches.
	unquotedAssignmentPattern := regexp.MustCompile(
		`(?i)([\w-]*(` + alternation + `)[\w-]*\s*` + assignOp + `\s*)([A-Za-z0-9+/_-][A-Za-z0-9+/_=-]{11,})`, //nolint:mnd
	)

	return &redactor{detector: detector, patterns: patterns, unquotedAssignmentPattern: unquotedAssignmentPattern}, nil
}

// redactHunk scrubs values that look like credentials (API keys, tokens,
// passwords, private keys) from a hunk so they are not sent to Bucketeer.
// Detection runs on the whole hunk so multi-line secrets such as PEM blocks
// are caught; the returned slice always has the same length as the input.
func (r *redactor) redactHunk(lines []string) []string {
	// guard the length invariant: splitting the joined empty string below
	// would yield [""] (length 1) for an empty input
	if len(lines) == 0 {
		return lines
	}
	joined := strings.Join(lines, "\n")
	for _, finding := range r.detector.DetectString(joined) {
		joined = redactSecret(joined, finding.Secret)
	}
	out := strings.Split(joined, "\n")

	r.redactPEMBlocks(out)

	for i, line := range out {
		out[i] = redactUnquotedAssignmentLine(r.unquotedAssignmentPattern, line)
	}

	for i := range out {
		for _, p := range r.patterns {
			out[i] = p.regex.ReplaceAllString(out[i], p.replacement)
		}
	}
	return out
}

// redactSecret replaces every occurrence of a (possibly multi-line) secret
// with a placeholder that preserves its newline count, so hunk lines keep
// mapping to the same file line numbers.
func redactSecret(text, secret string) string {
	if secret == "" {
		return text
	}
	count := strings.Count(secret, "\n") + 1
	parts := make([]string, count)
	for i := range parts {
		parts[i] = redactedPlaceholder
	}
	return strings.ReplaceAll(text, secret, strings.Join(parts, "\n"))
}

// looksLikeSecretValue reports whether v is shaped like a credential rather
// than an ordinary identifier, boolean, or English word: it must contain both
// a letter and a digit.
func looksLikeSecretValue(v string) bool {
	var hasLetter, hasDigit bool
	for _, c := range v {
		switch {
		case c >= '0' && c <= '9':
			hasDigit = true
		case (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z'):
			hasLetter = true
		}
		if hasLetter && hasDigit {
			return true
		}
	}
	return false
}

// redactUnquotedAssignmentLine redacts unquoted secret-shaped values (e.g.
// aws_secret_access_key = wJalr.../K7...) that gitleaks' generic-api-key rule
// and our own quoted-assignment patterns miss because they contain characters
// like `/` outside its charset. RE2 has no lookahead, so a value immediately
// followed by `.` or `(` -- e.g. `apiKey := opts.APIKey` or
// `tokenCount := len(tokens)` -- is excluded by inspecting match indices
// directly rather than via the regex itself.
func redactUnquotedAssignmentLine(re *regexp.Regexp, line string) string {
	matches := re.FindAllStringSubmatchIndex(line, -1)
	if matches == nil {
		return line
	}
	var b strings.Builder
	last := 0
	for _, m := range matches {
		valStart, valEnd := m[6], m[7]
		redact := looksLikeSecretValue(line[valStart:valEnd])
		if redact && valEnd < len(line) {
			if next := line[valEnd]; next == '.' || next == '(' {
				redact = false
			}
		}
		b.WriteString(line[last:m[0]])
		if redact {
			b.WriteString(line[m[0]:m[3]])
			b.WriteString(redactedPlaceholder)
		} else {
			b.WriteString(line[m[0]:m[1]])
		}
		last = m[1]
	}
	b.WriteString(line[last:])
	return b.String()
}

// redactPEMBlocks redacts everything from a private key header to the
// matching footer, or to the end of the hunk when the block is cut off by the
// hunk boundary. This backstops the gitleaks private-key rule, which needs
// the footer to be present to match.
func (r *redactor) redactPEMBlocks(lines []string) {
	inBlock := false
	for i, line := range lines {
		if inBlock {
			lines[i] = redactedPlaceholder
			if pemFooterRegex.MatchString(line) {
				inBlock = false
			}
			continue
		}
		if loc := pemHeaderRegex.FindStringIndex(line); loc != nil {
			lines[i] = line[:loc[0]] + redactedPlaceholder
			if !pemFooterRegex.MatchString(line[loc[1]:]) {
				inBlock = true
			}
			continue
		}
		// a footer without a preceding header means the block started before
		// the hunk boundary; redact from the top of the hunk to the footer
		if pemFooterRegex.MatchString(line) {
			for j := 0; j <= i; j++ {
				lines[j] = redactedPlaceholder
			}
		}
	}
}

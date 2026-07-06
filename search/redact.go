package search

import (
	"fmt"
	"regexp"
	"strings"
)

const redactedPlaceholder = "[REDACTED]"

// secretPattern pairs a compiled regex with the replacement applied to each
// match. Patterns that capture groups keep the non-sensitive portion (e.g. the
// variable name or the "Bearer" keyword) so the hunk stays readable.
type secretPattern struct {
	regex       *regexp.Regexp
	replacement string
}

// basePatterns match well-known credential formats. Patterns are applied per
// line, in order, so specific token formats are redacted before the generic
// assignment fallback.
var basePatterns = []secretPattern{
	// Private key material, e.g. "-----BEGIN RSA PRIVATE KEY-----"
	{regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY( BLOCK)?-----.*`), redactedPlaceholder},
	// AWS access key IDs
	{regexp.MustCompile(`\b(A3T[A-Z0-9]|AKIA|ASIA|ABIA|ACCA)[A-Z0-9]{16}\b`), redactedPlaceholder},
	// GitHub tokens (classic and fine-grained)
	{regexp.MustCompile(`\b(ghp|gho|ghu|ghs|ghr)_[A-Za-z0-9]{36,}\b`), redactedPlaceholder},
	{regexp.MustCompile(`\bgithub_pat_[A-Za-z0-9_]{22,}\b`), redactedPlaceholder},
	// GitLab personal access tokens
	{regexp.MustCompile(`\bglpat-[A-Za-z0-9_-]{20,}\b`), redactedPlaceholder},
	// Slack tokens
	{regexp.MustCompile(`\bxox[baprs]-[A-Za-z0-9-]{10,}\b`), redactedPlaceholder},
	// Stripe secret/restricted keys
	{regexp.MustCompile(`\b[sr]k_(live|test)_[A-Za-z0-9]{16,}\b`), redactedPlaceholder},
	// Google API keys
	{regexp.MustCompile(`\bAIza[0-9A-Za-z_-]{35}\b`), redactedPlaceholder},
	// JSON web tokens
	{regexp.MustCompile(`\beyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{5,}\b`), redactedPlaceholder},
	// Authorization header values, e.g. `Authorization: Bearer <token>`
	{regexp.MustCompile(`(?i)\b(bearer|basic)\s+[A-Za-z0-9\-._~+/]{16,}=*`), "$1 " + redactedPlaceholder},
}

// defaultKeywords are matched as substrings of variable names in generic
// secret-looking assignments, e.g. `apiKey = "..."`, `password: '...'`
var defaultKeywords = []string{"api[_-]?key", "secret", "token", "passwd", "password", "credential", "auth"}

type redactor struct {
	patterns []secretPattern
}

// newRedactor builds a redactor from the base patterns, user-provided regexes
// (whole match is replaced), and user-provided keywords merged with the
// default keywords for generic assignment matching.
func newRedactor(customPatterns, customKeywords []string) (*redactor, error) {
	patterns := make([]secretPattern, 0, len(basePatterns)+len(customPatterns)+3) //nolint:mnd
	patterns = append(patterns, basePatterns...)

	for _, p := range customPatterns {
		regex, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("invalid redaction pattern %q: %w", p, err)
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
	// One pattern per quote character since RE2 does not support backreferences.
	for _, q := range []string{`"`, `'`, "`"} {
		patterns = append(patterns, secretPattern{
			regex:       regexp.MustCompile(`(?i)([\w-]*(` + alternation + `)[\w-]*\s*[:=]+\s*` + q + `)[^` + q + `]{6,}` + q),
			replacement: "$1" + redactedPlaceholder + q,
		})
	}

	return &redactor{patterns: patterns}, nil
}

// redact scrubs values that look like credentials (API keys, tokens,
// passwords, private keys) from a line so they are not sent to Bucketeer.
func (r *redactor) redact(line string) string {
	for _, p := range r.patterns {
		line = p.regex.ReplaceAllString(line, p.replacement)
	}
	return line
}

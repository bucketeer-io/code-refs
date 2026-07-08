package search

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/betterleaks/betterleaks/detect"
)

const redactedPlaceholder = "[REDACTED]"

// secretPattern pairs a compiled regex with the replacement applied to each
// match. Patterns that capture groups keep the non-sensitive portion (e.g. the
// variable name or the "Bearer" keyword) so the hunk stays readable.
type secretPattern struct {
	regex       *regexp.Regexp
	replacement string
}

// Vendor-specific credential formats are detected by the betterleaks ruleset.
// The patterns below cover generic material that a leak scanner deliberately
// leaves to context-aware tools: authorization header values, passwords in
// URLs, and quoted assignments to secret-looking variable names.
var authHeaderPattern = secretPattern{
	regex:       regexp.MustCompile(`(?i)\b(bearer|basic)\s+[A-Za-z0-9\-._~+/]{16,}=*`),
	replacement: "$1 " + redactedPlaceholder,
}

// Passwords in URL userinfo, e.g. postgres://admin:hunter2@db:5432/prod
var urlCredentialsPattern = secretPattern{
	regex:       regexp.MustCompile(`(://[^:/?#@\s]+:)[^@/?#\s]+@`),
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
	detector *detect.Detector
	patterns []secretPattern
}

// newRedactor builds a redactor from the betterleaks default ruleset,
// user-provided regexes (whole match is replaced), and user-provided keywords
// merged with the default keywords for generic assignment matching.
func newRedactor(customPatterns, customKeywords []string) (*redactor, error) {
	detector, err := detect.NewDetectorDefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load betterleaks ruleset: %w", err)
	}

	patterns := make([]secretPattern, 0, len(customPatterns)+5) //nolint:mnd
	patterns = append(patterns, authHeaderPattern, urlCredentialsPattern)

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

	return &redactor{detector: detector, patterns: patterns}, nil
}

// redactHunk scrubs values that look like credentials (API keys, tokens,
// passwords, private keys) from a hunk so they are not sent to Bucketeer.
// Detection runs on the whole hunk so multi-line secrets such as PEM blocks
// are caught; the returned slice always has the same length as the input.
func (r *redactor) redactHunk(lines []string) []string {
	joined := strings.Join(lines, "\n")
	for _, finding := range r.detector.DetectString(joined) {
		joined = redactSecret(joined, finding.Secret)
		// composite rules (e.g. AWS key ID + secret access key) report the
		// components separately; redact those too
		for _, set := range finding.RequiredSets {
			for _, component := range set.Components {
				joined = redactSecret(joined, component.Secret)
			}
		}
	}
	out := strings.Split(joined, "\n")

	r.redactPEMBlocks(out)

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

// redactPEMBlocks redacts everything from a private key header to the
// matching footer, or to the end of the hunk when the block is cut off by the
// hunk boundary. This backstops the betterleaks private-key rule, which needs
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
		}
	}
}

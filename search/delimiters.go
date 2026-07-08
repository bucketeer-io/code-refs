package search

import (
	"github.com/bucketeer-io/code-refs/internal/helpers"
	"github.com/bucketeer-io/code-refs/options"
)

// Get a list of delimiters to use for flag key matching
// If defaults are disabled, only additional configured delimiters will be used
func GetDelimiters(opts options.Options) []string {
	delims := make([]string, 0, 3+len(opts.Delimiters.Additional)) //nolint:mnd
	if !opts.Delimiters.DisableDefaults {
		delims = append(delims, `"`, `'`, "`")
	}

	delims = append(delims, opts.Delimiters.Additional...)

	return helpers.Dedupe(delims)
}

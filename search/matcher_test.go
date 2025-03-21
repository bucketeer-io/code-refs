package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_buildFlagPatterns(t *testing.T) {
	testFlagKey := "testflag"
	patterns := buildElementPatterns([]string{testFlagKey}, defaultDelims)
	want := map[string][]string{"testflag": {"\"testflag\"", "\"testflag'", "\"testflag`", "'testflag\"", "'testflag'", "'testflag`", "`testflag\"", "`testflag'", "`testflag`"}}
	require.Equal(t, want, patterns)
}

func TestElementMatcher_FindAliases(t *testing.T) {
	t.Run("overlapping aliases are reported separately", func(t *testing.T) {
		matcher := NewElementMatcher("project", "", "", nil, map[string][]string{"flag": {"alias", "alias1"}})
		assert.ElementsMatch(t, []string{"alias", "alias1"}, matcher.FindAliases("alias1", "flag"))
	})
}

func TestElementMatcher_FindMatches(t *testing.T) {
	t.Run("overlapping flags are reported separately", func(t *testing.T) {
		matcher := NewElementMatcher("project", "", "", []string{"flag", "flag1"}, nil)
		assert.ElementsMatch(t, []string{"flag", "flag1"}, matcher.FindMatches("flag1"))
	})
}

func TestMatcher_MatchElement(t *testing.T) {
	specs := []struct {
		name     string
		expected bool
		line     string
		matcher  Matcher
		flagKey  string
	}{
		{
			name:     "match found",
			expected: true,
			line:     "var flagKey = 'testflag'",
			matcher:  Matcher{Element: NewElementMatcher("projKey", "", ",'\"", []string{"testflag"}, map[string][]string{"testflag": {"testFlag"}})},
			flagKey:  "testflag",
		},
		{
			name:     "no match found",
			expected: false,
			line:     "var flagKey = 'testflag'",
			matcher:  Matcher{Element: NewElementMatcher("projKey", "", ",'\"", []string{"anotherflag"}, map[string][]string{"anotherflag": {"anotherFlag"}})},
			flagKey:  "testflag",
		},
		{
			name:     "doesn't match when delimiters aren't present",
			expected: false,
			line:     "var TEST_FLAG",
			matcher:  Matcher{Element: NewElementMatcher("projKey", "", "'", []string{"TEST_FLAG"}, map[string][]string{"testflag": {}})},
			flagKey:  "TEST_FLAG",
		},
		{
			// This example is equivalent to setting `defaultsDisabled: true` in delimiters configuration
			name:     "matches without delimiters",
			expected: true,
			line:     "var TEST_FLAG",
			matcher:  Matcher{Element: NewElementMatcher("projKey", "", "", []string{"TEST_FLAG"}, map[string][]string{"testflag": {}})},
			flagKey:  "TEST_FLAG",
		},
	}

	for _, tt := range specs {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.matcher.MatchElement(tt.line, tt.flagKey))
		})
	}
}

package matcher

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegExpMatches(t *testing.T) {
	tcs := []struct {
		name       string
		expression any
		value      any
		expected   bool
	}{
		{"string pattern", "tEsT", "tEsT", true},
		{"string pattern using a non string argument", "10", 10, true},
		{"regular expression (pointer)", regexp.MustCompile("tEsT"), "tEsT", true},
		{"regular expression (non-pointer)", *regexp.MustCompile("tEsT"), "tEsT", true},
		{"regular expression (does not match)", regexp.MustCompile("tEsT"), "dev", false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Matches(tc.expression).Match(tc.value)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result.Pass)
		})
	}
}

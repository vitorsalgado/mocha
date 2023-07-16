package matcher

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegExpMatches(t *testing.T) {
	tcs := []struct {
		expression string
		value      any
		expected   bool
	}{
		{"tEsT", "tEsT", true},
		{"(?mi)hi", "/test?q=bye", false},
		{"(?mi)hi", "/test?q=hi", true},
	}

	for i, tc := range tcs {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			result, err := Matches(tc.expression).Match(tc.value)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result.Pass)
		})
	}
}

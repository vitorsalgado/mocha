package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHasPrefix(t *testing.T) {
	tcs := []struct {
		value    string
		expected bool
	}{
		{"hello", true},
		{"world", false},
	}

	for _, tc := range tcs {
		t.Run(tc.value, func(t *testing.T) {
			result, err := HasPrefix(tc.value).Match("hello world")
			require.NoError(t, err)
			require.Equal(t, tc.expected, result.Pass)
		})
	}
}

func TestHasPrefixMatcher_Name(t *testing.T) {
	require.NotEmpty(t, HasPrefix("").Name())
}

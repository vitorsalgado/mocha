package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHasSuffix(t *testing.T) {
	tcs := []struct {
		value    string
		expected bool
	}{
		{"world", true},
		{"hello", false},
	}

	for _, tc := range tcs {
		t.Run(tc.value, func(t *testing.T) {
			result, err := HasSuffix(tc.value).Match("hello world")

			require.NoError(t, err)
			require.Equal(t, tc.expected, result.Pass)
		})
	}
}

func TestHasSuffixMatcher_Name(t *testing.T) {
	require.NotEmpty(t, HasSuffix("").Name())
}

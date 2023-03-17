package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLen_String(t *testing.T) {
	testCases := []struct {
		name     string
		value    any
		size     int
		expected bool
	}{
		{"string", "hello world -  ", 15, true},
		{"array", []string{"hi", "bye"}, 2, true},
		{"array (no match)", []string{"hi", "bye"}, 1, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Len(tc.size).Match(tc.value)

			require.NoError(t, err)
			require.Equal(t, tc.expected, result.Pass)
		})
	}
}

func TestLenMatcher_Name(t *testing.T) {
	require.NotEmpty(t, Len(10).Name())
}

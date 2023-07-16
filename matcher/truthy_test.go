package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTruthy(t *testing.T) {
	testCases := []struct {
		name     string
		value    any
		expected bool
	}{
		{"bool", true, true},
		{"bool (false)", false, false},
		{"string (bool)", "true", true},
		{"string (bool) (2)", "false", false},
		{"string (txt)", "t", true},
		{"int", 0, false},
		{"int (2)", 1, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Truthy().Match(tc.value)

			require.NoError(t, err)
			require.Equal(t, tc.expected, result.Pass)
		})
	}
}

func TestTruthyInvalidSyntax(t *testing.T) {
	result, err := Truthy().Match("y")

	require.Error(t, err)
	require.False(t, result.Pass)
}

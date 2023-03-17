package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFalsy(t *testing.T) {
	testCases := []struct {
		name     string
		value    any
		expected bool
	}{
		{"bool", true, false},
		{"bool (false)", false, true},
		{"string (bool)", "true", false},
		{"string (bool) (2)", "false", true},
		{"string (txt)", "t", false},
		{"int", 0, true},
		{"int (2)", 1, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Falsy().Match(tc.value)

			require.NoError(t, err)
			require.Equal(t, tc.expected, result.Pass)
		})
	}
}

func TestFalsyInvalidSyntax(t *testing.T) {
	result, err := Truthy().Match("y")

	require.Error(t, err)
	require.Nil(t, result)
}

func TestFalsyMatcher_Name(t *testing.T) {
	require.NotEmpty(t, Falsy().Name())
}

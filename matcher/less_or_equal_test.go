package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrEqualLess(t *testing.T) {
	testCases := []struct {
		name     string
		expected float64
		value    any
		result   bool
	}{
		{"is less", 10, 5, true},
		{"is less fl", 9.9, 9.8, true},
		{"is equal", 10, 10, true},
		{"is equal fl", 9.9, 9.9, true},
		{"greater", 5, 10, false},
		{"greater (string)", 5, "10", false},
		{"greater (float32)", 5, float32(10), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := LessOrEqual(tc.expected).Match(tc.value)

			require.NoError(t, err)
			require.Equal(t, tc.result, res.Pass)
		})
	}
}

func TestLessOrEqualUnhandledType(t *testing.T) {
	res, err := LessOrEqual(10).Match(true)

	require.Error(t, err)
	require.Nil(t, res)
}

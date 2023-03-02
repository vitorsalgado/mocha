package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLess(t *testing.T) {
	testCases := []struct {
		name     string
		expected float64
		value    any
		result   bool
	}{
		{"is less", 10, 5, true},
		{"is less fl", 9.9, 9.8, true},
		{"equal", 10, 10, false},
		{"greater", 5, 10, false},
		{"is less (string)", 10, "5", true},
		{"is less (float32)", 10, float32(5), true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := LessThan(tc.expected).Match(tc.value)

			require.NoError(t, err)
			require.Equal(t, tc.result, res.Pass)
		})
	}
}

func TestLessUnhandledType(t *testing.T) {
	res, err := Lt(10).Match(true)

	require.Error(t, err)
	require.Nil(t, res)
}

package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGreaterOrEqual(t *testing.T) {
	testCases := []struct {
		name     string
		expected float64
		value    any
		result   bool
	}{
		{"is less", 10, 5, false},
		{"is less fl", 9.9, 9.8, false},
		{"equal", 10, 10, true},
		{"greater", 5, 10, true},
		{"greater fl", 5.5, 6.5, true},
		{"greater fl (string)", 5.5, "6.5", true},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res, err := GreaterThanOrEqual(tc.expected).Match(tc.value)

			require.NoError(t, err)
			require.Equal(t, tc.result, res.Pass)
		})
	}
}

func TestGreaterOrEqualUnhandledType(t *testing.T) {
	res, err := Gte(10).Match(true)

	require.Error(t, err)
	require.False(t, res.Pass)
}

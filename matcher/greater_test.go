package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGreater(t *testing.T) {
	testCases := []struct {
		name     string
		expected float64
		value    any
		result   bool
	}{
		{"is less", 10, 5, false},
		{"is less fl", 9.9, 9.8, false},
		{"equal", 10, 10, false},
		{"greater", 5, 10, true},
		{"greater fl", 5.5, 6.5, true},
		{"greater fl (string)", 5.5, "6.5", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := GreaterThan(tc.expected).Match(tc.value)

			require.NoError(t, err)
			require.Equal(t, tc.result, res.Pass)
		})
	}
}

func TestGreaterUnhandledType(t *testing.T) {
	res, err := Gt(10).Match(true)

	require.Error(t, err)
	require.Nil(t, res)
}

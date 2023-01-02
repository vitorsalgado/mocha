package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGreaterOrEqual(t *testing.T) {
	testCases := []struct {
		name     string
		expected float64
		value    float64
		result   bool
	}{
		{"is less", 10, 5, false},
		{"is less fl", 9.9, 9.8, false},
		{"equal", 10, 10, true},
		{"greater", 5, 10, true},
		{"greater fl", 5.5, 6.5, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := GreaterOrEqualThan(tc.expected).Match(tc.value)

			assert.NoError(t, err)
			assert.Equal(t, tc.result, res.Pass)
		})
	}
}

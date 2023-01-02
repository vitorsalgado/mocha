package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLess(t *testing.T) {
	testCases := []struct {
		name     string
		expected float64
		value    float64
		result   bool
	}{
		{"is less", 10, 5, true},
		{"is less fl", 9.9, 9.8, true},
		{"equal", 10, 10, false},
		{"greater", 5, 10, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := LessThan(tc.expected).Match(tc.value)

			assert.NoError(t, err)
			assert.Equal(t, tc.result, res.Pass)
		})
	}
}

package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	testCases := []struct {
		name     string
		expected any
		value    any
		result   bool
	}{
		{
			"should return true when value is contained in the matcher argument",
			"world",
			"hello world",
			true,
		},
		{
			"should return false when value is not contained in the matcher argument",
			"dev",
			"hello world",
			false,
		},
		{
			"should return true when expected value is contained in the given slice",
			"dev",
			[]string{"dev", "qa"},
			true,
		},
		{
			"should return false when expected value is not contained in the given slice",
			"po",
			[]string{"dev", "qa"},
			false,
		},
		{
			"should return true when expected value is a key present in the given map",
			"dev",
			map[string]string{"dev": "ok", "qa": "nok"},
			true,
		},
		{
			"should return false when expected value is a key not present in the given map",
			"unknown",
			map[string]string{"dev": "ok", "qa": "nok"},
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Contain(tc.expected).Match(tc.value)
			assert.Nil(t, err)
			assert.Equal(t, tc.result, result.OK)
		})
	}
}

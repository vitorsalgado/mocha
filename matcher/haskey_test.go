package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHasKey(t *testing.T) {
	m := map[string]any{
		"name": "someone",
		"age":  34,
		"address": map[string]any{
			"street": "very nice place",
		},
		"options": []string{},
		"active":  false,
		"none":    nil,
		"zero":    0,
	}

	tcs := []struct {
		key      string
		expected bool
	}{
		{"name", true},
		{"age", true},
		{"address", true},
		{"address.street", true},
		{"address.city", false},
		{"active", true},
		{"zero", true},
		{"options", true},
		{"none", false},
	}

	for _, tc := range tcs {
		t.Run(tc.key, func(t *testing.T) {
			result, err := HasKey(tc.key).Match(m)
			require.Equal(t, tc.expected, result.Pass)
			require.Nil(t, err)
		})
	}
}

func TestHasKeyMatcher_Name(t *testing.T) {
	require.NotEmpty(t, HasKey("").Name())
}

func TestHasKeyNew(t *testing.T) {
	require.Panics(t, func() {
		HasKey(".")
	})
}

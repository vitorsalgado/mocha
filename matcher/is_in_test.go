package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsIn(t *testing.T) {
	tcs := []struct {
		name     string
		items    []any
		value    any
		expected bool
	}{
		{"numbers", []any{10, 20, 30}, 20, true},
		{"string", []any{"10", "20", "30"}, "20", true},
		{"string -- number", []any{"10", "20", "30"}, 20, false},
		{"mixed", []any{"city", 100, true, false, 2000, "all", "test"}, "test", true},
		{"mixed", []any{"city", 100, true, false, 2000, "all", "test"}, "dev", false},
		{"mixed", []any{"city", 100, true, false, 2000, "all", "test"}, true, true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result, err := IsIn(tc.items).Match(tc.value)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result.Pass)
		})
	}
}

func TestInInWithInvalidItems(t *testing.T) {
	res, err := IsIn(true).Match(true)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestIsContainedInMatcher_Name(t *testing.T) {
	require.NotEmpty(t, IsIn([]string{}).Name())
}

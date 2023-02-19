package matcher

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSome(t *testing.T) {
	tcs := []struct {
		name     string
		matcher  Matcher
		items    []any
		expected bool
	}{
		{"numbers", Equal(10), []any{10, 20, 30}, true},
		{"string", Contain("30"), []any{"10", "20", "30"}, true},
		{"string -- number", Equal(40), []any{"10", "20", "30"}, false},
		{"mixed", Equal("city"), []any{"city", 100, true, false, 2000, "all", "test"}, true},
		{"mixed", Equal("none"), []any{"city", 100, true, false, 2000, "all", "test"}, false},
		{"mixed", Equal(100), []any{"city", 100, true, false, 2000, "all", "test"}, true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Some(tc.matcher).Match(tc.items)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result.Pass)
		})
	}
}

func TestSomeErr(t *testing.T) {
	res, err := Some(Func(func(_ any) (bool, error) {
		return false, errors.New("boom")
	})).Match(true)

	require.Error(t, err)
	require.Nil(t, res)
}

package matcher

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestItem(t *testing.T) {
	type model struct {
		name string
	}

	v := []any{"test", "dev", true, false, []string{"qa", "none"}, model{"go"}, 100}

	testCases := []struct {
		index    int
		matcher  Matcher
		expected bool
	}{
		{0, Equal("test"), true},
		{1, HasPrefix("d"), true},
		{1, Equal("test"), false},
		{2, Falsy(), false},
		{3, Falsy(), true},
		{4, Item(0, Equal("qa")), true},
		{5, Equal(model{"go"}), true},
		{6, GreaterThan(10), true},
		{6, LessThanOrEqual(10), false},
	}

	for _, tc := range testCases {
		t.Run(strconv.Itoa(tc.index), func(t *testing.T) {
			res, err := Item(tc.index, tc.matcher).Match(v)

			require.NoError(t, err)
			require.Equal(t, tc.expected, res.Pass)
		})
	}
}

func TestItemInvalidIndex(t *testing.T) {
	res, err := Item(-1, Equal("dev")).Match([]string{"dev"})
	require.Error(t, err)
	require.False(t, res.Pass)
}

func TestInvalidType(t *testing.T) {
	res, err := Item(0, Equal("dev")).Match(true)
	require.Error(t, err)
	require.False(t, res.Pass)
}

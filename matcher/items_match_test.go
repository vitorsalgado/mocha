package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEqualElementsMatcher(t *testing.T) {
	testCases := []struct {
		name          string
		value         any
		expectedValue any
		expected      bool
	}{
		{"arrays", []string{"a", "b", "c"}, []string{"c", "b", "a"}, true},
		{"arrays (more items)", []string{"a", "b", "c"}, []string{"c", "b", "a", "d", "e"}, false},
		{"arrays (more items)", []string{"c", "b", "a", "d", "e"}, []string{"a", "b", "c"}, false},
		{"arrays (repetitive)", []string{"a", "a", "b", "c"}, []string{"c", "b", "a", "a"}, true},
		{"arrays (repetitive)", []string{"a", "a", "b", "c"}, []string{"c", "b", "a", "a", "a"}, false},
		{"arrays (repetitive)", []string{"a", "a", "a", "b", "b", "c"}, []string{"c", "b", "a", "a", "a"}, false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res, err := ItemsMatch(tc.value).Match(tc.expectedValue)

			require.NoError(t, err)
			require.Equal(t, tc.expected, res.Pass)
		})
	}
}

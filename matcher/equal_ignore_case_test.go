package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEqualIgnoreCase(t *testing.T) {
	tcs := []struct {
		name              string
		expectedValue     string
		expectedValueArgs []any
		matchValue        any
		expected          bool
	}{
		{"nil value (empty)", "", nil, nil, true},
		{"nil value", "test", nil, nil, false},
		{"diff case", "TesT", nil, "test", true},
		{"equal values", "test", nil, "test", true},
		{"diff values", "TeST", nil, "TeST DEV", false},
		{"format diff case", "hello %s", []any{"WORLD"}, "hello world", true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			var result Result
			var err error

			if len(tc.expectedValueArgs) > 0 {
				result, err = EqualIgnoreCasef(tc.expectedValue, tc.expectedValueArgs...).Match(tc.matchValue)
			} else {
				result, err = Eqi(tc.expectedValue).Match(tc.matchValue)
			}

			require.NoError(t, err)
			require.Equal(t, tc.expected, result.Pass)
		})
	}
}

package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNil(t *testing.T) {
	tcs := []struct {
		name     string
		value    any
		expected bool
	}{
		{"nil", nil, true},
		{"empty string", "", false},
		{"zero", 0, false},
		{"bool", false, false},
		{"text", "txt", false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res, err := Nil().Match(tc.value)
			require.NoError(t, err)
			require.Equal(t, tc.expected, res.Pass)
		})
	}
}

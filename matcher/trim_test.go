package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrim(t *testing.T) {
	result, err := Trim(StrictEqual("test")).Match("  test  ")

	require.NoError(t, err)
	require.True(t, result.Pass)
}

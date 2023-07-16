package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToUpperCase(t *testing.T) {
	result, err := ToUpper(StrictEqual("TEST")).Match("tesT")

	require.NoError(t, err)
	require.True(t, result.Pass)
}

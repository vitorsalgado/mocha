package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvert(t *testing.T) {
	rr, err := StrictEqual(100).Match(float64(100))
	require.NoError(t, err)
	require.False(t, rr.Pass)

	result, err := ConvertTo[int](StrictEqual(100)).Match(float64(100))
	require.NoError(t, err)
	require.True(t, result.Pass)
}

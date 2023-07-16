package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnything(t *testing.T) {
	res, err := Anything().Match("")
	require.NoError(t, err)
	require.True(t, res.Pass)
}

package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsPresent(t *testing.T) {
	resYes, _ := Present().Match("test")
	resNo, _ := Present().Match("")
	require.True(t, resYes.Pass)
	require.False(t, resNo.Pass)

	resYes, _ = Present().Match(1)
	require.True(t, resYes.Pass)

	resYes, _ = Present().Match(0)
	require.True(t, resYes.Pass)

	p := "test"
	resYes, _ = Present().Match(&p)
	resNo, _ = Present().Match(nil)
	require.True(t, resYes.Pass)
	require.False(t, resNo.Pass)
}

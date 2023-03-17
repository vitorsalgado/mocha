package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsPresent(t *testing.T) {
	resYes, _ := Present().Match("test")
	resNo, _ := Present().Match("")
	assert.True(t, resYes.Pass)
	assert.False(t, resNo.Pass)

	resYes, _ = Present().Match(1)
	assert.True(t, resYes.Pass)

	resYes, _ = Present().Match(0)
	assert.True(t, resYes.Pass)

	p := "test"
	resYes, _ = Present().Match(&p)
	resNo, _ = Present().Match(nil)
	assert.True(t, resYes.Pass)
	assert.False(t, resNo.Pass)
}

func TestBePresentMatcher_Name(t *testing.T) {
	require.NotEmpty(t, Present().Name())
}

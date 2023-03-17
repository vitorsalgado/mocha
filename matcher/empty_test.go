package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsEmpty(t *testing.T) {
	resYes, err := Empty().Match("")
	assert.Nil(t, err)
	resNo, err := Empty().Match("test")
	assert.Nil(t, err)
	assert.True(t, resYes.Pass)
	assert.False(t, resNo.Pass)

	resYes, err = Empty().Match([]string{})
	assert.Nil(t, err)
	resNo, err = Empty().Match([]string{"test"})
	assert.Nil(t, err)
	assert.True(t, resYes.Pass)
	assert.False(t, resNo.Pass)

	resYes, err = Empty().Match(map[string]string{})
	assert.Nil(t, err)
	resNo, err = Empty().Match(map[string]string{"k": "v"})
	assert.Nil(t, err)
	assert.True(t, resYes.Pass)
	assert.False(t, resNo.Pass)
}

func TestEmptyMatcher_Name(t *testing.T) {
	require.NotEmpty(t, Empty().Name())
}

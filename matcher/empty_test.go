package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmpty(t *testing.T) {
	resYes, err := Empty().Match("")
	assert.Nil(t, err)
	resNo, err := Empty().Match("test")
	assert.Nil(t, err)
	assert.True(t, resYes.OK)
	assert.False(t, resNo.OK)

	resYes, err = Empty().Match([]string{})
	assert.Nil(t, err)
	resNo, err = Empty().Match([]string{"test"})
	assert.Nil(t, err)
	assert.True(t, resYes.OK)
	assert.False(t, resNo.OK)

	resYes, err = Empty().Match(map[string]string{})
	assert.Nil(t, err)
	resNo, err = Empty().Match(map[string]string{"k": "v"})
	assert.Nil(t, err)
	assert.True(t, resYes.OK)
	assert.False(t, resNo.OK)
}

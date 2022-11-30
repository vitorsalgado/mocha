package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPresent(t *testing.T) {
	resYes, _ := Present().Match("test")
	resNo, _ := Present().Match("")
	assert.True(t, resYes.OK)
	assert.False(t, resNo.OK)

	resYes, _ = Present().Match(1)
	assert.True(t, resYes.OK)

	resYes, _ = Present().Match(0)
	assert.True(t, resYes.OK)

	p := "test"
	resYes, _ = Present().Match(&p)
	resNo, _ = Present().Match(nil)
	assert.True(t, resYes.OK)
	assert.False(t, resNo.OK)
}

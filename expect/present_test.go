package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPresent(t *testing.T) {
	resYes, _ := ToBePresent().Match("test")
	resNo, _ := ToBePresent().Match("")
	assert.True(t, resYes.OK)
	assert.False(t, resNo.OK)

	resYes, _ = ToBePresent().Match(1)
	assert.True(t, resYes.OK)

	resYes, _ = ToBePresent().Match(0)
	assert.True(t, resYes.OK)

	p := "test"
	resYes, _ = ToBePresent().Match(&p)
	resNo, _ = ToBePresent().Match(nil)
	assert.True(t, resYes.OK)
	assert.False(t, resNo.OK)
}

package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmpty(t *testing.T) {
	resYes, err := ToBeEmpty().Match("")
	assert.Nil(t, err)
	resNo, err := ToBeEmpty().Match("test")
	assert.Nil(t, err)
	assert.True(t, resYes)
	assert.False(t, resNo)

	resYes, err = ToBeEmpty().Match([]string{})
	assert.Nil(t, err)
	resNo, err = ToBeEmpty().Match([]string{"test"})
	assert.Nil(t, err)
	assert.True(t, resYes)
	assert.False(t, resNo)

	resYes, err = ToBeEmpty().Match(map[string]string{})
	assert.Nil(t, err)
	resNo, err = ToBeEmpty().Match(map[string]string{"k": "v"})
	assert.Nil(t, err)
	assert.True(t, resYes)
	assert.False(t, resNo)
}

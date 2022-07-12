package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmpty(t *testing.T) {
	resYes, err := ToBeEmpty().Matches("", emptyArgs())
	assert.Nil(t, err)
	resNo, err := ToBeEmpty().Matches("test", emptyArgs())
	assert.Nil(t, err)
	assert.True(t, resYes)
	assert.False(t, resNo)

	resYes, err = ToBeEmpty().Matches([]string{}, emptyArgs())
	assert.Nil(t, err)
	resNo, err = ToBeEmpty().Matches([]string{"test"}, emptyArgs())
	assert.Nil(t, err)
	assert.True(t, resYes)
	assert.False(t, resNo)

	resYes, err = ToBeEmpty().Matches(map[string]string{}, emptyArgs())
	assert.Nil(t, err)
	resNo, err = ToBeEmpty().Matches(map[string]string{"k": "v"}, emptyArgs())
	assert.Nil(t, err)
	assert.True(t, resYes)
	assert.False(t, resNo)
}

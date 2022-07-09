package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmpty(t *testing.T) {
	resYes, err := ToBeEmpty[string]().Matches("", emptyArgs())
	assert.Nil(t, err)
	resNo, err := ToBeEmpty[string]().Matches("test", emptyArgs())
	assert.Nil(t, err)
	assert.True(t, resYes)
	assert.False(t, resNo)

	resYes, err = ToBeEmpty[[]string]().Matches([]string{}, emptyArgs())
	assert.Nil(t, err)
	resNo, err = ToBeEmpty[[]string]().Matches([]string{"test"}, emptyArgs())
	assert.Nil(t, err)
	assert.True(t, resYes)
	assert.False(t, resNo)

	resYes, err = ToBeEmpty[map[string]string]().Matches(map[string]string{}, emptyArgs())
	assert.Nil(t, err)
	resNo, err = ToBeEmpty[map[string]string]().Matches(map[string]string{"k": "v"}, emptyArgs())
	assert.Nil(t, err)
	assert.True(t, resYes)
	assert.False(t, resNo)
}

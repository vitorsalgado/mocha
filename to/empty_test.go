package to

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmpty(t *testing.T) {
	resYes, err := BeEmpty[string]().Matches("", emptyArgs())
	assert.Nil(t, err)
	resNo, err := BeEmpty[string]().Matches("test", emptyArgs())
	assert.Nil(t, err)
	assert.True(t, resYes)
	assert.False(t, resNo)

	resYes, err = BeEmpty[[]string]().Matches([]string{}, emptyArgs())
	assert.Nil(t, err)
	resNo, err = BeEmpty[[]string]().Matches([]string{"test"}, emptyArgs())
	assert.Nil(t, err)
	assert.True(t, resYes)
	assert.False(t, resNo)

	resYes, err = BeEmpty[map[string]string]().Matches(map[string]string{}, emptyArgs())
	assert.Nil(t, err)
	resNo, err = BeEmpty[map[string]string]().Matches(map[string]string{"k": "v"}, emptyArgs())
	assert.Nil(t, err)
	assert.True(t, resYes)
	assert.False(t, resNo)
}

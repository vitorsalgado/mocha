package to

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPresent(t *testing.T) {
	resYes, _ := BePresent[string]().Matches("test", emptyArgs())
	resNo, _ := BePresent[string]().Matches("", emptyArgs())
	assert.True(t, resYes)
	assert.False(t, resNo)

	resYes, _ = BePresent[int]().Matches(1, emptyArgs())
	assert.True(t, resYes)

	resYes, _ = BePresent[int]().Matches(0, emptyArgs())
	assert.True(t, resYes)

	p := "test"
	resYes, _ = BePresent[*string]().Matches(&p, emptyArgs())
	resNo, _ = BePresent[*string]().Matches(nil, emptyArgs())
	assert.True(t, resYes)
	assert.False(t, resNo)
}

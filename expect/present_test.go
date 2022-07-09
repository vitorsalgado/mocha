package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPresent(t *testing.T) {
	resYes, _ := ToBePresent[string]().Matches("test", emptyArgs())
	resNo, _ := ToBePresent[string]().Matches("", emptyArgs())
	assert.True(t, resYes)
	assert.False(t, resNo)

	resYes, _ = ToBePresent[int]().Matches(1, emptyArgs())
	assert.True(t, resYes)

	resYes, _ = ToBePresent[int]().Matches(0, emptyArgs())
	assert.True(t, resYes)

	p := "test"
	resYes, _ = ToBePresent[*string]().Matches(&p, emptyArgs())
	resNo, _ = ToBePresent[*string]().Matches(nil, emptyArgs())
	assert.True(t, resYes)
	assert.False(t, resNo)
}

package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPresent(t *testing.T) {
	resYes, _ := ToBePresent().Matches("test", emptyArgs())
	resNo, _ := ToBePresent().Matches("", emptyArgs())
	assert.True(t, resYes)
	assert.False(t, resNo)

	resYes, _ = ToBePresent().Matches(1, emptyArgs())
	assert.True(t, resYes)

	resYes, _ = ToBePresent().Matches(0, emptyArgs())
	assert.True(t, resYes)

	p := "test"
	resYes, _ = ToBePresent().Matches(&p, emptyArgs())
	resNo, _ = ToBePresent().Matches(nil, emptyArgs())
	assert.True(t, resYes)
	assert.False(t, resNo)
}

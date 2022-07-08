package to

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnyOf(t *testing.T) {
	t.Run("should return true if any of the given matchers returns true", func(t *testing.T) {
		result, err := BeAnyOf(
			Equal("test"),
			EqualFold("dev"),
			LowerCase(Equal("TEST")),
			Contains("qa")).
			Matches("test", emptyArgs())
		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return false if all of the given matchers returns false", func(t *testing.T) {
		result, err := BeAnyOf(
			Equal("abc"),
			EqualFold("def"),
			LowerCase(Equal("TEST")),
			Contains("dev")).
			Matches("test", emptyArgs())
		assert.Nil(t, err)
		assert.False(t, result)
	})
}

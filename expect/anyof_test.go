package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnyOf(t *testing.T) {
	t.Run("should return true if any of the given matchers returns true", func(t *testing.T) {
		result, err := AnyOf(
			ToEqual("test"),
			ToEqualFold("dev"),
			LowerCase(ToEqual("TEST")),
			ToContain("qa")).
			Matches("test", emptyArgs())
		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return false if all of the given matchers returns false", func(t *testing.T) {
		result, err := AnyOf(
			ToEqual("abc"),
			ToEqualFold("def"),
			LowerCase(ToEqual("TEST")),
			ToContain("dev")).
			Matches("test", emptyArgs())
		assert.Nil(t, err)
		assert.False(t, result)
	})
}

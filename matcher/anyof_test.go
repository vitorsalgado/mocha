package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnyOf(t *testing.T) {
	t.Run("should return true if any of the given matchers returns true", func(t *testing.T) {
		result, err := AnyOf(
			EqualTo("test"),
			EqualFold("dev"),
			ToLowerCase(EqualTo("TEST")),
			Contains("qa"))("test", Args{})
		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return false if all of the given matchers returns false", func(t *testing.T) {
		result, err := AnyOf(
			EqualTo("abc"),
			EqualFold("def"),
			ToLowerCase(EqualTo("TEST")),
			Contains("dev"))("test", Args{})
		assert.Nil(t, err)
		assert.False(t, result)
	})
}

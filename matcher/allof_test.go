package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllOf(t *testing.T) {
	t.Run("should return true when all matchers evaluates to true", func(t *testing.T) {
		result, err := AllOf(
			StrictEqual("test"),
			EqualIgnoreCase("TEST"),
			ToUpper(StrictEqual("TEST")),
			Contain("tes")).
			Match("test")
		assert.Nil(t, err)
		assert.True(t, result.Pass)
	})

	t.Run("should return false when just one matcher evaluates to false", func(t *testing.T) {
		result, err := AllOf(
			StrictEqual("test"),
			EqualIgnoreCase("dev"),
			ToUpper(StrictEqual("TEST")),
			Contain("tes")).
			Match("test")
		assert.Nil(t, err)
		assert.False(t, result.Pass)
	})

	t.Run("should return false when all matchers evaluates to false", func(t *testing.T) {
		result, err := AllOf(
			StrictEqual("dev"),
			EqualIgnoreCase("qa"),
			ToUpper(StrictEqual("none")),
			Contain("blah")).
			Match("test")
		assert.Nil(t, err)
		assert.False(t, result.Pass)
	})

	t.Run("mismatch description is not empty", func(t *testing.T) {
		assert.Panics(t, func() {
			AllOf()
		})
	})
}

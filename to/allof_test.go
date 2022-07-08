package to

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllOf(t *testing.T) {
	t.Parallel()

	t.Run("should return true when all matchers evaluates to true", func(t *testing.T) {
		result, err := BeAllOf(
			Equal("test"),
			EqualFold("TEST"),
			ToUpperCase(Equal("TEST")),
			Contains("tes")).
			Matches("test", emptyArgs())
		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return false when just one matcher evaluates to false", func(t *testing.T) {
		result, err := BeAllOf(
			Equal("test"),
			EqualFold("dev"),
			ToUpperCase(Equal("TEST")),
			Contains("tes")).
			Matches("test", emptyArgs())
		assert.Nil(t, err)
		assert.False(t, result)
	})

	t.Run("should return false when all matchers evaluates to false", func(t *testing.T) {
		result, err := BeAllOf(
			Equal("dev"),
			EqualFold("qa"),
			ToUpperCase(Equal("none")),
			Contains("blah")).
			Matches("test", emptyArgs())
		assert.Nil(t, err)
		assert.False(t, result)
	})
}

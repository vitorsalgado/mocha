package matchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllOf(t *testing.T) {
	t.Parallel()

	t.Run("should return true when all matchers evaluates to true", func(t *testing.T) {
		result, err := AllOf(
			EqualTo("test"),
			EqualFold("TEST"),
			ToUpperCase(EqualTo("TEST")),
			Contains("tes")).
			Matches("test", emptyArgs())
		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return false when just one matcher evaluates to false", func(t *testing.T) {
		result, err := AllOf(
			EqualTo("test"),
			EqualFold("dev"),
			ToUpperCase(EqualTo("TEST")),
			Contains("tes")).
			Matches("test", emptyArgs())
		assert.Nil(t, err)
		assert.False(t, result)
	})

	t.Run("should return false when all matchers evaluates to false", func(t *testing.T) {
		result, err := AllOf(
			EqualTo("dev"),
			EqualFold("qa"),
			ToUpperCase(EqualTo("none")),
			Contains("blah")).
			Matches("test", emptyArgs())
		assert.Nil(t, err)
		assert.False(t, result)
	})
}

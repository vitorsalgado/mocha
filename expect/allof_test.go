package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllOf(t *testing.T) {
	t.Parallel()

	t.Run("should return true when all matchers evaluates to true", func(t *testing.T) {
		result, err := AllOf(
			ToEqual("test"),
			ToEqualFold("TEST"),
			UpperCase(ToEqual("TEST")),
			ToContain("tes")).
			Matches("test", emptyArgs())
		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return false when just one matcher evaluates to false", func(t *testing.T) {
		result, err := AllOf(
			ToEqual("test"),
			ToEqualFold("dev"),
			UpperCase(ToEqual("TEST")),
			ToContain("tes")).
			Matches("test", emptyArgs())
		assert.Nil(t, err)
		assert.False(t, result)
	})

	t.Run("should return false when all matchers evaluates to false", func(t *testing.T) {
		result, err := AllOf(
			ToEqual("dev"),
			ToEqualFold("qa"),
			UpperCase(ToEqual("none")),
			ToContain("blah")).
			Matches("test", emptyArgs())
		assert.Nil(t, err)
		assert.False(t, result)
	})

	t.Run("mismatch description is not empty", func(t *testing.T) {
		assert.NotEmpty(t, AllOf().DescribeMismatch("any target", "value"))
	})
}

package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnyOf(t *testing.T) {
	t.Run("should return true if any of the given matchers returns true", func(t *testing.T) {
		result, err := AnyOf(
			Equal("test"),
			EqualIgnoreCase("dev"),
			ToLower(Equal("TEST")),
			Contain("qa")).
			Match("test")
		assert.Nil(t, err)
		assert.True(t, result.OK)
	})

	t.Run("should return false if all of the given matchers returns false", func(t *testing.T) {
		result, err := AnyOf(
			Equal("abc"),
			EqualIgnoreCase("def"),
			ToLower(Equal("TEST")),
			Contain("dev")).
			Match("test")
		assert.Nil(t, err)
		assert.False(t, result.OK)
	})

	t.Run("mismatch description is not empty", func(t *testing.T) {
		result, err := AnyOf().Match("")

		assert.NoError(t, err)
		assert.NotEmpty(t, result.DescribeFailure())
	})
}

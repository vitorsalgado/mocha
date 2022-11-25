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
			Match("test")
		assert.Nil(t, err)
		assert.True(t, result.OK)
	})

	t.Run("should return false if all of the given matchers returns false", func(t *testing.T) {
		result, err := AnyOf(
			ToEqual("abc"),
			ToEqualFold("def"),
			LowerCase(ToEqual("TEST")),
			ToContain("dev")).
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

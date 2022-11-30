package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllOf(t *testing.T) {
	t.Run("should return true when all matchers evaluates to true", func(t *testing.T) {
		result, err := AllOf(
			Equal("test"),
			EqualIgnoreCase("TEST"),
			ToUpper(Equal("TEST")),
			Contain("tes")).
			Match("test")
		assert.Nil(t, err)
		assert.True(t, result.OK)
	})

	t.Run("should return false when just one matcher evaluates to false", func(t *testing.T) {
		result, err := AllOf(
			Equal("test"),
			EqualIgnoreCase("dev"),
			ToUpper(Equal("TEST")),
			Contain("tes")).
			Match("test")
		assert.Nil(t, err)
		assert.False(t, result.OK)
	})

	t.Run("should return false when all matchers evaluates to false", func(t *testing.T) {
		result, err := AllOf(
			Equal("dev"),
			EqualIgnoreCase("qa"),
			ToUpper(Equal("none")),
			Contain("blah")).
			Match("test")
		assert.Nil(t, err)
		assert.False(t, result.OK)
	})

	t.Run("mismatch description is not empty", func(t *testing.T) {
		result, err := AllOf().Match("")
		assert.NoError(t, err)
		assert.NotEmpty(t, result.DescribeFailure())
	})
}

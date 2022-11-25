package expect

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegExpMatches(t *testing.T) {
	t.Run("should match the regular expression string pattern", func(t *testing.T) {
		result, err := ToMatchExpr("tEsT").Match("tEsT")

		assert.Nil(t, err)
		assert.True(t, result.OK)
	})

	t.Run("should match the regular expression string pattern using a non string argument", func(t *testing.T) {
		result, err := ToMatchExpr("10").Match(10)

		assert.NoError(t, err)
		assert.True(t, result.OK)
	})

	t.Run("should match the provided regular expression against matcher argument", func(t *testing.T) {
		re := regexp.MustCompile("tEsT")
		result, err := ToMatchExpr(re).Match("tEsT")

		assert.Nil(t, err)
		assert.True(t, result.OK)
	})

	t.Run("should accept a non pointer regular expression", func(t *testing.T) {
		re := regexp.MustCompile("tEsT")
		result, err := ToMatchExpr(*re).Match("tEsT")

		assert.Nil(t, err)
		assert.True(t, result.OK)
	})

	t.Run("should return false when regexp does not match", func(t *testing.T) {
		re := regexp.MustCompile("tEsT")
		result, err := ToMatchExpr(re).Match("dev")

		assert.Nil(t, err)
		assert.False(t, result.OK)
	})
}

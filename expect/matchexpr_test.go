package expect

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegExpMatches(t *testing.T) {
	params := Args{}

	t.Run("should match the regular expression string pattern", func(t *testing.T) {
		result, err := ToMatchExpr("tEsT").Matches("tEsT", params)

		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should match the regular expression string pattern using a non string argument", func(t *testing.T) {
		result, err := ToMatchExpr("10").Matches(10, params)

		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should match the provided regular expression against matcher argument", func(t *testing.T) {
		re := regexp.MustCompile("tEsT")
		result, err := ToMatchExpr(re).Matches("tEsT", params)

		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should accept a non pointer regular expression", func(t *testing.T) {
		re := regexp.MustCompile("tEsT")
		result, err := ToMatchExpr(*re).Matches("tEsT", params)

		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return false when regexp does not match", func(t *testing.T) {
		re := regexp.MustCompile("tEsT")
		result, err := ToMatchExpr(re).Matches("dev", params)

		assert.Nil(t, err)
		assert.False(t, result)
	})
}

package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasPrefix(t *testing.T) {
	t.Run("should return true when string has prefix", func(t *testing.T) {
		result, err := ToHavePrefix("hello").Match("hello world")

		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return true when string hasn't prefix", func(t *testing.T) {
		result, err := ToHavePrefix("world").Match("hello world")

		assert.Nil(t, err)
		assert.False(t, result)
	})
}

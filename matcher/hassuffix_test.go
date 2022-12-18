package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasSuffix(t *testing.T) {
	t.Run("should return true when string has suffix", func(t *testing.T) {
		result, err := HasSuffix("world").Match("hello world")

		assert.Nil(t, err)
		assert.True(t, result.Pass)
	})

	t.Run("should return true when string hasn't suffix", func(t *testing.T) {
		result, err := HasSuffix("hello").Match("hello world")

		assert.Nil(t, err)
		assert.False(t, result.Pass)
	})
}

package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasPrefix(t *testing.T) {
	t.Run("should return true when string has prefix", func(t *testing.T) {
		result, err := HasPrefix("hello")("hello world", Params{})

		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return true when string hasn't prefix", func(t *testing.T) {
		result, err := HasPrefix("world")("hello world", Params{})

		assert.Nil(t, err)
		assert.False(t, result)
	})
}

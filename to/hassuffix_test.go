package to

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasSuffix(t *testing.T) {
	t.Run("should return true when string has suffix", func(t *testing.T) {
		result, err := HaveSuffix("world").Matches("hello world", emptyArgs())

		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return true when string hasn't suffix", func(t *testing.T) {
		result, err := HaveSuffix("hello").Matches("hello world", emptyArgs())

		assert.Nil(t, err)
		assert.False(t, result)
	})
}

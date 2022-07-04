package matchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	params := Args{}

	t.Run("should return true when value is contained in the matcher argument", func(t *testing.T) {
		result, err := Contains("world")("hello world", params)
		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return false when value is not contained in the matcher argument", func(t *testing.T) {
		result, err := Contains("dev")("hello world", params)
		assert.Nil(t, err)
		assert.False(t, result)
	})
}

package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	t.Run("should return true when value is contained in the matcher argument", func(t *testing.T) {
		result, err := ToContain[string]("world").Matches("hello world", emptyArgs())
		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return false when value is not contained in the matcher argument", func(t *testing.T) {
		result, err := ToContain[string]("dev").Matches("hello world", emptyArgs())
		assert.Nil(t, err)
		assert.False(t, result)
	})

	t.Run("should ", func(t *testing.T) {
		result, err := ToContain[[]string]("dev").Matches([]string{"dev", "qa"}, emptyArgs())
		assert.Nil(t, err)
		assert.True(t, result)
	})
}

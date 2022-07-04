package matchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsItem(t *testing.T) {
	t.Run("should return true when value is contained inside matcher array argument", func(t *testing.T) {
		result, err := ContainsItem("apple").Matches([]string{"banana", "apple", "orange"}, emptyArgs())
		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return false when value is contained inside matcher array argument", func(t *testing.T) {
		result, err := ContainsItem(5).Matches([]int{1, 2, 6, 7}, emptyArgs())
		assert.Nil(t, err)
		assert.False(t, result)
	})
}

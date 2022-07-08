package to

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsItem(t *testing.T) {
	t.Run("should return true when value is contained inside matcher array argument", func(t *testing.T) {
		result, err := ContainItem("apple").Matches([]string{"banana", "apple", "orange"}, emptyArgs())
		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return false when value is contained inside matcher array argument", func(t *testing.T) {
		result, err := ContainItem(5).Matches([]int{1, 2, 6, 7}, emptyArgs())
		assert.Nil(t, err)
		assert.False(t, result)
	})
}

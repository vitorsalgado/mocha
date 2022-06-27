package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsItem(t *testing.T) {
	params := Params{}

	t.Run("should return true when value is contained inside matcher array argument", func(t *testing.T) {
		result, err := ContainsItem("apple")([]string{"banana", "apple", "orange"}, params)
		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return false when value is contained inside matcher array argument", func(t *testing.T) {
		result, err := ContainsItem(5)([]int{1, 2, 6, 7}, params)
		assert.Nil(t, err)
		assert.False(t, result)
	})
}

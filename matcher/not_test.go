package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNot(t *testing.T) {
	value := "test"
	params := Args{}

	t.Run("should return true when value is not equal", func(t *testing.T) {
		result, err := Not(EqualTo("dev"))(value, params)

		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return false when value is equal", func(t *testing.T) {
		result, err := Not(EqualTo("test"))(value, params)

		assert.Nil(t, err)
		assert.False(t, result)
	})
}

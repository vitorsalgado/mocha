package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNot(t *testing.T) {
	value := "test"

	t.Run("should return true when value is not equal", func(t *testing.T) {
		result, err := Not(StrictEqual("dev")).Match(value)

		assert.Nil(t, err)
		assert.True(t, result.Pass)
	})

	t.Run("should return false when value is equal", func(t *testing.T) {
		result, err := Not(StrictEqual("test")).Match(value)

		assert.Nil(t, err)
		assert.False(t, result.Pass)
	})
}

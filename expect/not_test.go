package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNot(t *testing.T) {
	value := "test"

	t.Run("should return true when value is not equal", func(t *testing.T) {
		result, err := Not(ToEqual("dev")).Matches(value, emptyArgs())

		assert.Nil(t, err)
		assert.True(t, result)
	})

	t.Run("should return false when value is equal", func(t *testing.T) {
		result, err := Not(ToEqual("test")).Matches(value, emptyArgs())

		assert.Nil(t, err)
		assert.False(t, result)
	})
}

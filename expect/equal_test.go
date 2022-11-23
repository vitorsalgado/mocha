package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqual(t *testing.T) {
	t.Parallel()

	t.Run("should compare expected string with nil value", func(t *testing.T) {
		exp := "test"
		res, err := ToEqual(&exp).Match(nil)

		assert.Nil(t, err)
		assert.False(t, res)
	})

	t.Run("should compare two byte arrays", func(t *testing.T) {
		value := []byte("test")
		res, err := ToEqual(value).Match(value)

		assert.Nil(t, err)
		assert.True(t, res)
	})
}

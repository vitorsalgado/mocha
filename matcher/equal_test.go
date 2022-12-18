package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqual(t *testing.T) {
	t.Run("should compare items string with nil value", func(t *testing.T) {
		exp := "test"
		res, err := Equal(&exp).Match(nil)

		assert.Nil(t, err)
		assert.False(t, res.Pass)
	})

	t.Run("should compare two byte arrays", func(t *testing.T) {
		value := []byte("test")
		res, err := Equal(value).Match(value)

		assert.Nil(t, err)
		assert.True(t, res.Pass)
	})
}

package mocha

import (
	"github.com/vitorsalgado/mocha/internal/assert"
	"testing"
)

func TestEqual(t *testing.T) {
	t.Run("should compare expected string with nil value", func(t *testing.T) {
		exp := "test"
		res, err := Equal(&exp)(nil, MatcherContext{})

		assert.Nil(t, err)
		assert.False(t, res)
	})

	t.Run("should compare two byte arrays", func(t *testing.T) {
		value := []byte("test")
		res, err := Equal(value)(value, MatcherContext{})

		assert.Nil(t, err)
		assert.True(t, res)
	})
}

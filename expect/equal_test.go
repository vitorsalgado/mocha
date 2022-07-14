package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqual(t *testing.T) {
	t.Parallel()

	t.Run("should compare expected string with nil value", func(t *testing.T) {
		exp := "test"
		res, err := ToEqual(&exp).Matches(nil, emptyArgs())

		assert.Nil(t, err)
		assert.False(t, res)
	})

	t.Run("should compare two byte arrays", func(t *testing.T) {
		value := []byte("test")
		res, err := ToEqual(value).Matches(value, emptyArgs())

		assert.Nil(t, err)
		assert.True(t, res)
	})
}

func TestToEqualJSON(t *testing.T) {
	t.Run("should return matcher error", func(t *testing.T) {
		c := make(chan bool, 1)
		body := map[string]interface{}{"ok": true, "name": "dev"}
		res, err := ToEqualJSON(c).Matches(body, emptyArgs())

		assert.Error(t, err)
		assert.False(t, res)
	})
}

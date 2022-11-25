package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONPath(t *testing.T) {
	m := map[string]any{
		"name": "someone",
		"age":  34,
		"address": map[string]any{
			"street": "very nice place",
		},
	}

	t.Run("should read match text field on object root", func(t *testing.T) {
		res, err := JSONPath("name", ToEqual("someone")).Match(m)
		assert.Nil(t, err)
		assert.True(t, res.OK)
	})

	t.Run("should match numeric field value", func(t *testing.T) {
		res, err := JSONPath("age", ToEqual(34)).Match(m)
		assert.Nil(t, err)
		assert.True(t, res.OK)
	})

	t.Run("should match nested object field", func(t *testing.T) {
		res, err := JSONPath("address.street", ToEqual("very nice place")).Match(m)
		assert.Nil(t, err)
		assert.True(t, res.OK)
	})

	t.Run("should return error when any error occurs", func(t *testing.T) {
		res, err := JSONPath("312nj.,", ToEqual("anything")).Match(m)
		assert.NotNil(t, err)
		assert.False(t, res.OK)
	})
}

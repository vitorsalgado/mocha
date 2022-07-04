package matchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONPath(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"name": "someone",
		"age":  34,
		"address": map[string]any{
			"street": "very nice place",
		},
	}

	t.Run("should read match text field on object root", func(t *testing.T) {
		res, err := JSONPath("name", EqualAny("someone"))(m, Args{})
		assert.Nil(t, err)
		assert.True(t, res)
	})

	t.Run("should match numeric field value", func(t *testing.T) {
		res, err := JSONPath("age", EqualTo(34))(m, Args{})
		assert.Nil(t, err)
		assert.True(t, res)
	})

	t.Run("should match nested object field", func(t *testing.T) {
		res, err := JSONPath("address.street", EqualTo("very nice place"))(m, Args{})
		assert.Nil(t, err)
		assert.True(t, res)
	})

	t.Run("should return error when any error occurs", func(t *testing.T) {
		res, err := JSONPath("312nj.,", EqualAny("anything"))(m, Args{})
		assert.NotNil(t, err)
		assert.False(t, res)
	})
}

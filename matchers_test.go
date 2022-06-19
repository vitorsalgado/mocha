package mocha

import (
	"testing"

	"github.com/vitorsalgado/mocha/internal/assert"
)

func TestEqual(t *testing.T) {
	t.Parallel()

	t.Run("should compare expected string with nil value", func(t *testing.T) {
		exp := "test"
		res, err := EqualTo(&exp)(nil, MatcherParams{})

		assert.Nil(t, err)
		assert.False(t, res)
	})

	t.Run("should compare two byte arrays", func(t *testing.T) {
		value := []byte("test")
		res, err := EqualTo(value)(value, MatcherParams{})

		assert.Nil(t, err)
		assert.True(t, res)
	})
}

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
		res, err := JSONPath("name", Equal("someone"))(m, MatcherParams{})
		assert.Nil(t, err)
		assert.True(t, res)
	})

	t.Run("should match numeric field value", func(t *testing.T) {
		res, err := JSONPath("age", EqualTo(34))(m, MatcherParams{})
		assert.Nil(t, err)
		assert.True(t, res)
	})

	t.Run("should match nested object field", func(t *testing.T) {
		res, err := JSONPath("address.street", EqualTo("very nice place"))(m, MatcherParams{})
		assert.Nil(t, err)
		assert.True(t, res)
	})

	t.Run("should return error when any error occurs", func(t *testing.T) {
		res, err := JSONPath("312nj.,", Equal("anything"))(m, MatcherParams{})
		assert.NotNil(t, err)
		assert.False(t, res)
	})
}

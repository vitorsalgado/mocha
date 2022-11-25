package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasKey(t *testing.T) {
	m := map[string]any{
		"name": "someone",
		"age":  34,
		"address": map[string]any{
			"street": "very nice place",
		},
		"options": []string{},
		"active":  false,
		"none":    nil,
		"zero":    0,
	}

	result, err := ToHaveKey("name").Match(m)
	assert.True(t, result.OK)
	assert.Nil(t, err)

	result, err = ToHaveKey("age").Match(m)
	assert.True(t, result.OK)
	assert.Nil(t, err)

	result, err = ToHaveKey("address").Match(m)
	assert.True(t, result.OK)
	assert.Nil(t, err)

	result, err = ToHaveKey("active").Match(m)
	assert.True(t, result.OK)
	assert.Nil(t, err)

	result, err = ToHaveKey("zero").Match(m)
	assert.True(t, result.OK)
	assert.Nil(t, err)

	result, err = ToHaveKey("options").Match(m)
	assert.True(t, result.OK)
	assert.Nil(t, err)

	result, err = ToHaveKey("none").Match(m)
	assert.False(t, result.OK)
	assert.Nil(t, err)
}

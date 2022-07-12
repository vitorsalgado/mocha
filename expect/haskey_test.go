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

	result, err := ToHaveKey("name").Matches(m, emptyArgs())
	assert.True(t, result)
	assert.Nil(t, err)

	result, err = ToHaveKey("age").Matches(m, emptyArgs())
	assert.True(t, result)
	assert.Nil(t, err)

	result, err = ToHaveKey("address").Matches(m, emptyArgs())
	assert.True(t, result)
	assert.Nil(t, err)

	result, err = ToHaveKey("active").Matches(m, emptyArgs())
	assert.True(t, result)
	assert.Nil(t, err)

	result, err = ToHaveKey("zero").Matches(m, emptyArgs())
	assert.True(t, result)
	assert.Nil(t, err)

	result, err = ToHaveKey("options").Matches(m, emptyArgs())
	assert.True(t, result)
	assert.Nil(t, err)

	result, err = ToHaveKey("none").Matches(m, emptyArgs())
	assert.False(t, result)
	assert.Nil(t, err)
}

package matcher

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

	result, err := HaveKey("name").Match(m)
	assert.True(t, result.Pass)
	assert.Nil(t, err)

	result, err = HaveKey("age").Match(m)
	assert.True(t, result.Pass)
	assert.Nil(t, err)

	result, err = HaveKey("address").Match(m)
	assert.True(t, result.Pass)
	assert.Nil(t, err)

	result, err = HaveKey("active").Match(m)
	assert.True(t, result.Pass)
	assert.Nil(t, err)

	result, err = HaveKey("zero").Match(m)
	assert.True(t, result.Pass)
	assert.Nil(t, err)

	result, err = HaveKey("options").Match(m)
	assert.True(t, result.Pass)
	assert.Nil(t, err)

	result, err = HaveKey("none").Match(m)
	assert.False(t, result.Pass)
	assert.Nil(t, err)
}

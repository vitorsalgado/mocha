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

	params := Params{}

	result, err := HasKey[any]("name")(m, params)
	assert.True(t, result)
	assert.Nil(t, err)

	result, err = HasKey[any]("age")(m, params)
	assert.True(t, result)
	assert.Nil(t, err)

	result, err = HasKey[any]("address")(m, params)
	assert.True(t, result)
	assert.Nil(t, err)

	result, err = HasKey[any]("active")(m, params)
	assert.True(t, result)
	assert.Nil(t, err)

	result, err = HasKey[any]("zero")(m, params)
	assert.True(t, result)
	assert.Nil(t, err)

	result, err = HasKey[any]("options")(m, params)
	assert.True(t, result)
	assert.Nil(t, err)

	result, err = HasKey[any]("none")(m, params)
	assert.False(t, result)
	assert.Nil(t, err)
}

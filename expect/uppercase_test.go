package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToUpperCase(t *testing.T) {
	result, err := UpperCase(ToEqual("TEST")).Matches("tesT", Args{})

	assert.Nil(t, err)
	assert.True(t, result)
}

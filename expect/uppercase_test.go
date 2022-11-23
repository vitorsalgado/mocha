package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToUpperCase(t *testing.T) {
	result, err := UpperCase(ToEqual("TEST")).Match("tesT")

	assert.Nil(t, err)
	assert.True(t, result)
}

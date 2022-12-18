package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToUpperCase(t *testing.T) {
	result, err := ToUpper(Equal("TEST")).Match("tesT")

	assert.Nil(t, err)
	assert.True(t, result.Pass)
}

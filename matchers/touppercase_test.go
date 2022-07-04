package matchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToUpperCase(t *testing.T) {
	result, err := ToUpperCase(EqualTo("TEST")).Matches("tesT", Args{})

	assert.Nil(t, err)
	assert.True(t, result)
}

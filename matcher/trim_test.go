package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrim(t *testing.T) {
	result, err := Trim(EqualTo("test"))("  test  ", Args{})

	assert.Nil(t, err)
	assert.True(t, result)
}

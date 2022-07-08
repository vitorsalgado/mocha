package to

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrim(t *testing.T) {
	result, err := Trim(Equal("test")).Matches("  test  ", Args{})

	assert.Nil(t, err)
	assert.True(t, result)
}

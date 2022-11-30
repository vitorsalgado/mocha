package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrim(t *testing.T) {
	result, err := Trim(Equal("test")).Match("  test  ")

	assert.Nil(t, err)
	assert.True(t, result.OK)
}

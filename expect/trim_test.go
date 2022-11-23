package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrim(t *testing.T) {
	result, err := Trim(ToEqual("test")).Match("  test  ")

	assert.Nil(t, err)
	assert.True(t, result)
}

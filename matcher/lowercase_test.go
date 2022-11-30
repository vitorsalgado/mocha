package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToLowerCase(t *testing.T) {
	result, err := ToLower(Equal("test")).Match("TeST")

	assert.Nil(t, err)
	assert.True(t, result.OK)
}

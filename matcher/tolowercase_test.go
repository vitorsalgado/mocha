package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToLowerCase(t *testing.T) {
	result, err := ToLowerCase(EqualTo("test"))("TeST", Args{})

	assert.Nil(t, err)
	assert.True(t, result)
}

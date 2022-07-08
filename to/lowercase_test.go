package to

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToLowerCase(t *testing.T) {
	result, err := ToLowerCase(Equal("test")).Matches("TeST", Args{})

	assert.Nil(t, err)
	assert.True(t, result)
}

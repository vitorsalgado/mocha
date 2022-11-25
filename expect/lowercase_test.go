package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToLowerCase(t *testing.T) {
	result, err := LowerCase(ToEqual("test")).Match("TeST")

	assert.Nil(t, err)
	assert.True(t, result.OK)
}

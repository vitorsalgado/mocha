package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToLowerCase(t *testing.T) {
	result, err := ToLower(StrictEqual("test")).Match("TeST")

	assert.Nil(t, err)
	assert.True(t, result.Pass)
}

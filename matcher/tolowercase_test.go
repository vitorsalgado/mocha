package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToLowerCase(t *testing.T) {
	result, err := ToLowerCase(EqualTo("test"))("TeST", Params{})

	assert.Nil(t, err)
	assert.True(t, result)
}

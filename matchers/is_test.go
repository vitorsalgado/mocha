package matchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIs(t *testing.T) {
	result, err := Is(EqualTo("test"))("test", Args{})

	assert.Nil(t, err)
	assert.True(t, result)
}

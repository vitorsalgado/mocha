package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIs(t *testing.T) {
	result, err := Is(EqualTo("test"))("test", Params{})

	assert.Nil(t, err)
	assert.True(t, result)
}

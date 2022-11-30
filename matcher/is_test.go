package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIs(t *testing.T) {
	result, err := Is(Equal("test")).Match("test")

	assert.NoError(t, err)
	assert.True(t, result.OK)
}

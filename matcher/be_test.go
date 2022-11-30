package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBe(t *testing.T) {
	result, err := Be(Equal("test")).Match("test")

	assert.NoError(t, err)
	assert.True(t, result.OK)
}

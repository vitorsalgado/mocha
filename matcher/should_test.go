package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldBe(t *testing.T) {
	result, err := Should(Equal("test")).Match("test")

	assert.NoError(t, err)
	assert.True(t, result.OK)
}

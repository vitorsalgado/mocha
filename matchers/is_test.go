package matchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIs(t *testing.T) {
	result, err := Is(EqualTo("test")).Matches("test", emptyArgs())

	assert.Nil(t, err)
	assert.True(t, result)
}

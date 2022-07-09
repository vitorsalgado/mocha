package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTo(t *testing.T) {
	result, err := ToBe(ToEqual("test")).Matches("test", emptyArgs())

	assert.Nil(t, err)
	assert.True(t, result)
}

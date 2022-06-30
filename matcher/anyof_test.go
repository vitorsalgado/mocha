package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnyOf(t *testing.T) {
	result, err := AnyOf(
		EqualTo("test"),
		EqualFold("dev"),
		ToLowerCase(EqualTo("TEST")),
		Contains("qa"))("test", Args{})
	assert.Nil(t, err)
	assert.True(t, result)
}

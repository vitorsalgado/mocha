package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrim(t *testing.T) {
	result, err := Trim(StrictEqual("test")).Match("  test  ")

	assert.Nil(t, err)
	assert.True(t, result.Pass)
}

func TestTrimMatcher_Name(t *testing.T) {
	require.NotEmpty(t, Trim(Eq("")).Name())
}

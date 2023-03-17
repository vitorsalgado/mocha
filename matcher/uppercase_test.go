package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToUpperCase(t *testing.T) {
	result, err := ToUpper(StrictEqual("TEST")).Match("tesT")

	assert.Nil(t, err)
	assert.True(t, result.Pass)
}

func TestUpperCaseMatcher_Name(t *testing.T) {
	require.NotEmpty(t, ToUpper(Eq("")).Name())
}

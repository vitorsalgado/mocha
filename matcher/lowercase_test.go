package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToLowerCase(t *testing.T) {
	result, err := ToLower(StrictEqual("test")).Match("TeST")

	assert.Nil(t, err)
	assert.True(t, result.Pass)
}

func TestLowerCaseMatcher_Name(t *testing.T) {
	require.NotEmpty(t, ToLower(Eq("")).Name())
}

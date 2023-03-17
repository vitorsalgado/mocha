package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplit(t *testing.T) {
	v := "test, testing, tested"
	r, err := Split(",", Each(Trim(HasPrefix("test")))).Match(v)

	assert.NoError(t, err)
	assert.True(t, r.Pass)
}

func TestSplitMatcher_Name(t *testing.T) {
	require.NotEmpty(t, Split(".", Eq("")).Name())
}

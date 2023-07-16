package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplit(t *testing.T) {
	v := "test, testing, tested"
	r, err := Split(",", Each(Trim(HasPrefix("test")))).Match(v)

	require.NoError(t, err)
	require.True(t, r.Pass)
}

package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	v := "test, testing, tested"
	r, err := Split(",", Each(Trim(HasPrefix("test")))).Match(v)

	assert.NoError(t, err)
	assert.True(t, r.Pass)
}

package matcher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAnythingMatcher(t *testing.T) {
	res, err := Anything[any]()(nil, Params{})
	assert.Nil(t, err)
	assert.True(t, res)
}

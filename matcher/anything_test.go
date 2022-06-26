package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnythingMatcher(t *testing.T) {
	res, err := Anything[any]()(nil, Params{})
	assert.Nil(t, err)
	assert.True(t, res)
}

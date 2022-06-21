package matcher

import (
	"testing"

	"github.com/vitorsalgado/mocha/internal/assert"
)

func TestAnythingMatcher(t *testing.T) {
	res, err := Anything[any]()(nil, Params{})
	assert.Nil(t, err)
	assert.True(t, res)
}

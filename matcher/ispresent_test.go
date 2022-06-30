package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPresent(t *testing.T) {
	params := Args{}

	resy, _ := IsPresent[string]()("test", params)
	resn, _ := IsPresent[string]()("", params)
	assert.True(t, resy)
	assert.False(t, resn)

	resy, _ = IsPresent[int]()(1, params)
	assert.True(t, resy)

	resy, _ = IsPresent[int]()(0, params)
	assert.True(t, resy)

	p := "test"
	resy, _ = IsPresent[*string]()(&p, params)
	resn, _ = IsPresent[*string]()(nil, params)
	assert.True(t, resy)
	assert.False(t, resn)
}

package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmpty(t *testing.T) {
	params := Args{}

	resy, err := IsEmpty[string]()("", params)
	assert.Nil(t, err)
	resn, err := IsEmpty[string]()("test", params)
	assert.Nil(t, err)
	assert.True(t, resy)
	assert.False(t, resn)

	resy, err = IsEmpty[[]string]()([]string{}, params)
	assert.Nil(t, err)
	resn, err = IsEmpty[[]string]()([]string{"test"}, params)
	assert.Nil(t, err)
	assert.True(t, resy)
	assert.False(t, resn)

	resy, err = IsEmpty[map[string]string]()(map[string]string{}, params)
	assert.Nil(t, err)
	resn, err = IsEmpty[map[string]string]()(map[string]string{"k": "v"}, params)
	assert.Nil(t, err)
	assert.True(t, resy)
	assert.False(t, resn)
}

package matchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmpty(t *testing.T) {
	resy, err := IsEmpty[string]().Matches("", emptyArgs())
	assert.Nil(t, err)
	resn, err := IsEmpty[string]().Matches("test", emptyArgs())
	assert.Nil(t, err)
	assert.True(t, resy)
	assert.False(t, resn)

	resy, err = IsEmpty[[]string]().Matches([]string{}, emptyArgs())
	assert.Nil(t, err)
	resn, err = IsEmpty[[]string]().Matches([]string{"test"}, emptyArgs())
	assert.Nil(t, err)
	assert.True(t, resy)
	assert.False(t, resn)

	resy, err = IsEmpty[map[string]string]().Matches(map[string]string{}, emptyArgs())
	assert.Nil(t, err)
	resn, err = IsEmpty[map[string]string]().Matches(map[string]string{"k": "v"}, emptyArgs())
	assert.Nil(t, err)
	assert.True(t, resy)
	assert.False(t, resn)
}

package to

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmpty(t *testing.T) {
	resy, err := BeEmpty[string]().Matches("", emptyArgs())
	assert.Nil(t, err)
	resn, err := BeEmpty[string]().Matches("test", emptyArgs())
	assert.Nil(t, err)
	assert.True(t, resy)
	assert.False(t, resn)

	resy, err = BeEmpty[[]string]().Matches([]string{}, emptyArgs())
	assert.Nil(t, err)
	resn, err = BeEmpty[[]string]().Matches([]string{"test"}, emptyArgs())
	assert.Nil(t, err)
	assert.True(t, resy)
	assert.False(t, resn)

	resy, err = BeEmpty[map[string]string]().Matches(map[string]string{}, emptyArgs())
	assert.Nil(t, err)
	resn, err = BeEmpty[map[string]string]().Matches(map[string]string{"k": "v"}, emptyArgs())
	assert.Nil(t, err)
	assert.True(t, resy)
	assert.False(t, resn)
}

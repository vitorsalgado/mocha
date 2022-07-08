package to

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPresent(t *testing.T) {
	resy, _ := BePresent[string]().Matches("test", emptyArgs())
	resn, _ := BePresent[string]().Matches("", emptyArgs())
	assert.True(t, resy)
	assert.False(t, resn)

	resy, _ = BePresent[int]().Matches(1, emptyArgs())
	assert.True(t, resy)

	resy, _ = BePresent[int]().Matches(0, emptyArgs())
	assert.True(t, resy)

	p := "test"
	resy, _ = BePresent[*string]().Matches(&p, emptyArgs())
	resn, _ = BePresent[*string]().Matches(nil, emptyArgs())
	assert.True(t, resy)
	assert.False(t, resn)
}

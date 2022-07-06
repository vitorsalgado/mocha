package matchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPresent(t *testing.T) {
	resy, _ := IsPresent[string]().Matches("test", emptyArgs())
	resn, _ := IsPresent[string]().Matches("", emptyArgs())
	assert.True(t, resy)
	assert.False(t, resn)

	resy, _ = IsPresent[int]().Matches(1, emptyArgs())
	assert.True(t, resy)

	resy, _ = IsPresent[int]().Matches(0, emptyArgs())
	assert.True(t, resy)

	p := "test"
	resy, _ = IsPresent[*string]().Matches(&p, emptyArgs())
	resn, _ = IsPresent[*string]().Matches(nil, emptyArgs())
	assert.True(t, resy)
	assert.False(t, resn)
}

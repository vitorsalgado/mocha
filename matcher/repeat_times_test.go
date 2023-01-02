package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepeat(t *testing.T) {
	rep := Repeat(2)

	res, err := rep.Match(nil)
	assert.NoError(t, rep.After())

	assert.NoError(t, err)
	assert.True(t, res.Pass)

	res, err = rep.Match(nil)
	assert.NoError(t, rep.After())

	assert.NoError(t, err)
	assert.True(t, res.Pass)

	res, err = rep.Match(nil)
	assert.NoError(t, rep.After())

	assert.NoError(t, err)
	assert.False(t, res.Pass)

	res, err = rep.Match(nil)
	assert.NoError(t, rep.After())

	assert.NoError(t, err)
	assert.False(t, res.Pass)
}

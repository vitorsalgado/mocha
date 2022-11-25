package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepeat(t *testing.T) {
	rep := Repeat(2)

	res, err := rep.Match(nil)
	assert.NoError(t, rep.OnMockServed())

	assert.NoError(t, err)
	assert.True(t, res.OK)

	res, err = rep.Match(nil)
	assert.NoError(t, rep.OnMockServed())

	assert.NoError(t, err)
	assert.True(t, res.OK)

	res, err = rep.Match(nil)
	assert.NoError(t, rep.OnMockServed())

	assert.NoError(t, err)
	assert.False(t, res.OK)

	res, err = rep.Match(nil)
	assert.NoError(t, rep.OnMockServed())

	assert.NoError(t, err)
	assert.False(t, res.OK)
}

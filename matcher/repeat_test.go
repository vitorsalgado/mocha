package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepeat(t *testing.T) {
	m := Repeat[any](5)
	result, _ := m(nil, Args{})
	assert.True(t, result)
	result, _ = m(nil, Args{})
	assert.True(t, result)
	result, _ = m(nil, Args{})
	assert.True(t, result)
	result, _ = m(nil, Args{})
	assert.True(t, result)
	result, _ = m(nil, Args{})
	assert.True(t, result)

	result, _ = m(nil, Args{})
	assert.False(t, result)
	result, _ = m(nil, Args{})
	assert.False(t, result)
}

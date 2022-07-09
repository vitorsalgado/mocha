package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepeat(t *testing.T) {
	m := Repeat(5)
	result, _ := m.Matches(nil, Args{})
	assert.True(t, result)
	result, _ = m.Matches(nil, Args{})
	assert.True(t, result)
	result, _ = m.Matches(nil, Args{})
	assert.True(t, result)
	result, _ = m.Matches(nil, Args{})
	assert.True(t, result)
	result, _ = m.Matches(nil, Args{})
	assert.True(t, result)

	result, _ = m.Matches(nil, Args{})
	assert.False(t, result)
	result, _ = m.Matches(nil, Args{})
	assert.False(t, result)
}

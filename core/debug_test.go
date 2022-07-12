package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/expect"
)

func TestDebug(t *testing.T) {
	mk := NewMock()
	mk.Name = "test"
	result, err := debug(mk, expect.ToEqual("test")).Matches("test", expect.Args{})

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestDebugErr(t *testing.T) {
	mk := NewMock()
	mk.Name = "test"
	result, err := debug(mk, expect.Func(
		func(v any, params expect.Args) (bool, error) {
			return false, fmt.Errorf("failed")
		})).
		Matches("test", expect.Args{})

	assert.NotNil(t, err)
	assert.False(t, result)
}

func TestDebugNotMatched(t *testing.T) {
	mk := NewMock()
	mk.Name = "test"
	result, err := debug(mk, expect.ToEqual("test")).Matches("dev", expect.Args{})

	assert.Nil(t, err)
	assert.False(t, result)
}

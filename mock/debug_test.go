package mock

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/to"
)

func TestDebug(t *testing.T) {
	mk := New()
	mk.Name = "test"
	result, err := Debug("equalTo", *mk, to.Equal("test")).Matches("test", to.Args{})

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestDebugErr(t *testing.T) {
	mk := New()
	mk.Name = "test"
	result, err := Debug("err", *mk, to.Fn(
		func(v string, params to.Args) (bool, error) {
			return false, fmt.Errorf("failed")
		})).
		Matches("test", to.Args{})

	assert.NotNil(t, err)
	assert.False(t, result)
}

func TestDebugNotMatched(t *testing.T) {
	mk := New()
	mk.Name = "test"
	result, err := Debug("equalTo", *mk, to.Equal("test")).Matches("dev", to.Args{})

	assert.Nil(t, err)
	assert.False(t, result)
}

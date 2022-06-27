package mock

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorsalgado/mocha/matcher"
)

func TestDebug(t *testing.T) {
	mk := New()
	mk.Name = "test"
	result, err := Debug("equalTo", *mk, matcher.EqualTo("test"))("test", matcher.Params{})

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestDebugErr(t *testing.T) {
	mk := New()
	mk.Name = "test"
	result, err := Debug("err", *mk, func(v string, params matcher.Params) (bool, error) { return false, fmt.Errorf("failed") })("test", matcher.Params{})

	assert.NotNil(t, err)
	assert.False(t, result)
}

func TestDebugNotMatched(t *testing.T) {
	mk := New()
	mk.Name = "test"
	result, err := Debug("equalTo", *mk, matcher.EqualTo("test"))("dev", matcher.Params{})

	assert.Nil(t, err)
	assert.False(t, result)
}

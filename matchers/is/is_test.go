package is

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorsalgado/mocha/matchers"
)

func TestIs(t *testing.T) {
	var res = false
	var err error = nil

	res, err = AllOf(EqualTo("TEST"), EqualFold("test")).Matches("TEST", matchers.Args{})
	assert.Nil(t, err)
	assert.True(t, res)

	res, err = AnyOf(EqualTo("test"), EqualFold("dev")).Matches("test", matchers.Args{})
	assert.Nil(t, err)
	assert.True(t, res)

	res, err = Both(EqualTo("dev-test")).And(Not(Empty[string]())).Matches("dev-test", matchers.Args{})
	assert.Nil(t, err)
	assert.True(t, res)

	res, err = Either(EqualTo("dev")).Or(EqualTo("test")).Matches("dev", matchers.Args{})
	assert.Nil(t, err)
	assert.True(t, res)

	res, err = XOR(EqualTo("test"), EqualTo("dev")).Matches("test", matchers.Args{})
	assert.Nil(t, err)
	assert.True(t, res)
}

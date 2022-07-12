package is

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/expect"
)

func TestIs(t *testing.T) {
	var res = false
	var err error = nil

	res, err = AllOf(EqualTo("TEST"), EqualFold("test"), Present()).Matches("TEST", expect.Args{})
	assert.Nil(t, err)
	assert.True(t, res)

	res, err = AnyOf(EqualTo("test"), EqualFold("dev")).Matches("test", expect.Args{})
	assert.Nil(t, err)
	assert.True(t, res)

	res, err = Both(EqualTo("dev-test")).And(Not(Empty())).Matches("dev-test", expect.Args{})
	assert.Nil(t, err)
	assert.True(t, res)

	res, err = Either(EqualTo("dev")).Or(EqualTo("test")).Matches("dev", expect.Args{})
	assert.Nil(t, err)
	assert.True(t, res)

	res, err = XOR(EqualTo("test"), EqualTo("dev")).Matches("test", expect.Args{})
	assert.Nil(t, err)
	assert.True(t, res)
}

package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpect_Compositions(t *testing.T) {
	m, err := Compose(StrictEqual("hello world")).
		And(Contain("hello")).
		And(HasPrefix("hello")).
		And(HasSuffix("world")).
		Match("hello world")

	assert.True(t, m.Pass)
	assert.Nil(t, err)

	m, err = Compose(StrictEqual("dev qa")).
		And(HasSuffix("dev")).
		Or(Contain("qa")).
		Match("dev qa")

	assert.True(t, m.Pass)
	assert.Nil(t, err)

	m, err = Compose(StrictEqual("testing")).
		Xor(Contain("test")).
		Match("testing")

	assert.False(t, m.Pass)
	assert.Nil(t, err)

	m, err = Compose(StrictEqual("hello world")).
		And(Contain("hello")).
		And(HasPrefix("world")).
		And(HasSuffix("hello")).
		Match("hello world")

	assert.False(t, m.Pass)
	assert.Nil(t, err)
}

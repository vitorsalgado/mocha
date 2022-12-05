package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpect_Compositions(t *testing.T) {
	m, err := Compose(Equal("hello world")).
		And(Contain("hello")).
		And(HasPrefix("hello")).
		And(HasSuffix("world")).
		Match("hello world")

	assert.True(t, m.OK)
	assert.Nil(t, err)

	m, err = Compose(Equal("dev qa")).
		And(HasSuffix("dev")).
		Or(Contain("qa")).
		Match("dev qa")

	assert.True(t, m.OK)
	assert.Nil(t, err)

	m, err = Compose(Equal("testing")).
		Xor(Contain("test")).
		Match("testing")

	assert.False(t, m.OK)
	assert.Nil(t, err)

	m, err = Compose(Equal("hello world")).
		And(Contain("hello")).
		And(HasPrefix("world")).
		And(HasSuffix("hello")).
		Match("hello world")

	assert.False(t, m.OK)
	assert.Nil(t, err)
}

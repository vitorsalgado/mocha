package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpect_Compositions(t *testing.T) {
	m, err := Compose(ToEqual("hello world")).
		And(ToContain("hello")).
		And(ToHavePrefix("hello")).
		And(ToHaveSuffix("world")).
		Match("hello world")

	assert.True(t, m.OK)
	assert.Nil(t, err)

	m, err = Compose(ToEqual("dev qa")).
		And(ToHaveSuffix("dev")).
		Or(ToContain("qa")).
		Match("dev qa")

	assert.True(t, m.OK)
	assert.Nil(t, err)

	m, err = Compose(ToEqual("testing")).
		Xor(ToContain("test")).
		Match("testing")

	assert.False(t, m.OK)
	assert.Nil(t, err)

	m, err = Compose(ToEqual("hello world")).
		And(ToContain("hello")).
		And(ToHavePrefix("world")).
		And(ToHaveSuffix("hello")).
		Match("hello world")

	assert.False(t, m.OK)
	assert.Nil(t, err)
}

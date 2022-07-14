package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpect_Compositions(t *testing.T) {
	m, err := ToEqual("hello world").
		And(ToContain("hello")).
		And(ToHavePrefix("hello")).
		And(ToHaveSuffix("world")).
		Matches("hello world", emptyArgs())

	assert.True(t, m)
	assert.Nil(t, err)

	m, err = ToEqual("dev qa").
		And(ToHaveSuffix("dev")).
		Or(ToContain("qa")).
		Matches("dev qa", emptyArgs())

	assert.True(t, m)
	assert.Nil(t, err)

	m, err = ToEqual("testing").
		Xor(ToContain("test")).
		Matches("testing", emptyArgs())

	assert.False(t, m)
	assert.Nil(t, err)

	m, err = ToEqual("hello world").
		And(ToContain("hello")).
		And(ToHavePrefix("world")).
		And(ToHaveSuffix("hello")).
		Matches("hello world", emptyArgs())

	assert.False(t, m)
	assert.Nil(t, err)
}

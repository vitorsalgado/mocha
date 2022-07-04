package matchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLen(t *testing.T) {
	str := "hello world -  "
	result, err := Len[string](15).Matches(str, emptyArgs())

	assert.Nil(t, err)
	assert.True(t, result)
}

package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLen(t *testing.T) {
	str := "hello world -  "
	result, err := Len[string](15)(str, Params{})

	assert.Nil(t, err)
	assert.True(t, result)
}

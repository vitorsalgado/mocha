package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLen(t *testing.T) {
	str := "hello world -  "
	result, err := ToHaveLen[string](15).Matches(str, emptyArgs())

	assert.Nil(t, err)
	assert.True(t, result)
}

package expect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLen(t *testing.T) {
	str := "hello world -  "
	result, err := ToHaveLen(15).Match(str)

	assert.Nil(t, err)
	assert.True(t, result.OK)
}

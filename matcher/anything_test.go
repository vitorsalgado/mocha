package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnything(t *testing.T) {
	res, err := Anything().Match("")
	assert.NoError(t, err)
	assert.True(t, res.Pass)
}

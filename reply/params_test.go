package reply

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParameters(t *testing.T) {
	params := Parameters()
	key1 := "k1"
	val1 := "test"
	key2 := "k2"
	val2 := 100

	params.Set(key1, val1)
	params.Set(key2, val2)

	v, ok := params.Get(key1)
	assert.Equal(t, val1, v)
	assert.True(t, ok)

	all := params.GetAll()
	assert.Equal(t, 2, len(all))

	v, ok = params.Get(key2)
	assert.Equal(t, val2, v)
	assert.True(t, ok)

	v, ok = params.Get("unknown")
	assert.Nil(t, v)
	assert.False(t, ok)

	assert.True(t, params.Has(key1))
	assert.True(t, params.Has(key2))
	assert.False(t, params.Has("unknown"))

	params.Remove(key1)
	assert.False(t, params.Has(key1))
}

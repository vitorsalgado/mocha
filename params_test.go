package mocha

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParameters(t *testing.T) {
	ctx := context.Background()
	params := Parameters()
	key1 := "k1"
	val1 := "test"
	key2 := "k2"
	val2 := 100

	params.Set(ctx, key1, val1)
	params.Set(ctx, key2, val2)

	v, ok := params.Get(ctx, key1)
	assert.Equal(t, val1, v)
	assert.True(t, ok)

	all := params.GetAll(ctx)
	assert.Equal(t, 2, len(all))

	v, ok = params.Get(ctx, key2)
	assert.Equal(t, val2, v)
	assert.True(t, ok)

	v, ok = params.Get(ctx, "unknown")
	assert.Nil(t, v)
	assert.False(t, ok)

	assert.True(t, params.Has(ctx, key1))
	assert.True(t, params.Has(ctx, key2))
	assert.False(t, params.Has(ctx, "unknown"))

	params.Remove(ctx, key1)
	assert.False(t, params.Has(ctx, key1))
}

package mocha

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParameters(t *testing.T) {
	ctx := context.Background()
	params := newInMemoryParameters()
	key1 := "k1"
	val1 := "test"
	key2 := "k2"
	val2 := 100

	require.NoError(t, params.Set(ctx, key1, val1))
	require.NoError(t, params.Set(ctx, key2, val2))

	v, ok, err := params.Get(ctx, key1)
	require.NoError(t, err)
	require.Equal(t, val1, v)
	require.True(t, ok)

	all, err := params.GetAll(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(all))

	v, ok, err = params.Get(ctx, key2)
	require.NoError(t, err)
	require.Equal(t, val2, v)
	require.True(t, ok)

	v, ok, err = params.Get(ctx, "unknown")
	require.NoError(t, err)
	require.Nil(t, v)
	require.False(t, ok)

	h1, err := params.Has(ctx, key1)
	require.NoError(t, err)
	h2, err := params.Has(ctx, key2)
	require.NoError(t, err)
	h3, err := params.Has(ctx, "unknown")
	require.NoError(t, err)

	require.True(t, h1)
	require.True(t, h2)
	require.False(t, h3)

	require.NoError(t, params.Remove(ctx, key1))

	h1, err = params.Has(ctx, key1)

	require.NoError(t, err)
	require.False(t, h1)
}

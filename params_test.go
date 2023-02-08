package mocha

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

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

func TestParameters_Concurrency(t *testing.T) {
	ctx := context.Background()
	params := newInMemoryParameters()
	jobs := 10
	wg := sync.WaitGroup{}

	for i := 0; i < jobs; i++ {
		wg.Add(1)

		key := "key--" + strconv.FormatInt(int64(i), 10)
		value := "value"

		go func(index int) {
			if index%2 == 0 {
				time.Sleep(100 * time.Millisecond)
			}

			kk := "key--" + strconv.FormatInt(int64(i), 10)
			vv := "value"

			err := params.Set(ctx, kk, vv)

			require.NoError(t, err)

			v, exists, err := params.Get(ctx, kk)

			require.NoError(t, err)
			require.True(t, exists)
			require.Equal(t, vv, v)

			_ = params.Remove(ctx, "k001")
			_ = params.Remove(ctx, kk)

			wg.Done()
		}(i)

		err := params.Set(ctx, key, value)
		require.NoError(t, err)
	}

	_, _, _ = params.Get(ctx, "key--100")
	_, _, _ = params.Get(ctx, "key--0")

	_ = params.Set(ctx, "k1", "v1")
	_ = params.Set(ctx, "k2", "v2")
	_ = params.Set(ctx, "k1", "v2")

	wg.Wait()
}

package foundation

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParameters(t *testing.T) {
	params := NewInMemoryParameters()
	key1 := "k1"
	val1 := "test"
	key2 := "k2"
	val2 := 100

	require.NoError(t, params.Set(key1, val1))
	require.NoError(t, params.Set(key2, val2))

	v, err := params.Get(key1)
	require.NoError(t, err)
	require.Equal(t, val1, v)

	all, err := params.GetAll()
	require.NoError(t, err)
	require.Equal(t, 2, len(all))

	v, err = params.Get(key2)
	require.NoError(t, err)
	require.Equal(t, val2, v)

	v, err = params.Get("unknown")
	require.NoError(t, err)
	require.Nil(t, v)

	h1, err := params.Has(key1)
	require.NoError(t, err)
	h2, err := params.Has(key2)
	require.NoError(t, err)
	h3, err := params.Has("unknown")
	require.NoError(t, err)

	require.True(t, h1)
	require.True(t, h2)
	require.False(t, h3)

	require.NoError(t, params.Remove(key1))

	h1, err = params.Has(key1)

	require.NoError(t, err)
	require.False(t, h1)
}

func TestParametersConcurrency(t *testing.T) {
	params := NewInMemoryParameters()
	jobs := 10
	wg := sync.WaitGroup{}

	for i := 0; i < jobs; i++ {
		wg.Add(1)

		key := "key--" + strconv.FormatInt(int64(i), 10)
		value := "value"

		go func(index int) {
			kk := "key--" + strconv.FormatInt(int64(index), 10)
			vv := "value"

			err := params.Set(kk, vv)

			assert.NoError(t, err)

			v, err := params.Get(kk)

			assert.NoError(t, err)
			assert.Equal(t, vv, v)

			_ = params.Remove("k001")
			_ = params.Remove(kk)

			wg.Done()
		}(i)

		err := params.Set(key, value)
		require.NoError(t, err)
	}

	_, _ = params.Get("key--100")
	_, _ = params.Get("key--0")

	_ = params.Set("k1", "v1")
	_ = params.Set("k2", "v2")
	_ = params.Set("k1", "v2")

	wg.Wait()
}

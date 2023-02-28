package mocha

import (
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSequentialReply(t *testing.T) {
	t.Parallel()

	t.Run("should return replies based configure sequence and return error when over", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		rv := &RequestValues{RawRequest: req, URL: req.URL}
		builder := Seq(InternalServerError(), BadRequest(), OK(), NotFound())

		res, err := builder.Build(nil, rv)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

		res, err = builder.Build(nil, rv)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		res, err = builder.Build(nil, rv)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		res, err = builder.Build(nil, rv)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)

		_, err = builder.Build(nil, rv)
		assert.NotNil(t, err)
	})

	t.Run("should return replies based configure sequence and return error when over", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		rv := &RequestValues{RawRequest: req, URL: req.URL}
		require.NoError(t, err)

		builder := Seq().Add(OK()).OnSequenceEnded(NotFound())

		res, err := builder.Build(nil, rv)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		res, err = builder.Build(nil, rv)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

func TestSequentialReplyShouldReturnErrorWhenSequenceDoesNotContainReplies(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	rv := &RequestValues{RawRequest: req, URL: req.URL}
	res, err := Seq().Build(nil, rv)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestSequentialReplyValidate(t *testing.T) {
	seq := Seq()

	require.Error(t, seq.beforeBuild(nil))

	seq.Add(OK())

	require.NoError(t, seq.beforeBuild(nil))
}

func TestSeqRace(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	rv := &RequestValues{RawRequest: req, URL: req.URL}
	builder := Seq(InternalServerError(), BadRequest(), OK(), NotFound())

	jobs := 3
	wg := sync.WaitGroup{}

	for i := 0; i < jobs; i++ {
		wg.Add(1)
		go func(index int) {
			res, err := builder.Build(nil, rv)
			require.NoError(t, err)
			require.True(t, res.StatusCode != StatusNoMatch)

			builder.curHits()

			wg.Done()
		}(i)

		builder.curHits()
	}

	res, err := builder.Build(nil, rv)
	require.NoError(t, err)
	require.True(t, res.StatusCode != StatusNoMatch)

	builder.curHits()

	require.Eventually(t, func() bool {
		wg.Wait()
		return true
	}, 1*time.Second, 100*time.Millisecond)
	require.Equal(t, 4, builder.curHits())
}

package reply

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSequential(t *testing.T) {
	t.Run("should return replies based configure sequence and return error when over", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		builder := Seq(InternalServerError(), BadRequest(), OK(), NotFound())

		res, err := builder.Build(nil, newReqValues(req))
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

		res, err = builder.Build(nil, newReqValues(req))
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		res, err = builder.Build(nil, newReqValues(req))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		res, err = builder.Build(nil, newReqValues(req))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)

		_, err = builder.Build(nil, newReqValues(req))
		assert.NotNil(t, err)
	})

	t.Run("should return replies based configure sequence and return error when over", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		require.NoError(t, err)

		builder := Seq().Add(OK()).AfterSequenceEnded(NotFound())

		res, err := builder.Build(nil, newReqValues(req))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		res, err = builder.Build(nil, newReqValues(req))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

func TestShouldReturnErrorWhenSequenceDoesNotContainReplies(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	res, err := Seq().Build(nil, newReqValues(req))
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

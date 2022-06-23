package reply

import (
	"github.com/stretchr/testify/assert"
	"github.com/vitorsalgado/mocha/mock"
	"net/http"
	"testing"
)

func TestSequential(t *testing.T) {
	t.Run("should return replies based configure sequence and return error when over", func(t *testing.T) {
		m := mock.Mock{Name: "mock_test"}
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		builder := Sequential().
			Add(InternalServerError(), BadRequest(), OK()).
			Then(NotFound())

		m.Hit()
		res, err := builder.Build(req, &m)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.Status)

		m.Hit()
		res, err = builder.Build(req, &m)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.Status)

		m.Hit()
		res, err = builder.Build(req, &m)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.Status)

		m.Hit()
		res, err = builder.Build(req, &m)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.Status)

		m.Hit()
		res, err = builder.Build(req, &m)
		assert.NotNil(t, err)
	})

	t.Run("should return replies based configure sequence and return error when over", func(t *testing.T) {
		m := mock.Mock{Name: "mock_test"}
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		builder := Sequential().Add(OK()).ReplyOnSequenceEnded(NotFound())

		m.Hit()
		res, err := builder.Build(req, &m)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.Status)

		m.Hit()
		res, err = builder.Build(req, &m)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.Status)
	})
}

func TestShouldReturnErrorWhenSequenceDoesNotContainReplies(t *testing.T) {
	res, err := Sequential().Build(nil, nil)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

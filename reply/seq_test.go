package reply

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSequential(t *testing.T) {
	t.Run("should return replies based configure sequence and return error when over", func(t *testing.T) {
		m := &mmock{}
		handler := m.On("Hits").Return(0)
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		builder := Seq().
			Add(InternalServerError(), BadRequest(), OK(), NotFound())

		handler.Unset()
		handler = m.On("Hits").Return(1)

		res, err := builder.Build(req, m, nil)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.Status)

		handler.Unset()
		handler = m.On("Hits").Return(2)

		res, err = builder.Build(req, m, nil)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.Status)

		handler.Unset()
		handler = m.On("Hits").Return(3)

		res, err = builder.Build(req, m, nil)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.Status)

		handler.Unset()
		handler = m.On("Hits").Return(4)

		res, err = builder.Build(req, m, nil)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.Status)

		handler.Unset()
		m.On("Hits").Return(5)

		_, err = builder.Build(req, m, nil)
		assert.NotNil(t, err)
	})

	t.Run("should return replies based configure sequence and return error when over", func(t *testing.T) {
		m := &mmock{}
		handler := m.On("Hits").Return(0)
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		builder := Seq().Add(OK()).AfterEnded(NotFound())

		handler.Unset()
		handler = m.On("Hits").Return(1)

		res, err := builder.Build(req, m, nil)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.Status)

		handler.Unset()
		m.On("Hits").Return(2)

		res, err = builder.Build(req, m, nil)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.Status)
	})
}

func TestShouldReturnErrorWhenSequenceDoesNotContainReplies(t *testing.T) {
	m := &mmock{}
	m.On("Hits").Return(0)

	res, err := Seq().Build(nil, m, nil)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

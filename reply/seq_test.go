package reply

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSequential(t *testing.T) {
	t.Run("should return replies based configure sequence and return error when over", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), KArg, &Arg{M: M{1}})
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080", nil)
		builder := Seq(InternalServerError(), BadRequest(), OK(), NotFound())

		res, err := builder.Build(nil, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.Status)

		ctx = context.WithValue(context.Background(), KArg, &Arg{M: M{2}})
		req, _ = http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080", nil)
		res, err = builder.Build(nil, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.Status)

		ctx = context.WithValue(context.Background(), KArg, &Arg{M: M{3}})
		req, _ = http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080", nil)
		res, err = builder.Build(nil, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.Status)

		ctx = context.WithValue(context.Background(), KArg, &Arg{M: M{4}})
		req, _ = http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080", nil)
		res, err = builder.Build(nil, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.Status)

		ctx = context.WithValue(context.Background(), KArg, &Arg{M: M{5}})
		req, _ = http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080", nil)
		_, err = builder.Build(nil, req)
		assert.NotNil(t, err)
	})

	t.Run("should return replies based configure sequence and return error when over", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), KArg, &Arg{M: M{0}})
		_, _ = http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080", nil)
		builder := Seq().Add(OK()).AfterEnded(NotFound())

		ctx = context.WithValue(context.Background(), KArg, &Arg{M: M{1}})
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080", nil)
		res, err := builder.Build(nil, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.Status)

		ctx = context.WithValue(context.Background(), KArg, &Arg{M: M{2}})
		req, _ = http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080", nil)
		res, err = builder.Build(nil, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.Status)
	})
}

func TestShouldReturnErrorWhenSequenceDoesNotContainReplies(t *testing.T) {
	ctx := context.WithValue(context.Background(), KArg, &Arg{M: M{0}})
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080", nil)
	res, err := Seq().Build(nil, req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

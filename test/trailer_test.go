package test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/internal/mimetype"
	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/reply"
)

func TestTrailer_WithBody(t *testing.T) {
	m := mocha.New()
	m.MustStart()

	defer m.Close()

	scoped := m.MustMock(
		mocha.Get(matcher.URLPath("/test")).
			Reply(reply.
				OK().
				PlainText("hello world").
				Header(header.ContentType, mimetype.TextPlain).
				Trailer("trailer-1", "trailer-1-value").
				Trailer("trailer-2", "trailer-2-value")))

	req := testutil.Get(m.URL() + "/test")

	res, err := req.Do()
	require.NoError(t, err)

	defer res.Body.Close()

	scoped.AssertCalled(t)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, mimetype.TextPlain, res.Header.Get(header.ContentType))
	assert.Len(t, res.Trailer, 2)

	b, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Len(t, res.Trailer, 2)
	assert.Equal(t, "hello world", string(b))
	assert.Equal(t, "trailer-1-value", res.Trailer.Get("trailer-1"))
	assert.Equal(t, "trailer-2-value", res.Trailer.Get("trailer-2"))
}

func TestTrailer_WithoutBody(t *testing.T) {
	m := mocha.New()
	m.MustStart()

	defer m.Close()

	scoped := m.MustMock(
		mocha.Get(matcher.URLPath("/test")).
			Reply(reply.
				OK().
				Header(header.ContentType, mimetype.TextPlain).
				Trailer("trailer-1", "trailer-1-value").
				Trailer("trailer-2", "trailer-2-value")))

	req := testutil.Get(m.URL() + "/test")

	res, err := req.Do()
	require.NoError(t, err)

	scoped.AssertCalled(t)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, mimetype.TextPlain, res.Header.Get(header.ContentType))
	assert.Len(t, res.Trailer, 2)

	_, err = io.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Len(t, res.Trailer, 2)
	assert.Equal(t, "trailer-1-value", res.Trailer.Get("trailer-1"))
	assert.Equal(t, "trailer-2-value", res.Trailer.Get("trailer-2"))
}

package test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/matcher"
	mhttp2 "github.com/vitorsalgado/mocha/v3/mhttp"
	"github.com/vitorsalgado/mocha/v3/misc"
)

func TestTrailer_WithBody(t *testing.T) {
	m := mhttp2.NewAPI()
	m.MustStart()

	defer m.Close()

	scoped := m.MustMock(
		mhttp2.Get(matcher.URLPath("/test")).
			Reply(mhttp2.OK().
				PlainText("hello world").
				Header(misc.HeaderContentType, misc.MIMETextPlain).
				Trailer("trailer-1", "trailer-1-value").
				Trailer("trailer-2", "trailer-2-value")))

	res, err := http.Get(m.URL() + "/test")
	require.NoError(t, err)

	defer res.Body.Close()

	require.True(t, scoped.AssertCalled(t))
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, misc.MIMETextPlain, res.Header.Get(misc.HeaderContentType))
	require.Len(t, res.Trailer, 2)

	b, err := io.ReadAll(res.Body)

	require.NoError(t, err)
	require.Len(t, res.Trailer, 2)
	require.Equal(t, "hello world", string(b))
	require.Equal(t, "trailer-1-value", res.Trailer.Get("trailer-1"))
	require.Equal(t, "trailer-2-value", res.Trailer.Get("trailer-2"))
}

func TestTrailer_WithoutBody(t *testing.T) {
	m := mhttp2.NewAPI()
	m.MustStart()

	defer m.Close()

	scoped := m.MustMock(
		mhttp2.Get(matcher.URLPath("/test")).
			Reply(mhttp2.OK().
				Header(misc.HeaderContentType, misc.MIMETextPlain).
				Trailer("trailer-1", "trailer-1-value").
				Trailer("trailer-2", "trailer-2-value")))

	res, err := http.Get(m.URL() + "/test")

	require.NoError(t, err)
	require.True(t, scoped.AssertCalled(t))
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, misc.MIMETextPlain, res.Header.Get(misc.HeaderContentType))
	require.Len(t, res.Trailer, 2)

	_, err = io.ReadAll(res.Body)

	require.NoError(t, err)
	require.Len(t, res.Trailer, 2)
	require.Equal(t, "trailer-1-value", res.Trailer.Get("trailer-1"))
	require.Equal(t, "trailer-2-value", res.Trailer.Get("trailer-2"))
}

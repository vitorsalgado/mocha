package test

import (
	"compress/gzip"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	. "github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/reply"
)

func TestCompressedResponse_GZIP(t *testing.T) {
	m := mocha.New()
	m.MustStart()

	defer m.Close()

	m.MustMock(mocha.Get(URLPath("/test")).
		Reply(reply.
			OK().
			BodyText("hello world").
			Gzip()))

	req := testutil.Get(m.URL() + "/test")
	res, err := req.Do()
	require.NoError(t, err)

	defer res.Body.Close()

	gz, err := gzip.NewReader(res.Body)
	require.NoError(t, err)

	b, err := io.ReadAll(gz)
	require.NoError(t, err)

	assert.False(t, res.Uncompressed)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "hello world", string(b))
}

func Test_GZIPProxiedResponse(t *testing.T) {
	p := mocha.New()
	p.MustStart()

	defer p.Close()

	ps := p.MustMock(mocha.Get(URLPath("/test")).
		Reply(reply.
			OK().
			BodyText("hello world").
			Gzip()))

	m := mocha.New()
	m.MustStart()

	defer m.Close()

	ms := m.MustMock(mocha.Get(URLPath("/test")).
		Reply(reply.From(p.URL())))

	req := testutil.Get(m.URL() + "/test")
	res, err := req.Do()
	require.NoError(t, err)

	defer res.Body.Close()

	gz, err := gzip.NewReader(res.Body)
	require.NoError(t, err)

	b, err := io.ReadAll(gz)
	require.NoError(t, err)

	assert.False(t, res.Uncompressed)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "hello world", string(b))
	assert.True(t, ps.HasBeenCalled())
	assert.True(t, ms.HasBeenCalled())
}

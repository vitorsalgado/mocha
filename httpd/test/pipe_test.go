package test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitorsalgado/mocha/v3/foundation"
	"github.com/vitorsalgado/mocha/v3/httpd"
)

type pipingTest struct {
}

func (a *pipingTest) Pipe(conduit *foundation.Conduit) {
	for {
		select {
		case chunk, ok := <-conduit.In:
			if !ok {
				close(conduit.Out)
				return
			}

			txt := string(chunk.Data)

			conduit.Out <- &foundation.Chunk{
				Data: []byte(txt + "(extra)"),
			}
		}
	}
}

type pipingTest2 struct {
}

func (a *pipingTest2) Pipe(conduit *foundation.Conduit) {
	for {
		select {
		case chunk, ok := <-conduit.In:
			if !ok {
				close(conduit.Out)
				return
			}

			txt := string(chunk.Data)

			conduit.Out <- &foundation.Chunk{
				Data: []byte(txt + "(other)"),
			}
		}
	}
}

func TestPipes(t *testing.T) {
	m := mhttp.NewAPIWithT(t, mhttp.Setup().LogVerbosity(mhttp.LogBody))
	m.MustStart()

	scope := m.MustMock(mhttp.Getf("/test").Pipe(&pipingTest{}).Pipe(&pipingTest2{}).Reply(mhttp.OK().BodyText("hi")))

	client := &http.Client{}
	res, err := client.Get(m.URL("/test"))
	require.NoError(t, err)

	b, err := io.ReadAll(res.Body)

	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.True(t, scope.AssertCalled(t))
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "hi(extra)(other)", string(b))
}

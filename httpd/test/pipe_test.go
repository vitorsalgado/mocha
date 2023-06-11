package test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitorsalgado/mocha/v3/httpd"
	"github.com/vitorsalgado/mocha/v3/lib"
)

type pipingTest struct {
}

func (a *pipingTest) Pipe(conduit *lib.Conduit) {
	for chunk := range conduit.In {
		txt := string(chunk.Data)

		conduit.Out <- &lib.Chunk{
			Data: []byte(txt + "(extra)"),
		}
	}

	close(conduit.Out)
}

type pipingTest2 struct {
}

func (a *pipingTest2) Pipe(conduit *lib.Conduit) {
	for chunk := range conduit.In {
		txt := string(chunk.Data)

		conduit.Out <- &lib.Chunk{
			Data: []byte(txt + "(other)"),
		}
	}

	close(conduit.Out)

}

func TestPipes(t *testing.T) {
	m := httpd.NewAPIWithT(t, httpd.Setup().LogVerbosity(httpd.LogBody))
	m.MustStart()

	scope := m.MustMock(httpd.Getf("/test").Pipe(&pipingTest{}).Pipe(&pipingTest2{}).Reply(httpd.OK().BodyText("hi")))

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

package test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/dzhttp"
	"github.com/vitorsalgado/mocha/v3/dzstd"
)

type pipingTest struct {
}

func (a *pipingTest) Pipe(conduit *dzstd.Conduit) {
	for chunk := range conduit.In {
		txt := string(chunk.Data)

		conduit.Out <- &dzstd.Chunk{
			Data: []byte(txt + "(extra)"),
		}
	}

	close(conduit.Out)
}

type pipingTest2 struct {
}

func (a *pipingTest2) Pipe(conduit *dzstd.Conduit) {
	for chunk := range conduit.In {
		txt := string(chunk.Data)

		conduit.Out <- &dzstd.Chunk{
			Data: []byte(txt + "(other)"),
		}
	}

	close(conduit.Out)

}

func TestPipes(t *testing.T) {
	m := dzhttp.NewAPIWithT(t, dzhttp.Setup().LogVerbosity(dzhttp.LogBody))
	m.MustStart()

	scope := m.MustMock(dzhttp.Getf("/test").Pipe(&pipingTest{}).Pipe(&pipingTest2{}).Reply(dzhttp.OK().BodyText("hi")))

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

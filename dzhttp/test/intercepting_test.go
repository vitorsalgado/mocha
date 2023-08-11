package test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/dzhttp"
	"github.com/vitorsalgado/mocha/v3/dzstd"
)

type TestInterceptor struct{ txt string }

func (in *TestInterceptor) Intercept(p []byte, chain dzstd.Chain) (n int, err error) {
	return chain.Next(append(p, []byte(in.txt)...))
}

func TestInterceptors(t *testing.T) {
	m := dzhttp.NewAPI(dzhttp.Setup().LogVerbosity(dzhttp.LogBody)).CloseWithT(t)
	m.MustStart()

	in1 := &TestInterceptor{"(extra)"}
	in2 := &TestInterceptor{"(other)"}

	scope := m.MustMock(dzhttp.Getf("/test").Intercept(in1, in2).Reply(dzhttp.OK().BodyText("hi")))

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

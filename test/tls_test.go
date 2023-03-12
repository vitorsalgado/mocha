package test

import (
	"crypto/tls"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3"
	. "github.com/vitorsalgado/mocha/v3/matcher"
)

func TestTLS(t *testing.T) {
	certFile := "testdata/tls/srv.crt"
	keyFile := "testdata/tls/srv.key"
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	require.NoError(t, err)

	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}

	testCases := []struct {
		name   string
		config mocha.Configurer
	}{
		{"basic", mocha.Setup()},
		{"custom key/pair", mocha.Setup().TLSLoadX509KeyPair(certFile, keyFile)},
		{"tls config", mocha.Setup().TLSConfig(&tls.Config{Certificates: []tls.Certificate{cert}})},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := mocha.NewT(t, tc.config)
			m.MustStartTLS()

			scoped := m.MustMock(
				mocha.Get(URLPath("/test")).
					Header("test", StrictEqual("hello")).
					Reply(mocha.OK().PlainText("hello+world")))

			req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)
			req.Header.Add("test", "hello")

			res, err := client.Do(req)
			require.NoError(t, err)

			b, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			require.NoError(t, res.Body.Close())
			require.Equal(t, "hello+world", string(b))
			require.True(t, scoped.HasBeenCalled())
		})
	}
}

func TestTLS_FileSetup(t *testing.T) {
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	m := mocha.NewT(t, mocha.UseConfig("testdata/tls/tls_app.yaml"))
	m.MustStartTLS()

	scoped := m.MustMock(
		mocha.Get(URLPath("/test")).
			Header("test", StrictEqual("hello")).
			Reply(mocha.OK().PlainText("hello+world")))

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)
	req.Header.Add("test", "hello")

	res, err := client.Do(req)
	require.NoError(t, err)

	b, err := io.ReadAll(res.Body)

	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, "hello+world", string(b))
	require.True(t, scoped.HasBeenCalled())
}

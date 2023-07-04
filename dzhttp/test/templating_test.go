package test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/dzhttp"
)

func TestTemplating(t *testing.T) {
	m := dzhttp.NewAPI()
	m.MustStart()
	defer m.Close()

	m.Parameters().MustSet("test", "hi")
	m.MustMock(dzhttp.FromFile("testdata/templating/templating_01.yaml"))

	httpClient := &http.Client{}

	req, _ := http.NewRequest(http.MethodPost, m.URL()+"/test/templating", strings.NewReader("hello world -> hi"))
	req.Header.Add("template", "true")
	req.AddCookie(&http.Cookie{Name: "cookie_test", Value: "cookie_value"})

	res, err := httpClient.Do(req)
	require.NoError(t, err)

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "test --> go --> POST --> true cookie_value Path 0 test Path 1 templating", string(b))
}

func TestTemplating_BodyFilename(t *testing.T) {
	m := dzhttp.NewAPI()
	m.MustStart()
	defer m.Close()

	m.MustMock(dzhttp.FromFile("testdata/templating/templating_02.yaml"))

	httpClient := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test/templating_02_data.txt", nil)
	res, err := httpClient.Do(req)
	require.NoError(t, err)

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "hello world\n", string(b))
}

func TestTemplating_BodyFilename_BodyTemplate(t *testing.T) {
	m := dzhttp.NewAPI()
	m.MustStart()
	defer m.Close()

	m.MustMock(dzhttp.FromFile("testdata/templating/templating_03.yaml"))
	_ = m.Parameters().Set("message", "hello world")

	httpClient := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/templating_03_data.txt", nil)
	res, err := httpClient.Do(req)
	require.NoError(t, err)

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "Method is: GET\nMessage is: hello world\n", string(b))
}

func TestTemplating_Header(t *testing.T) {
	m := dzhttp.NewAPI()
	m.MustStart()
	defer m.Close()

	m.MustMock(dzhttp.FromFile("testdata/templating/templating_04.yaml"))
	m.Parameters().MustSet("test", "ok")
	m.Parameters().MustSet("context", "test")

	httpClient := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test/templating", nil)
	res, err := httpClient.Do(req)
	require.NoError(t, err)

	defer res.Body.Close()

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "ok", res.Header.Get("test"))
	require.Equal(t, "test", res.Header.Get("ctx"))
	require.Equal(t, "hi", res.Header.Get("message"))
}

func TestTemplating_All(t *testing.T) {
	m := dzhttp.NewAPI()
	m.MustStart()
	defer m.Close()

	m.Parameters().MustSet("test", "ok")
	m.Parameters().MustSet("context", "test")
	m.MustMock(dzhttp.FromFile("testdata/templating/templating_05.yaml"))

	httpClient := &http.Client{}

	req, _ := http.NewRequest(http.MethodPost, m.URL()+"/test/templating", strings.NewReader("hello world -> ok"))
	req.Header.Add("template", "true")
	req.AddCookie(&http.Cookie{Name: "cookie_test", Value: "cookie_value"})

	res, err := httpClient.Do(req)
	require.NoError(t, err)

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "test --> go --> POST --> true cookie_value Path 0 test Path 1 templating", string(b))
	require.Equal(t, "ok", res.Header.Get("test"))
	require.Equal(t, "test", res.Header.Get("ctx"))
	require.Equal(t, "hi", res.Header.Get("message"))
}

package webhook

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/httpd"
	. "github.com/vitorsalgado/mocha/v3/httpd/httpval"
	. "github.com/vitorsalgado/mocha/v3/matcher"
)

const (
	_tlsCertFile       = "../test/testdata/cert/cert.pem"
	_tlsKeyFile        = "../test/testdata/cert/key.pem"
	_tlsClientCertFile = "../test/testdata/cert/cert_client.pem"
)

func TestWebHook_Run(t *testing.T) {
	key := "test_key"
	target := httpd.NewAPIWithT(t)
	target.MustStart()

	m := httpd.NewAPIWithT(t, httpd.Setup().PostAction(Name, New()))
	m.MustStart()

	testCases := []struct {
		name          string
		targetMock    *httpd.HTTPMockBuilder
		webhookDef    *httpd.PostActionDef
		expectedCalls int
	}{
		{"basic with default method",
			httpd.Getf("/third_party/hook").
				Reply(httpd.OK()),
			Setup().
				URL(target.URL("/third_party/hook")).
				Build(),
			1,
		},
		{"complex",
			httpd.Postf("/third_party/hook").
				Headerf("X-Key", key).
				Headerf(HeaderContentType, MIMETextPlain).
				Body(Eq("hi")).
				Reply(httpd.OK().PlainText("bye")),
			Setup().
				URL(target.URL("/third_party/hook")).
				Method(http.MethodPost).
				Header("X-Key", key).
				Header(HeaderContentType, MIMETextPlain).
				Body("hi").
				Build(),
			1,
		},
		{
			"no body",
			httpd.Postf("/third_party/hook").
				Headerf("X-Key", key).
				Headerf(HeaderContentType, MIMETextPlain).
				Reply(httpd.OK().PlainText("bye")),
			Setup().
				URL(target.URL("/third_party/hook")).
				Method(http.MethodPost).
				Header("X-Key", key).
				Header(HeaderContentType, MIMETextPlain).
				Build(),
			1,
		},
		{
			"bad request from target",
			httpd.Postf("/third_party/hook").
				Headerf(HeaderContentType, MIMETextPlain).
				Reply(httpd.BadRequest()),
			Setup().
				URL(target.URL("/third_party/hook")).
				Method(http.MethodPost).
				Header(HeaderContentType, MIMETextPlain).
				Build(),
			1,
		},
		{
			"transform",
			httpd.Putf("/third_party/hook/transformed").
				Headerf(HeaderContentType, MIMETextPlain).
				Body(Eq("hello world")).
				Reply(httpd.BadRequest()),
			Setup().
				URL(target.URL("/third_party/hook")).
				Method(http.MethodPost).
				Transform(func(input *httpd.PostActionInput, args *Input) error {
					args.URL += "/transformed"
					args.Method = http.MethodPut
					args.Body = "hello world"
					return nil
				}).
				Header(HeaderContentType, MIMETextPlain).
				Build(),
			1,
		},
		{
			"transform with error",
			httpd.Getf("/third_party/hook/transformed").
				Reply(httpd.BadRequest()),
			Setup().
				URL(target.URL("/third_party/hook")).
				Method(http.MethodGet).
				Transform(func(input *httpd.PostActionInput, args *Input) error {
					return errors.New("boom")
				}).
				Build(),
			0,
		},
		{
			"malformed url",
			httpd.Postf("/third_party/hook").
				Reply(httpd.OK()),
			Setup().
				URL(" -   " + string(rune(0x7f))).
				Method(http.MethodPost).
				Build(),
			0,
		},
		{
			"unable to build http request",
			httpd.Postf("/third_party/hook").Reply(httpd.OK()),
			Setup().
				URL(target.URL("/third_party/hook")).
				Method("ghjk%&^*()").
				Build(),
			0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer target.Clean()
			defer m.Clean()

			target.MustMock(tc.targetMock)

			m.MustMock(httpd.Getf("/test").
				PostAction(tc.webhookDef).
				Reply(httpd.NoContent()))

			client := &http.Client{}
			req, _ := http.NewRequest(http.MethodGet, m.URL("/test"), nil)
			res, err := client.Do(req)

			require.NoError(t, err)
			require.Equal(t, http.StatusNoContent, res.StatusCode)
			require.True(t, target.AssertNumberOfCalls(t, tc.expectedCalls))
			require.True(t, m.AssertNumberOfCalls(t, 1))
		})
	}
}

func TestWebHook_TLS(t *testing.T) {
	key := "test_key"
	cert, err := tls.LoadX509KeyPair(_tlsCertFile, _tlsKeyFile)
	require.NoError(t, err)

	caCert, err := os.ReadFile(filepath.Clean(_tlsClientCertFile))
	require.NoError(t, err)

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	target := httpd.NewAPIWithT(t, httpd.Setup().TLSMutual(_tlsCertFile, _tlsKeyFile, _tlsClientCertFile))
	target.MustStartTLS()
	target.MustMock(httpd.Postf("/third_party/hook/tls").
		Headerf("X-Key", key).
		Headerf(HeaderContentType, MIMETextPlainCharsetUTF8).
		Reply(httpd.OK().PlainText("hello")))

	m := httpd.NewAPIWithT(t, httpd.Setup().
		TLSMutual(_tlsCertFile, _tlsKeyFile, _tlsClientCertFile).
		PostAction(Name, New()))
	m.MustStartTLS()
	m.MustMock(httpd.Putf("/test").
		PostAction(Setup().
			URL(target.URL("/third_party/hook/tls")).
			Method(http.MethodPost).
			Header("X-Key", key).
			Header(HeaderContentType, MIMETextPlainCharsetUTF8).
			SSLVerify(true).
			Body("hi").
			Build()).
		Reply(httpd.NoContent()))

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				RootCAs:      certPool,
			}},
	}
	req, _ := http.NewRequest(http.MethodPut, m.URL("/test"), nil)
	res, err := client.Do(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, res.StatusCode)
	require.True(t, target.AssertNumberOfCalls(t, 1))
	require.True(t, m.AssertNumberOfCalls(t, 1))
}

func TestWebHook_FaultTarget(t *testing.T) {
	target := httpd.NewAPIWithT(t)
	target.MustStart()
	target.MustMock(httpd.Getf("/hook/fault").
		Delay(1 * time.Minute).
		Reply(httpd.OK().PlainText("hello")))

	m := httpd.NewAPIWithT(t, httpd.Setup().PostAction(Name, New()))
	m.MustStart()
	m.MustMock(httpd.Getf("/test").
		PostAction(Setup().
			URL(target.URL("/hook/fault")).
			Method(http.MethodGet).
			Build()).
		Reply(httpd.NoContent()))

	client := &http.Client{}
	timeout := 2 * time.Second
	exit := make(chan struct{}, 1)
	ti := time.AfterFunc(timeout, func() { exit <- struct{}{} })
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	defer cancel()
	defer ti.Stop()

	go func() {
		for {
			select {
			case <-exit:
				return
			case <-ctx.Done():
				target.CloseNow()
				return
			}
		}
	}()

	ctxReq, cancelReq := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelReq()

	req, _ := http.NewRequestWithContext(ctxReq, http.MethodGet, m.URL("/test"), nil)
	res, err := client.Do(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, res.StatusCode)
}

func TestWebHook_InvalidArgs(t *testing.T) {
	client := &http.Client{}
	testCases := []any{make(chan struct{}), nil}

	target := httpd.NewAPIWithT(t, httpd.Setup().HTTPClient(func() *http.Client {
		return client
	}))
	target.MustStart()
	target.MustMock(httpd.Getf("/hook").Reply(httpd.OK()))

	m := httpd.NewAPIWithT(t, httpd.Setup().PostAction(Name, New()))
	m.MustStart()

	for i, tc := range testCases {
		tc := tc
		t.Run(strconv.FormatInt(int64(i), 10), func(t *testing.T) {
			m.MustMock(httpd.Getf("/test").
				PostAction(&httpd.PostActionDef{Name: Name, RawParameters: tc}).
				Reply(httpd.NoContent()))

			res, err := client.Get(m.URL("/test"))

			require.NoError(t, err)
			require.Equal(t, http.StatusNoContent, res.StatusCode)

			m.Clean()
		})
	}
}

func TestWebHook_FileSetup(t *testing.T) {
	target := httpd.NewAPIWithT(t)
	target.MustStart()
	target.MustMock(httpd.Postf("/fs/hook").
		Headerf("hello", "world").
		Headerf("dev", "ok").
		ContentType(MIMEApplicationJSON).
		Body(Eq(`{"task": "done"}`)).
		Reply(httpd.NoContent()))

	m := httpd.NewAPIWithT(t, httpd.Setup().PostAction(Name, New()))
	m.SetData(map[string]any{"webhook_target": target.URL("/fs/hook")})

	m.MustStart()
	m.MustMock(httpd.FromFile("testdata/1_webhook_complete_setup.yaml"))

	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, m.URL("/test"), nil)
	res, err := client.Do(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, res.StatusCode)
	require.True(t, target.AssertNumberOfCalls(t, 1))
	require.True(t, m.AssertNumberOfCalls(t, 1))
}

func TestWebHook_InvalidFiles(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
	}{
		{"post actions must be an an array", "testdata/2_webhook_object_not_array.yaml"},
		{"item is not an object", "testdata/3_webhook_invalid_item_type.yaml"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := httpd.NewAPIWithT(t, httpd.Setup().PostAction(Name, New()))
			_, err := m.Mock(httpd.FromFile(tc.filename))

			require.Error(t, err)
		})
	}
}
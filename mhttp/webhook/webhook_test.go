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

	. "github.com/vitorsalgado/mocha/v3/matcher"
	mhttp2 "github.com/vitorsalgado/mocha/v3/mhttp"
	. "github.com/vitorsalgado/mocha/v3/misc"
)

const (
	_tlsCertFile       = "../test/testdata/cert/cert.pem"
	_tlsKeyFile        = "../test/testdata/cert/key.pem"
	_tlsClientCertFile = "../test/testdata/cert/cert_client.pem"
)

func TestWebHook_Run(t *testing.T) {
	key := "test_key"
	target := mhttp2.NewAPIWithT(t)
	target.MustStart()

	m := mhttp2.NewAPIWithT(t, mhttp2.Setup().PostAction(Name, New()))
	m.MustStart()

	testCases := []struct {
		name          string
		targetMock    *mhttp2.HTTPMockBuilder
		webhookDef    *mhttp2.PostActionDef
		expectedCalls int
	}{
		{"basic with default method",
			mhttp2.Getf("/third_party/hook").
				Reply(mhttp2.OK()),
			Setup().
				URL(target.URL("/third_party/hook")).
				Build(),
			1,
		},
		{"complex",
			mhttp2.Postf("/third_party/hook").
				Headerf("X-Key", key).
				Headerf(HeaderContentType, MIMETextPlain).
				Body(Eq("hi")).
				Reply(mhttp2.OK().PlainText("bye")),
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
			mhttp2.Postf("/third_party/hook").
				Headerf("X-Key", key).
				Headerf(HeaderContentType, MIMETextPlain).
				Reply(mhttp2.OK().PlainText("bye")),
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
			mhttp2.Postf("/third_party/hook").
				Headerf(HeaderContentType, MIMETextPlain).
				Reply(mhttp2.BadRequest()),
			Setup().
				URL(target.URL("/third_party/hook")).
				Method(http.MethodPost).
				Header(HeaderContentType, MIMETextPlain).
				Build(),
			1,
		},
		{
			"transform",
			mhttp2.Putf("/third_party/hook/transformed").
				Headerf(HeaderContentType, MIMETextPlain).
				Body(Eq("hello world")).
				Reply(mhttp2.BadRequest()),
			Setup().
				URL(target.URL("/third_party/hook")).
				Method(http.MethodPost).
				Transform(func(input *mhttp2.PostActionInput, args *Input) error {
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
			mhttp2.Getf("/third_party/hook/transformed").
				Reply(mhttp2.BadRequest()),
			Setup().
				URL(target.URL("/third_party/hook")).
				Method(http.MethodGet).
				Transform(func(input *mhttp2.PostActionInput, args *Input) error {
					return errors.New("boom")
				}).
				Build(),
			0,
		},
		{
			"malformed url",
			mhttp2.Postf("/third_party/hook").
				Reply(mhttp2.OK()),
			Setup().
				URL(" -   " + string(rune(0x7f))).
				Method(http.MethodPost).
				Build(),
			0,
		},
		{
			"unable to build http request",
			mhttp2.Postf("/third_party/hook").Reply(mhttp2.OK()),
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

			m.MustMock(mhttp2.Getf("/test").
				PostAction(tc.webhookDef).
				Reply(mhttp2.NoContent()))

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

	target := mhttp2.NewAPIWithT(t, mhttp2.Setup().TLSMutual(_tlsCertFile, _tlsKeyFile, _tlsClientCertFile))
	target.MustStartTLS()
	target.MustMock(mhttp2.Postf("/third_party/hook/tls").
		Headerf("X-Key", key).
		Headerf(HeaderContentType, MIMETextPlainCharsetUTF8).
		Reply(mhttp2.OK().PlainText("hello")))

	m := mhttp2.NewAPIWithT(t, mhttp2.Setup().
		TLSMutual(_tlsCertFile, _tlsKeyFile, _tlsClientCertFile).
		PostAction(Name, New()))
	m.MustStartTLS()
	m.MustMock(mhttp2.Putf("/test").
		PostAction(Setup().
			URL(target.URL("/third_party/hook/tls")).
			Method(http.MethodPost).
			Header("X-Key", key).
			Header(HeaderContentType, MIMETextPlainCharsetUTF8).
			SSLVerify(true).
			Body("hi").
			Build()).
		Reply(mhttp2.NoContent()))

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
	target := mhttp2.NewAPIWithT(t)
	target.MustStart()
	target.MustMock(mhttp2.Getf("/hook/fault").
		Delay(1 * time.Minute).
		Reply(mhttp2.OK().PlainText("hello")))

	m := mhttp2.NewAPIWithT(t, mhttp2.Setup().PostAction(Name, New()))
	m.MustStart()
	m.MustMock(mhttp2.Getf("/test").
		PostAction(Setup().
			URL(target.URL("/hook/fault")).
			Method(http.MethodGet).
			Build()).
		Reply(mhttp2.NoContent()))

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

	target := mhttp2.NewAPIWithT(t, mhttp2.Setup().HTTPClient(func() *http.Client {
		return client
	}))
	target.MustStart()
	target.MustMock(mhttp2.Getf("/hook").Reply(mhttp2.OK()))

	m := mhttp2.NewAPIWithT(t, mhttp2.Setup().PostAction(Name, New()))
	m.MustStart()

	for i, tc := range testCases {
		tc := tc
		t.Run(strconv.FormatInt(int64(i), 10), func(t *testing.T) {
			m.MustMock(mhttp2.Getf("/test").
				PostAction(&mhttp2.PostActionDef{Name: Name, RawParameters: tc}).
				Reply(mhttp2.NoContent()))

			res, err := client.Get(m.URL("/test"))

			require.NoError(t, err)
			require.Equal(t, http.StatusNoContent, res.StatusCode)

			m.Clean()
		})
	}
}

func TestWebHook_FileSetup(t *testing.T) {
	target := mhttp2.NewAPIWithT(t)
	target.MustStart()
	target.MustMock(mhttp2.Postf("/fs/hook").
		Headerf("hello", "world").
		Headerf("dev", "ok").
		ContentType(MIMEApplicationJSON).
		Body(Eq(`{"task": "done"}`)).
		Reply(mhttp2.NoContent()))

	m := mhttp2.NewAPIWithT(t, mhttp2.Setup().PostAction(Name, New()))
	m.SetData(map[string]any{"webhook_target": target.URL("/fs/hook")})

	m.MustStart()
	m.MustMock(mhttp2.FromFile("testdata/1_webhook_complete_setup.yaml"))

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
			m := mhttp2.NewAPIWithT(t, mhttp2.Setup().PostAction(Name, New()))
			_, err := m.Mock(mhttp2.FromFile(tc.filename))

			require.Error(t, err)
		})
	}
}

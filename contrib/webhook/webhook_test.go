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

	"github.com/vitorsalgado/mocha/v3"
	. "github.com/vitorsalgado/mocha/v3/matcher"
	. "github.com/vitorsalgado/mocha/v3/misc"
)

const (
	_tlsCertFile       = "../../test/testdata/cert/cert.pem"
	_tlsKeyFile        = "../../test/testdata/cert/key.pem"
	_tlsClientCertFile = "../../test/testdata/cert/cert_client.pem"
)

func TestWebHook_Run(t *testing.T) {
	key := "test_key"
	target := mocha.NewAPIWithT(t)
	target.MustStart()

	m := mocha.NewAPIWithT(t, mocha.Setup().PostAction(Name, New()))
	m.MustStart()

	testCases := []struct {
		name          string
		targetMock    *mocha.MockBuilder
		webhookDef    *mocha.PostActionDef
		expectedCalls int
	}{
		{"basic with default method",
			mocha.Getf("/third_party/hook").
				Reply(mocha.OK()),
			Setup().
				URL(target.URL("/third_party/hook")).
				Build(),
			1,
		},
		{"complex",
			mocha.Postf("/third_party/hook").
				Headerf("X-Key", key).
				Headerf(HeaderContentType, MIMETextPlain).
				Body(Eq("hi")).
				Reply(mocha.OK().PlainText("bye")),
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
			mocha.Postf("/third_party/hook").
				Headerf("X-Key", key).
				Headerf(HeaderContentType, MIMETextPlain).
				Reply(mocha.OK().PlainText("bye")),
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
			mocha.Postf("/third_party/hook").
				Headerf(HeaderContentType, MIMETextPlain).
				Reply(mocha.BadRequest()),
			Setup().
				URL(target.URL("/third_party/hook")).
				Method(http.MethodPost).
				Header(HeaderContentType, MIMETextPlain).
				Build(),
			1,
		},
		{
			"transform",
			mocha.Putf("/third_party/hook/transformed").
				Headerf(HeaderContentType, MIMETextPlain).
				Body(Eq("hello world")).
				Reply(mocha.BadRequest()),
			Setup().
				URL(target.URL("/third_party/hook")).
				Method(http.MethodPost).
				Transform(func(input *mocha.PostActionInput, args *Input) error {
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
			mocha.Getf("/third_party/hook/transformed").
				Reply(mocha.BadRequest()),
			Setup().
				URL(target.URL("/third_party/hook")).
				Method(http.MethodGet).
				Transform(func(input *mocha.PostActionInput, args *Input) error {
					return errors.New("boom")
				}).
				Build(),
			0,
		},
		{
			"malformed url",
			mocha.Postf("/third_party/hook").
				Reply(mocha.OK()),
			Setup().
				URL(" -   " + string(rune(0x7f))).
				Method(http.MethodPost).
				Build(),
			0,
		},
		{
			"unable to build http request",
			mocha.Postf("/third_party/hook").Reply(mocha.OK()),
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

			m.MustMock(mocha.Getf("/test").
				PostAction(tc.webhookDef).
				Reply(mocha.NoContent()))

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

	target := mocha.NewAPIWithT(t, mocha.Setup().TLSMutual(_tlsCertFile, _tlsKeyFile, _tlsClientCertFile))
	target.MustStartTLS()
	target.MustMock(mocha.Postf("/third_party/hook/tls").
		Headerf("X-Key", key).
		Headerf(HeaderContentType, MIMETextPlainCharsetUTF8).
		Reply(mocha.OK().PlainText("hello")))

	m := mocha.NewAPIWithT(t, mocha.Setup().
		TLSMutual(_tlsCertFile, _tlsKeyFile, _tlsClientCertFile).
		PostAction(Name, New()))
	m.MustStartTLS()
	m.MustMock(mocha.Putf("/test").
		PostAction(Setup().
			URL(target.URL("/third_party/hook/tls")).
			Method(http.MethodPost).
			Header("X-Key", key).
			Header(HeaderContentType, MIMETextPlainCharsetUTF8).
			SSLVerify(true).
			Body("hi").
			Build()).
		Reply(mocha.NoContent()))

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
	target := mocha.NewAPIWithT(t)
	target.MustStart()
	target.MustMock(mocha.Getf("/hook/fault").
		Delay(1 * time.Minute).
		Reply(mocha.OK().PlainText("hello")))

	m := mocha.NewAPIWithT(t, mocha.Setup().PostAction(Name, New()))
	m.MustStart()
	m.MustMock(mocha.Getf("/test").
		PostAction(Setup().
			URL(target.URL("/hook/fault")).
			Method(http.MethodGet).
			Build()).
		Reply(mocha.NoContent()))

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
	testCases := []any{make(chan struct{}), nil}

	target := mocha.NewAPIWithT(t)
	target.MustStart()
	target.MustMock(mocha.Getf("/hook").Reply(mocha.OK()))

	m := mocha.NewAPIWithT(t, mocha.Setup().PostAction(Name, New()))
	m.MustStart()

	for i, tc := range testCases {
		t.Run(strconv.FormatInt(int64(i), 10), func(t *testing.T) {
			m.MustMock(mocha.Getf("/test").
				PostAction(&mocha.PostActionDef{Name: Name, RawParameters: tc}).
				Reply(mocha.NoContent()))

			client := &http.Client{}
			res, err := client.Get(m.URL("/test"))

			require.NoError(t, err)
			require.Equal(t, http.StatusNoContent, res.StatusCode)
		})
	}
}

func TestWebHook_FileSetup(t *testing.T) {
	target := mocha.NewAPIWithT(t)
	target.MustStart()
	target.MustMock(mocha.
		Postf("/fs/hook").
		Headerf("hello", "world").
		Headerf("dev", "ok").
		ContentType(MIMEApplicationJSON).
		Body(Eq(`{"task": "done"}`)).
		Reply(mocha.NoContent()))

	m := mocha.NewAPIWithT(t, mocha.Setup().PostAction(Name, New()))
	m.SetData(map[string]any{"webhook_target": target.URL("/fs/hook")})

	m.MustStart()
	m.MustMock(mocha.FromFile("testdata/1_webhook_complete_setup.yaml"))

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
			m := mocha.NewAPIWithT(t, mocha.Setup().PostAction(Name, New()))
			_, err := m.Mock(mocha.FromFile(tc.filename))

			require.Error(t, err)
		})
	}
}

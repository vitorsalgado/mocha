package mhttp

import (
	"crypto/tls"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"text/template"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3/foundation"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/mhttp/cors"
	"github.com/vitorsalgado/mocha/v3/misc"
)

type testBodyParser struct{}

func (p testBodyParser) CanParse(content string, r *http.Request) bool {
	return content == misc.MIMETextPlain && r.Header.Get("x-test") == "num"
}

func (p testBodyParser) Parse(body []byte, _ *http.Request) (any, error) {
	return strconv.Atoi(string(body))
}

type customTestServer struct {
	decorated Server
}

func (s *customTestServer) Setup(app *HTTPMockApp, handler http.Handler) error {
	return s.decorated.Setup(app, handler)
}

func (s *customTestServer) Start() error {
	return s.decorated.Start()
}

func (s *customTestServer) StartTLS() error {
	return s.decorated.StartTLS()
}

func (s *customTestServer) Close() error {
	return s.decorated.Close()
}

func (s *customTestServer) CloseNow() error {
	return s.decorated.CloseNow()
}

func (s *customTestServer) Info() *ServerInfo {
	return s.decorated.Info()
}

func (s *customTestServer) S() any {
	return s.decorated.S()
}

func TestConfig(t *testing.T) {
	client := &http.Client{}

	t.Run("should run server with the custom given address", func(t *testing.T) {
		addr := os.Getenv("TEST_CUSTOM_ADDR")
		if addr == "" {
			addr = "127.0.0.1:3000"
		}

		m := NewAPI(Setup().Addr(addr)).CloseWithT(t)
		m.MustStart()

		defer m.Close()

		scoped := m.MustMock(
			Getf("/test").
				Reply(OK()))

		req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)
		res, err := client.Do(req)

		assert.NoError(t, err)
		assert.True(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Contains(t, m.server.Info().URL, addr)
	})

	t.Run("request body parsers from config should take precedence", func(t *testing.T) {
		m := NewAPI(Setup().RequestBodyParsers(&testBodyParser{}))
		m.MustStart()

		defer m.Close()

		scoped := m.MustMock(Post(matcher.URLPath("/test")).
			Body(matcher.StrictEqual(10)).
			Reply(OK()))

		req, _ := http.NewRequest(http.MethodPost, m.URL()+"/test", strings.NewReader("10"))
		req.Header.Add(misc.HeaderContentType, misc.MIMETextPlain)
		req.Header.Add("x-test", "num")

		res, err := client.Do(req)

		assert.NoError(t, err)
		assert.True(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("middlewares from config should take precedence", func(t *testing.T) {
		msg := "hello world"
		middleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.Header().Add("intercepted", "true")
					w.Header().Add(misc.HeaderContentType, misc.MIMETextPlain)
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(msg))
				})
		}

		m := NewAPI(Setup().Middlewares(middleware))
		m.MustStart()

		defer m.Close()

		scoped := m.MustMock(
			Get(matcher.URLPath("/test")).
				Reply(OK()))

		req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)
		res, err := client.Do(req)

		assert.NoError(t, err)
		assert.False(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "true", res.Header.Get("intercepted"))
	})

	t.Run("configure custom server", func(t *testing.T) {
		m := NewAPI(Setup().Server(&customTestServer{decorated: newServer()}))
		m.MustStart()

		defer m.Close()

		scoped := m.MustMock(
			Get(matcher.URLPath("/test")).
				Reply(OK()))

		req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)
		res, err := client.Do(req)

		assert.NoError(t, err)
		assert.True(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func TestConfigBuilder(t *testing.T) {
	addr := ""
	nm := "test"
	customLogger := zerolog.Nop()
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	m := NewAPI(Setup().
		Name(nm).
		Addr(addr).
		RootDir("test_root_dir").
		MockNotFoundStatusCode(http.StatusNotFound).
		RequestBodyParsers(&plainTextParser{}).
		Middlewares().
		CORS(&cors.DefaultConfig).
		Server(&httpTestServer{}).
		HandlerDecorator(func(handler http.Handler) http.Handler { return handler }).
		LogVerbosity(LogBasic).
		LogLevel(LogLevelInfo).
		LogPretty(false).
		LogBodyMaxSize(100).
		Logger(&customLogger).
		UseDescriptiveLogger().
		RedactHeader("header-1", "HEADER-2", "hEaDeR-3").
		Parameters(foundation.NewInMemoryParameters()).
		MockFilePatterns("test", "dev").
		Loader(&fileLoader{}).
		MockFileHandlers(&customMockFileHandler{}).
		TemplateEngine(newGoTemplate()).
		TemplateEngineFunctions(template.FuncMap{"trim": strings.TrimSpace}).
		Proxy(&ProxyConfig{}, &ProxyConfig{}).
		TLSConfig(tlsConfig).
		TLSCertKeyPair("test/testdata/cert/cert.pem", "test/testdata/cert/key.pem").
		TLSMutual("test/testdata/cert/cert.pem", "test/testdata/cert/key.pem", "test/testdata/cert/cert_client.pem"))
	conf := m.Config()

	assert.Equal(t, nm, conf.Name)
	assert.Equal(t, addr, conf.Addr)
	assert.Equal(t, "test_root_dir", conf.RootDir)
	assert.Equal(t, http.StatusNotFound, conf.RequestWasNotMatchedStatusCode)
	assert.Len(t, conf.RequestBodyParsers, 1)
	assert.Len(t, conf.Middlewares, 0)
	assert.Equal(t, &cors.DefaultConfig, conf.CORS)
	assert.NotNil(t, conf.HandlerDecorator)
	assert.Equal(t, LogBasic, conf.LogVerbosity)
	assert.Equal(t, LogLevelInfo, conf.LogLevel)
	assert.False(t, conf.LogPretty)
	assert.Equal(t, &customLogger, conf.Logger)
	assert.Equal(t, int64(100), conf.LogBodyMaxSize)
	assert.Equal(t, foundation.NewInMemoryParameters(), conf.Parameters)
	assert.Equal(t, []string{"test", "dev"}, conf.MockFileSearchPatterns)
	assert.True(t, conf.UseDescriptiveLogger)
	assert.Len(t, conf.HeaderNamesToRedact, 3)
	assert.Equal(t, conf.HeaderNamesToRedact["header-1"], struct{}{})
	assert.Equal(t, conf.HeaderNamesToRedact["header-2"], struct{}{})
	assert.Equal(t, conf.HeaderNamesToRedact["header-3"], struct{}{})
	assert.Len(t, conf.Loaders, 1)
	assert.Len(t, conf.MockFileHandlers, 1)
	assert.IsType(t, &builtInGoTemplate{}, conf.TemplateEngine)
	assert.Len(t, conf.TemplateFunctions, 1)
	assert.NotNil(t, conf.Proxy)
	assert.Equal(t, tlsConfig, conf.TLSConfig)
	assert.NotNil(t, conf.TLSCertificates)
	assert.NotNil(t, conf.TLSClientCAs)
}

func TestLogLevelString(t *testing.T) {
	assert.Equal(t, LogBasic.String(), "basic")
	assert.Equal(t, LogHeader.String(), "header")
	assert.Equal(t, LogBody.String(), "body")
}

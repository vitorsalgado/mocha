package dzhttp

import (
	"crypto/tls"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"text/template"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/dzhttp/cors"
	"github.com/vitorsalgado/mocha/v3/dzhttp/httpval"
	"github.com/vitorsalgado/mocha/v3/dzstd"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

type testBodyParser struct{}

func (p testBodyParser) CanParse(content string, r *http.Request) bool {
	return content == httpval.MIMETextPlain && r.Header.Get("x-test") == "num"
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
		require.NoError(t, err)
		require.NoError(t, res.Body.Close())
		require.True(t, scoped.HasBeenCalled())
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.Contains(t, m.server.Info().URL, addr)
	})

	t.Run("request body parsers from config should take precedence", func(t *testing.T) {
		m := NewAPI(Setup().RequestBodyParsers(&testBodyParser{}))
		m.MustStart()

		defer m.Close()

		scoped := m.MustMock(Post(matcher.URLPath("/test")).
			Body(matcher.StrictEqual(10)).
			Reply(OK()))

		req, _ := http.NewRequest(http.MethodPost, m.URL()+"/test", strings.NewReader("10"))
		req.Header.Add(httpval.HeaderContentType, httpval.MIMETextPlain)
		req.Header.Add("x-test", "num")

		res, err := client.Do(req)
		require.NoError(t, err)
		require.NoError(t, res.Body.Close())
		require.True(t, scoped.HasBeenCalled())
		require.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("middlewares from config should take precedence", func(t *testing.T) {
		msg := "hello world"
		middleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.Header().Add("intercepted", "true")
					w.Header().Add(httpval.HeaderContentType, httpval.MIMETextPlain)
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
		require.NoError(t, err)
		require.NoError(t, res.Body.Close())
		require.False(t, scoped.HasBeenCalled())
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
		require.Equal(t, "true", res.Header.Get("intercepted"))
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
		require.NoError(t, err)
		require.NoError(t, res.Body.Close())
		require.True(t, scoped.HasBeenCalled())
		require.Equal(t, http.StatusOK, res.StatusCode)
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
		MaxBodyParsingLimit(10).
		NoBodyParsing().
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
		Parameters(dzstd.NewInMemoryParameters()).
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

	require.Equal(t, nm, conf.Name)
	require.Equal(t, addr, conf.Addr)
	require.Equal(t, "test_root_dir", conf.RootDir)
	require.Equal(t, http.StatusNotFound, conf.RequestWasNotMatchedStatusCode)
	require.Len(t, conf.RequestBodyParsers, 1)
	require.EqualValues(t, conf.MaxBodyParsingLimit, 10)
	require.True(t, conf.NoBodyParsing)
	require.Len(t, conf.Middlewares, 0)
	require.Equal(t, &cors.DefaultConfig, conf.CORS)
	require.NotNil(t, conf.HandlerDecorator)
	require.Equal(t, LogBasic, conf.LogVerbosity)
	require.Equal(t, LogLevelInfo, conf.LogLevel)
	require.False(t, conf.LogPretty)
	require.Equal(t, &customLogger, conf.Logger)
	require.Equal(t, int64(100), conf.LogBodyMaxSize)
	require.Equal(t, dzstd.NewInMemoryParameters(), conf.Parameters)
	require.Equal(t, []string{"test", "dev"}, conf.MockFileSearchPatterns)
	require.True(t, conf.Debug)
	require.Len(t, conf.HeaderNamesToRedact, 3)
	require.Equal(t, conf.HeaderNamesToRedact["header-1"], struct{}{})
	require.Equal(t, conf.HeaderNamesToRedact["header-2"], struct{}{})
	require.Equal(t, conf.HeaderNamesToRedact["header-3"], struct{}{})
	require.Len(t, conf.Loaders, 1)
	require.Len(t, conf.MockFileHandlers, 1)
	require.IsType(t, &builtInGoTemplate{}, conf.TemplateEngine)
	require.Len(t, conf.TemplateFunctions, 1)
	require.NotNil(t, conf.Proxy)
	require.Equal(t, tlsConfig, conf.TLSConfig)
	require.NotNil(t, conf.TLSCertificates)
	require.NotNil(t, conf.TLSClientCAs)
}

func TestLogLevelString(t *testing.T) {
	require.Equal(t, LogBasic.String(), "basic")
	require.Equal(t, LogHeader.String(), "header")
	require.Equal(t, LogBody.String(), "body")
}

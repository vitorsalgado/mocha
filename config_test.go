package mocha

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/internal/mimetype"
	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

type testBodyParser struct{}

func (p testBodyParser) CanParse(content string, r *http.Request) bool {
	return content == mimetype.TextPlain && r.Header.Get("x-test") == "num"
}

func (p testBodyParser) Parse(body []byte, _ *http.Request) (any, error) {
	return strconv.Atoi(string(body))
}

type customTestServer struct {
	decorated Server
}

func (s *customTestServer) Setup(config *Config, handler http.Handler) error {
	return s.decorated.Setup(config, handler)
}

func (s *customTestServer) Start() (ServerInfo, error) {
	return s.decorated.Start()
}

func (s *customTestServer) StartTLS() (ServerInfo, error) {
	return s.decorated.StartTLS()
}

func (s *customTestServer) Close() error {
	return s.decorated.Close()
}

func (s *customTestServer) Info() ServerInfo {
	return s.decorated.Info()
}

func TestConfig(t *testing.T) {
	t.Run("should run server with the custom given address", func(t *testing.T) {
		addr := os.Getenv("TEST_CUSTOM_ADDR")
		if addr == "" {
			addr = "127.0.0.1:3000"
		}

		m := New(Configure().Addr(addr)).CloseWithT(t)
		m.MustStart()

		defer m.Close()

		scoped := m.MustMock(
			Getf("/test").
				Reply(OK()))

		req := testutil.Get(m.URL() + "/test")
		res, err := req.Do()

		assert.NoError(t, err)
		assert.True(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Contains(t, m.server.Info().URL, addr)
	})

	t.Run("request body parsers from config should take precedence", func(t *testing.T) {
		m := New(Configure().RequestBodyParsers(&testBodyParser{}))
		m.MustStart()

		defer m.Close()

		scoped := m.MustMock(Post(matcher.URLPath("/test")).
			Body(matcher.StrictEqual(10)).
			Reply(OK()))

		req := testutil.Post(m.URL()+"/test", strings.NewReader("10"))
		req.Header(header.ContentType, mimetype.TextPlain)
		req.Header("x-test", "num")

		res, err := req.Do()

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
					w.Header().Add(header.ContentType, mimetype.TextPlain)
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(msg))
				})
		}

		m := New(Configure().Middlewares(middleware))
		m.MustStart()

		defer m.Close()

		scoped := m.MustMock(
			Get(matcher.URLPath("/test")).
				Reply(OK()))

		req := testutil.Get(m.URL() + "/test")
		res, err := req.Do()

		assert.NoError(t, err)
		assert.False(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "true", res.Header.Get("intercepted"))
	})

	t.Run("configure custom server", func(t *testing.T) {
		m := New(Configure().Server(&customTestServer{decorated: newServer()}))
		m.MustStart()

		defer m.Close()

		scoped := m.MustMock(
			Get(matcher.URLPath("/test")).
				Reply(OK()))

		req := testutil.Get(m.URL() + "/test")
		res, err := req.Do()

		assert.NoError(t, err)
		assert.True(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func TestConfigWithFunctions(t *testing.T) {
	addr := ""
	nm := "test"

	m := New(
		WithName(nm),
		WithAddr(addr),
		WithMockNotFoundStatusCode(http.StatusNotFound),
		WithRequestBodyParsers(&jsonBodyParser{}, &plainTextParser{}),
		WithMiddlewares(),
		WithCORS(&_defaultCORSConfig),
		WithServer(&httpTestServer{}),
		WithHandlerDecorator(func(handler http.Handler) http.Handler { return handler }),
		WithLogLevel(LogInfo),
		WithParams(newInMemoryParameters()),
		WithDirs("test", "dev"),
		WithLoader(&fileLoader{}),
		WithProxy(&ProxyConfig{}, &ProxyConfig{}),
		WithMockFileHandlers(&customMockFileHandler{}),
		WithTemplateEngine(newGoTemplate()),
		WithTemplateFunctions(template.FuncMap{"trim": strings.TrimSpace}),
		WithHTTPClient(func() (*http.Client, error) {
			return nil, nil
		}))
	conf := m.Config()

	assert.Equal(t, nm, conf.Name)
	assert.Equal(t, addr, conf.Addr)
	assert.Equal(t, http.StatusNotFound, conf.MockNotFoundStatusCode)
	assert.Len(t, conf.RequestBodyParsers, 2)
	assert.Len(t, conf.Middlewares, 0)
	assert.Equal(t, &_defaultCORSConfig, conf.CORS)
	assert.NotNil(t, conf.HandlerDecorator)
	assert.Equal(t, LogInfo, conf.LogLevel)
	assert.Equal(t, newInMemoryParameters(), conf.Parameters)
	assert.Equal(t, []string{ConfigMockFilePattern, "test", "dev"}, conf.Directories)
	assert.Len(t, conf.Loaders, 1)
	assert.Len(t, conf.MockFileHandlers, 1)
	assert.NotNil(t, conf.Proxy)
	assert.IsType(t, &builtInGoTemplate{}, conf.TemplateEngine)
	assert.Len(t, conf.TemplateFunctions, 1)
	assert.NotNil(t, conf.HTTPClientFactory)
}

func TestConfigBuilder(t *testing.T) {
	addr := ""
	nm := "test"

	m := New(Configure().
		Name(nm).
		Addr(addr).
		MockNotFoundStatusCode(http.StatusNotFound).
		RequestBodyParsers(&jsonBodyParser{}, &plainTextParser{}).
		Middlewares().
		CORS(&_defaultCORSConfig).
		Server(&httpTestServer{}).
		HandlerDecorator(func(handler http.Handler) http.Handler { return handler }).
		LogLevel(LogInfo).
		Parameters(newInMemoryParameters()).
		Dirs("test", "dev").
		Loader(&fileLoader{}).
		MockFileHandlers(&customMockFileHandler{}).
		TemplateEngine(newGoTemplate()).
		BuiltInTemplateEngineFunctions(template.FuncMap{"trim": strings.TrimSpace}).
		Proxy(&ProxyConfig{}, &ProxyConfig{}))
	conf := m.Config()

	assert.Equal(t, nm, conf.Name)
	assert.Equal(t, addr, conf.Addr)
	assert.Equal(t, http.StatusNotFound, conf.MockNotFoundStatusCode)
	assert.Len(t, conf.RequestBodyParsers, 2)
	assert.Len(t, conf.Middlewares, 0)
	assert.Equal(t, &_defaultCORSConfig, conf.CORS)
	assert.NotNil(t, conf.HandlerDecorator)
	assert.Equal(t, LogInfo, conf.LogLevel)
	assert.Equal(t, newInMemoryParameters(), conf.Parameters)
	assert.Equal(t, []string{ConfigMockFilePattern, "test", "dev"}, conf.Directories)
	assert.Len(t, conf.Loaders, 1)
	assert.Len(t, conf.MockFileHandlers, 1)
	assert.IsType(t, &builtInGoTemplate{}, conf.TemplateEngine)
	assert.Len(t, conf.TemplateFunctions, 1)
	assert.NotNil(t, conf.Proxy)
}

func TestWithNewFiles(t *testing.T) {
	m := New(WithNewDirs("test", "dev"))

	assert.Equal(t, []string{"test", "dev"}, m.config.Directories)
}

func TestUseColors(t *testing.T) {
	SetColors(false)
	SetColors(true)
}

func TestLogLevelString(t *testing.T) {
	assert.Equal(t, LogSilently.String(), "silent")
	assert.Equal(t, LogInfo.String(), "info")
	assert.Equal(t, LogVerbose.String(), "verbose")
}

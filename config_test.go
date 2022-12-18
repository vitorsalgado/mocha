package mocha

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/internal/mimetype"
	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/reply"
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

		m := New(t, Configure().Addr(addr)).CloseWithT(t)
		m.MustStart()

		defer m.Close()

		scoped := m.MustMock(
			Getf("/test").
				Reply(reply.OK()))

		req := testutil.Get(m.URL() + "/test")
		res, err := req.Do()

		assert.NoError(t, err)
		assert.True(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Contains(t, m.server.Info().URL, addr)
	})

	t.Run("request body parsers from config should take precedence", func(t *testing.T) {
		m := New(t, Configure().RequestBodyParsers(&testBodyParser{}))
		m.MustStart()

		defer m.Close()

		scoped := m.MustMock(Post(matcher.URLPath("/test")).
			Body(matcher.Equal(10)).
			Reply(reply.OK()))

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

		m := New(t, Configure().Middlewares(middleware))
		m.MustStart()

		defer m.Close()

		scoped := m.MustMock(
			Get(matcher.URLPath("/test")).
				Reply(reply.OK()))

		req := testutil.Get(m.URL() + "/test")
		res, err := req.Do()

		assert.NoError(t, err)
		assert.False(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "true", res.Header.Get("intercepted"))
	})

	t.Run("configure custom server", func(t *testing.T) {
		m := New(t, Configure().Server(&customTestServer{decorated: newServer()}))
		m.MustStart()

		defer m.Close()

		scoped := m.MustMock(
			Get(matcher.URLPath("/test")).
				Reply(reply.OK()))

		req := testutil.Get(m.URL() + "/test")
		res, err := req.Do()

		assert.NoError(t, err)
		assert.True(t, scoped.HasBeenCalled())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func TestConfig_WithFunctions(t *testing.T) {
	addr := ":3000"
	nm := "test"

	m := New(t,
		WithName(nm),
		WithAddr(addr),
		WithRequestBodyParsers(&jsonBodyParser{}, &plainTextParser{}),
		WithMiddlewares(),
		WithCORS(&_defaultCORSConfig),
		WithServer(&httpTestServer{}),
		WithHandlerDecorator(func(handler http.Handler) http.Handler { return handler }),
		WithLogLevel(LogInfo),
		WithParams(reply.Parameters()),
		WithDirs("test", "dev"),
		WithLoader(&FileLoader{}),
		WithProxy(&ProxyConfig{}, &ProxyConfig{}))
	conf := m.Config()

	assert.Equal(t, nm, conf.Name)
	assert.Equal(t, addr, conf.Addr)
	assert.Len(t, conf.RequestBodyParsers, 2)
	assert.Len(t, conf.Middlewares, 0)
	assert.Equal(t, &_defaultCORSConfig, conf.CORS)
	assert.NotNil(t, conf.HandlerDecorator)
	assert.Equal(t, LogInfo, conf.LogLevel)
	assert.Equal(t, reply.Parameters(), conf.Parameters)
	assert.Equal(t, []string{ConfigMockFilePattern, "test", "dev"}, conf.Directories)
	assert.Len(t, conf.Loaders, 1)
	assert.NotNil(t, conf.Proxy)
}

func TestWithNewFiles(t *testing.T) {
	m := New(t, WithNewDirs("test", "dev"))

	assert.Equal(t, []string{"test", "dev"}, m.config.Directories)
}

func TestUseColors(t *testing.T) {
	UseColors(false)
	UseColors(true)
}

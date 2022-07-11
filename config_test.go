package mocha

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/expect"
	"github.com/vitorsalgado/mocha/internal/testutil"
	"github.com/vitorsalgado/mocha/reply"
	"github.com/vitorsalgado/mocha/x/headers"
	"github.com/vitorsalgado/mocha/x/mimetypes"
)

type testBodyParser struct{}

func (p testBodyParser) CanParse(content string, r *http.Request) bool {
	return content == mimetypes.TextPlain && r.Header.Get("x-test") == "num"
}

func (p testBodyParser) Parse(body []byte, _ *http.Request) (any, error) {
	return strconv.Atoi(string(body))
}

type customTestServer struct {
	decorated Server
}

func (s *customTestServer) Configure(config Config, handler http.Handler) error {
	return s.decorated.Configure(config, handler)
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

		m := New(t, Configure().
			Addr(addr).
			Build())
		m.Start()

		scoped := m.Mock(
			Get(expect.URLPath("/test")).
				Reply(reply.OK()))

		req := testutil.Get(m.URL() + "/test")
		res, err := req.Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.True(t, scoped.Called())
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Contains(t, m.server.Info().URL, addr)
	})

	t.Run("request body parsers from config should take precedence", func(t *testing.T) {
		m := New(t, Configure().
			RequestBodyParsers(&testBodyParser{}).
			Build())
		m.Start()

		scoped := m.Mock(Post(expect.URLPath("/test")).
			Body(expect.ToEqual[any](10)).
			Reply(reply.OK()))

		req := testutil.Post(m.URL()+"/test", strings.NewReader("10"))
		req.Header(headers.ContentType, mimetypes.TextPlain)
		req.Header("x-test", "num")

		res, err := req.Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.True(t, scoped.Called())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("middlewares from config should take precedence", func(t *testing.T) {
		msg := "hello world"
		middleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.Header().Add("intercepted", "true")
					w.Header().Add(headers.ContentType, mimetypes.TextPlain)
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(msg))
				})
		}

		m := New(t, Configure().
			Middlewares(middleware).
			Build())
		m.Start()

		scoped := m.Mock(
			Get(expect.URLPath("/test")).
				Reply(reply.OK()))

		req := testutil.Get(m.URL() + "/test")
		res, err := req.Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.False(t, scoped.Called())
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, "true", res.Header.Get("intercepted"))
	})

	t.Run("configure custom server", func(t *testing.T) {
		m := New(t, Configure().
			Server(&customTestServer{decorated: newServer()}).
			Build())
		m.Start()

		scoped := m.Mock(
			Get(expect.URLPath("/test")).
				Reply(reply.OK()))

		req := testutil.Get(m.URL() + "/test")
		res, err := req.Do()
		if err != nil {
			log.Fatal(err)
		}

		assert.True(t, scoped.Called())
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

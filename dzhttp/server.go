package dzhttp

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
)

var _ Server = (*httpTestServer)(nil)

// ServerInfo holds HTTP server information, like its URL.
type ServerInfo struct {
	URL string
}

// Server defines HTTP mock server operations.
type Server interface {
	// Setup configures the HTTP mock.
	// It is the first method called during initialization.
	Setup(*HTTPMockApp, http.Handler) error

	// Start starts a server.
	Start() error

	// StartTLS starts a server with TLS.
	StartTLS() error

	// Close closes the server.
	Close() error

	// CloseNow force close the server.
	CloseNow() error

	// Info returns server information.
	Info() *ServerInfo

	// S returns the server implementation that is being used by this component.
	// For example, the built-in Server will return a *httptest.Server.
	S() any
}

type httpTestServer struct {
	app       *HTTPMockApp
	handler   http.Handler
	server    *httptest.Server
	info      *ServerInfo
	needSetup bool
	rwMutex   sync.RWMutex
}

func newServer() Server {
	return &httpTestServer{info: &ServerInfo{}}
}

func (s *httpTestServer) Setup(app *HTTPMockApp, handler http.Handler) error {
	r := chi.NewRouter()
	r.Mount(app.config.AdminPath, app.admin.Init())
	r.Handle("/*", handler)

	s.app = app
	s.handler = handler
	s.server = httptest.NewUnstartedServer(r)

	if app.config.Addr != "" {
		addr := app.config.Addr
		_, err := strconv.Atoi(addr)
		if err == nil {
			addr = ":" + addr
		}

		err = s.server.Listener.Close()
		if err != nil {
			return err
		}

		listener, err := net.Listen("tcp", addr)
		if err != nil {
			return err
		}

		s.server.Listener = listener
	}

	if app.config.TLSConfig != nil {
		s.server.TLS = app.config.TLSConfig
		s.server.EnableHTTP2 = app.config.UseHTTP2
	} else {
		s.server.TLS = &tls.Config{Certificates: app.config.TLSCertificates}

		if app.config.TLSClientCAs != nil {
			s.server.TLS.ClientCAs = app.config.TLSClientCAs
			s.server.TLS.ClientAuth = tls.RequireAndVerifyClientCert
		}

		s.server.EnableHTTP2 = app.config.UseHTTP2
	}

	return nil
}

func (s *httpTestServer) Start() error {
	err := s.beforeStart()
	if err != nil {
		return err
	}

	s.server.Start()
	s.info.URL = s.server.URL

	return nil
}

func (s *httpTestServer) StartTLS() error {
	err := s.beforeStart()
	if err != nil {
		return err
	}

	s.server.EnableHTTP2 = true
	s.server.StartTLS()
	s.info.URL = s.server.URL

	return nil
}

func (s *httpTestServer) Close() error {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	s.server.CloseClientConnections()
	s.server.Close()
	s.needSetup = true

	return nil
}

func (s *httpTestServer) CloseNow() error {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	s.server.CloseClientConnections()
	s.server.Close()
	s.needSetup = true

	return nil
}

func (s *httpTestServer) Info() *ServerInfo {
	return s.info
}

func (s *httpTestServer) S() any {
	return s.server
}

func (s *httpTestServer) beforeStart() error {
	if s.needSetup {
		return s.Setup(s.app, s.handler)
	}

	return nil
}

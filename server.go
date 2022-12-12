package mocha

import (
	"net"
	"net/http"
	"net/http/httptest"
)

// ServerInfo holds HTTP server information, like its URL.
type ServerInfo struct {
	URL string
}

// Server defines HTTP mock server operations.
type Server interface {
	// Setup configures the HTTP mock.
	// It is the first method called during initialization.
	Setup(*Config, http.Handler) error

	// Start starts a server.
	Start() (ServerInfo, error)

	// StartTLS starts a TLS server.
	StartTLS() (ServerInfo, error)

	// Close closes the server.
	Close() error

	// Info returns server information.
	Info() ServerInfo
}

type httpTestServer struct {
	server *httptest.Server
	info   ServerInfo
}

func newServer() Server {
	return &httpTestServer{info: ServerInfo{}}
}

func (s *httpTestServer) Setup(config *Config, handler http.Handler) error {
	s.server = httptest.NewUnstartedServer(handler)
	s.server.EnableHTTP2 = true

	if config.Addr != "" {
		err := s.server.Listener.Close()
		if err != nil {
			return err
		}

		listener, err := net.Listen("tcp", config.Addr)
		if err != nil {
			return err
		}

		s.server.Listener = listener
	}

	return nil
}

func (s *httpTestServer) Start() (ServerInfo, error) {
	s.server.Start()
	s.info.URL = s.server.URL

	return s.info, nil
}

func (s *httpTestServer) StartTLS() (ServerInfo, error) {
	s.server.StartTLS()
	s.info.URL = s.server.URL

	return s.info, nil
}

func (s *httpTestServer) Close() error {
	s.server.Close()
	return nil
}

func (s *httpTestServer) Info() ServerInfo {
	return s.info
}

package mocha

import (
	"net"
	"net/http"
	"net/http/httptest"

	"github.com/vitorsalgado/mocha/core"
)

type (
	// ServerInfo holds HTTP server information, like its URL.
	ServerInfo struct {
		URL string
	}

	// Server defines HTTP mock server operations.
	Server interface {
		// Configure configures the HTTP mock.
		// It is the first method called by Mocha during initialization.
		Configure(Config, http.Handler) error

		// Start starts a server.
		Start() (ServerInfo, error)

		// StartTLS starts a TLS server.
		StartTLS() (ServerInfo, error)

		// Close the server.
		Close() error

		// Info returns server information.
		Info() ServerInfo
	}

	testServer struct {
		server *httptest.Server
		info   ServerInfo
		t      core.T
	}
)

func newServer() Server {
	return &testServer{info: ServerInfo{}}
}

func (s *testServer) Configure(config Config, handler http.Handler) error {
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

func (s *testServer) Start() (ServerInfo, error) {
	s.server.Start()
	s.info.URL = s.server.URL

	return s.info, nil
}

func (s *testServer) StartTLS() (ServerInfo, error) {
	s.server.StartTLS()
	s.info.URL = s.server.URL

	return s.info, nil
}

func (s *testServer) Close() error {
	s.server.Close()
	return nil
}

func (s *testServer) Info() ServerInfo {
	return s.info
}

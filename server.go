package mocha

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
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
	Setup(*Mocha, http.Handler) error

	// Start starts a server.
	Start() error

	// StartTLS starts a server with TLS.
	StartTLS() error

	// Close closes the server.
	Close() error

	// Info returns server information.
	Info() *ServerInfo
}

type httpTestServer struct {
	app       *Mocha
	handler   http.Handler
	server    *httptest.Server
	info      *ServerInfo
	needSetup bool
}

func newServer() Server {
	return &httpTestServer{info: &ServerInfo{}}
}

func (s *httpTestServer) Setup(app *Mocha, handler http.Handler) error {
	s.app = app
	s.handler = handler
	s.server = httptest.NewUnstartedServer(handler)

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
	} else {
		if (app.config.TLSCertificateFs == "" || app.config.TLSKeyFs == "") && app.config.TLSRootCAs == nil {
			return nil
		}

		var certs []tls.Certificate

		if app.config.TLSCertificateFs != "" && app.config.TLSKeyFs != "" {
			cert, err := tls.LoadX509KeyPair(app.config.TLSCertificateFs, app.config.TLSKeyFs)
			if err != nil {
				return err
			}

			certs = make([]tls.Certificate, 1)
			certs[0] = cert
		}

		s.server.TLS = &tls.Config{Certificates: certs, RootCAs: app.config.TLSRootCAs}
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
	s.server.Close()
	s.needSetup = true

	return nil
}

func (s *httpTestServer) Info() *ServerInfo {
	return s.info
}

func (s *httpTestServer) beforeStart() error {
	if s.needSetup {
		return s.Setup(s.app, s.handler)
	}

	return nil
}

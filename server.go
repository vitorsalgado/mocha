package mocha

import (
	"net/http"
	"net/http/httptest"
)

type (
	ServerInfo struct {
		URL string
	}

	Server interface {
		Start() ServerInfo
		Close()
	}

	ServerBuilder interface {
		Build(root http.Handler) Server
	}

	HTTPTestServerBuilder struct {
	}

	StandardServer struct {
		server *httptest.Server
	}
)

func (s StandardServer) Start() ServerInfo {
	s.server.Start()
	return ServerInfo{}
}

func (b HTTPTestServerBuilder) Build() Server {
	return nil
}

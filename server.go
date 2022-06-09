package mocha

import (
	"net/http/httptest"
)

type Server struct {
}

func (s Server) Start() *httptest.Server {
	return httptest.NewServer(&Handler{})
}

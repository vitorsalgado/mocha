package mocha

import (
	"net/http"
	"net/url"
)

type MockRequest struct {
	Method string
	URL    url.URL
	Host   string
	Header http.Header
	Body   []byte
}

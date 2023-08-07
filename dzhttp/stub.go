package dzhttp

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"

	"github.com/vitorsalgado/mocha/v3/dzhttp/internal/httprec"
)

// MockedResponse defines the HTTP response that will be served once a Mock is matched for an HTTP Request.
type MockedResponse struct {
	StatusCode int
	Header     http.Header
	Cookies    []*http.Cookie
	Body       []byte
	Trailer    http.Header
	Encoding   string
}

func newStub() *MockedResponse {
	return &MockedResponse{Cookies: make([]*http.Cookie, 0), Header: make(http.Header), Trailer: make(http.Header)}
}

// Gunzip decompresses Gzip body.
func (s *MockedResponse) Gunzip() ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(s.Body))
	if err != nil {
		return nil, err
	}

	defer gz.Close()

	return io.ReadAll(gz)
}

func newResponse(w http.ResponseWriter, stub *MockedResponse) error {
	rw := w.(*httprec.HTTPRec)
	result := rw.Result()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return err
	}

	defer result.Body.Close()

	if len(w.Header()) != len(result.Header) {
		for k, v := range w.Header() {
			if result.Header.Get(k) == "" {
				for _, vv := range v {
					result.Header.Add(k, vv)
				}
			}
		}
	}

	stub.StatusCode = result.StatusCode
	stub.Header = result.Header.Clone()
	stub.Cookies = result.Cookies()

	if len(body) > 0 {
		stub.Body = body
	}

	return nil
}

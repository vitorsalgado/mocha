package dzhttp

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"

	"github.com/vitorsalgado/mocha/v3/dzhttp/internal/httprec"
)

// Stub defines the HTTP response that will be served once a Mock is matched for an HTTP Request.
type Stub struct {
	StatusCode int
	Header     http.Header
	Cookies    []*http.Cookie
	Body       []byte
	Trailer    http.Header
	Encoding   string
}

func newStub() *Stub {
	return &Stub{Cookies: make([]*http.Cookie, 0), Header: make(http.Header), Trailer: make(http.Header)}
}

// Gunzip decompresses Gzip body.
func (s *Stub) Gunzip() ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(s.Body))
	if err != nil {
		return nil, err
	}

	defer gz.Close()

	return io.ReadAll(gz)
}

func newResponseStub(w http.ResponseWriter, stub *Stub) error {
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
package mocha

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/vitorsalgado/mocha/v3/internal/headers"
	"github.com/vitorsalgado/mocha/v3/internal/mimetypes"
)

// RequestBodyParser parses request body if CanParse returns true.
// Multiple implementations of RequestBodyParser can be provided to Mocha using options.
type RequestBodyParser interface {
	// CanParse checks if current request body should be parsed by this component.
	// First parameter is the incoming content-type.
	CanParse(contentType string, r *http.Request) bool

	// Parse parses the request body.
	Parse(body []byte, r *http.Request) (any, error)
}

// parseRequestBody tests given parsers until it finds one that can parse the request body.
// User provided RequestBodyParser takes precedence.
func parseRequestBody(r *http.Request, parsers []RequestBodyParser) (any, error) {
	if r.Body != nil && r.Method != http.MethodGet && r.Method != http.MethodHead {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(b))

		contentType := r.Header.Get(headers.ContentType)

		for _, parse := range parsers {
			if parse.CanParse(contentType, r) {
				body, err := parse.Parse(b, r)
				if err != nil {
					return nil, err
				}

				return body, nil
			}
		}
	}

	return nil, nil
}

// jsonBodyParser parses requests with content type header containing "application/json"
type jsonBodyParser struct{}

func (parser *jsonBodyParser) CanParse(content string, _ *http.Request) bool {
	return strings.Contains(content, mimetypes.JSON)
}

func (parser *jsonBodyParser) Parse(body []byte, _ *http.Request) (data any, err error) {
	err = json.Unmarshal(body, &data)
	return data, err
}

// jsonBodyParser parses requests with content type header containing "application/x-www-form-urlencoded"
type formURLEncodedParser struct{}

func (parser *formURLEncodedParser) CanParse(content string, _ *http.Request) bool {
	return strings.Contains(content, mimetypes.FormURLEncoded)
}

func (parser *formURLEncodedParser) Parse(_ []byte, r *http.Request) (any, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	return r.Form, nil
}

// plainTextParser parses requests with content type header containing "text/plain"
type plainTextParser struct{}

func (parser *plainTextParser) CanParse(content string, _ *http.Request) bool {
	return strings.Contains(content, mimetypes.TextPlain)
}

func (parser *plainTextParser) Parse(body []byte, _ *http.Request) (any, error) {
	return string(body), nil
}

// bytesParser is default parser when none can parse.
type bytesParser struct{}

func (parser *bytesParser) CanParse(_ string, _ *http.Request) bool {
	return true
}

func (parser *bytesParser) Parse(body []byte, _ *http.Request) (any, error) {
	return body, nil
}

package mocha

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/vitorsalgado/mocha/internal/header"
	"github.com/vitorsalgado/mocha/internal/mime"
)

type (
	// RequestBodyParser parses request body if CanParse returns true.
	// Multiple implementations of RequestBodyParser can be provided to Mocha using options.
	RequestBodyParser interface {
		// CanParse checks if current request body should be parsed by this component.
		CanParse(content string, r *http.Request) bool

		// Parse parses the request body.
		Parse(r *http.Request) (any, error)
	}
)

// parseRequestBody tests given parsers until it finds one that can parse the request body.
// User provided RequestBodyParser takes precedence.
func parseRequestBody(r *http.Request, parsers []RequestBodyParser) (any, error) {
	if r.Body != nil && r.Method != http.MethodGet && r.Method != http.MethodHead {
		var content = r.Header.Get(header.ContentType)

		for _, parse := range parsers {
			if parse.CanParse(content, r) {
				body, err := parse.Parse(r)
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

func (parser jsonBodyParser) CanParse(content string, _ *http.Request) bool {
	return strings.Contains(content, mime.ContentTypeJSON)
}

func (parser jsonBodyParser) Parse(r *http.Request) (data any, err error) {
	err = json.NewDecoder(r.Body).Decode(&data)
	return data, err
}

// jsonBodyParser parses requests with content type header containing "application/x-www-form-urlencoded"
type formURLEncodedParser struct{}

func (parser formURLEncodedParser) CanParse(content string, _ *http.Request) bool {
	return strings.Contains(content, mime.ContentTypeFormURLEncoded)
}

func (parser *formURLEncodedParser) Parse(r *http.Request) (any, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	return r.Form.Encode(), nil
}

// plainTextParser parses requests with content type header containing "text/plain"
type plainTextParser struct{}

func (parser *plainTextParser) CanParse(content string, _ *http.Request) bool {
	return strings.Contains(content, mime.ContentTypeTextPlain)
}

func (parser *plainTextParser) Parse(r *http.Request) (any, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	return r.Form.Encode(), nil
}

// bytesParser is default parser when none can parse.
type bytesParser struct{}

func (parser *bytesParser) CanParse(content string, _ *http.Request) bool {
	return strings.Contains(content, mime.ContentTypeTextPlain)
}

func (parser *bytesParser) Parse(r *http.Request) (any, error) {
	return ioutil.ReadAll(r.Body)
}

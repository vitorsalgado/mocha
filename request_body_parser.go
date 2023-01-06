package mocha

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/internal/mimetype"
)

// RequestBodyParser parses request body if CanParse returns true.
// Multiple implementations of RequestBodyParser can be provided to Mocha using options.
type RequestBodyParser interface {
	// CanParse checks if current request body should be parsed by this component.
	// first parameter is the incoming content-type.
	CanParse(contentType string, r *http.Request) bool

	// Parse parses the request body.
	Parse(body []byte, r *http.Request) (any, error)
}

// parseRequestBody tests given parsers until it finds one that can parse the request body.
// User provided RequestBodyParser takes precedence.
func parseRequestBody(r *http.Request, parsers []RequestBodyParser) (parsedBody any, rawBody []byte, err error) {
	if r.Body != nil && r.Method != http.MethodGet && r.Method != http.MethodHead {
		rawBody, err = io.ReadAll(r.Body)
		if err != nil {
			return nil, nil, err
		}

		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(rawBody))

		contentType := r.Header.Get(header.ContentType)

		for _, parse := range parsers {
			if parse.CanParse(contentType, r) {
				parsedBody, err = parse.Parse(rawBody, r)
				if err != nil {
					return nil, rawBody, err
				}

				return parsedBody, rawBody, nil
			}
		}
	}

	return nil, nil, nil
}

// jsonBodyParser parses requests with content type header containing "application/json"
type jsonBodyParser struct{}

func (parser *jsonBodyParser) CanParse(content string, _ *http.Request) bool {
	return strings.Contains(content, mimetype.JSON)
}

func (parser *jsonBodyParser) Parse(body []byte, _ *http.Request) (data any, err error) {
	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&data)
	return data, err
}

// jsonBodyParser parses requests with content type header containing "application/x-www-form-urlencoded"
type formURLEncodedParser struct{}

func (parser *formURLEncodedParser) CanParse(content string, _ *http.Request) bool {
	return strings.Contains(content, mimetype.FormURLEncoded)
}

func (parser *formURLEncodedParser) Parse(_ []byte, r *http.Request) (any, error) {
	err := r.ParseForm()
	return r.Form, err
}

// plainTextParser parses requests with content type header containing "text/plain"
type plainTextParser struct{}

func (parser *plainTextParser) CanParse(content string, _ *http.Request) bool {
	return strings.Contains(content, mimetype.TextPlain)
}

func (parser *plainTextParser) Parse(body []byte, _ *http.Request) (any, error) {
	return string(body), nil
}

// noopParser is default parser and runs when none is selected.
// It basically returns the body []byte.
type noopParser struct{}

func (parser *noopParser) CanParse(_ string, _ *http.Request) bool         { return true }
func (parser *noopParser) Parse(body []byte, _ *http.Request) (any, error) { return body, nil }

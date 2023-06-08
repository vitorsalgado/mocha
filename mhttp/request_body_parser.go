package mhttp

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/vitorsalgado/mocha/v3/mhttpv"
)

// RequestBodyParser parses the request body if CanParse returns true.
// Multiple implementations of RequestBodyParser can be provided via configuration.
type RequestBodyParser interface {
	// CanParse checks if the current request body should be parsed.
	// The first parameter is the incoming content type.
	CanParse(contentType string, r *http.Request) bool

	// Parse parses the request body.
	Parse(body []byte, r *http.Request) (any, error)
}

// parseRequestBody tests given parsers until it finds one that can parse the request body.
// The user provided RequestBodyParser takes precedence.
func parseRequestBody(r *http.Request, parsers []RequestBodyParser) (parsedBody any, rawBody []byte, err error) {
	if r.Body == nil {
		return
	}

	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		return
	}

	rawBody, err = io.ReadAll(r.Body)
	if err != nil {
		return nil, nil, err
	}

	if len(rawBody) == 0 {
		return nil, nil, nil
	}

	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	contentType := r.Header.Get(mhttpv.HeaderContentType)

	for _, parser := range parsers {
		if parser.CanParse(contentType, r) {
			parsedBody, err = parser.Parse(rawBody, r)
			if err != nil {
				return nil, rawBody,
					fmt.Errorf("%T failed: %w", parser, err)
			}

			return parsedBody, rawBody, nil
		}
	}

	return nil, nil, nil
}

// jsonBodyParser parses requests with content type header containing "application/x-www-form-urlencoded"
type formURLEncodedParser struct{}

func (parser *formURLEncodedParser) CanParse(content string, _ *http.Request) bool {
	return strings.Contains(content, mhttpv.MIMEFormURLEncoded)
}

func (parser *formURLEncodedParser) Parse(_ []byte, r *http.Request) (any, error) {
	err := r.ParseForm()
	return r.Form, err
}

// plainTextParser parses requests with content type header containing "text/plain"
type plainTextParser struct{}

func (parser *plainTextParser) CanParse(content string, _ *http.Request) bool {
	return strings.Contains(content, mhttpv.MIMETextPlain) || strings.Contains(content, mhttpv.MIMEApplicationJSON)
}

func (parser *plainTextParser) Parse(body []byte, _ *http.Request) (any, error) {
	return string(body), nil
}

// noopParser is the default parser and runs when none is selected.
// It just returns the body []byte.
type noopParser struct{}

func (parser *noopParser) CanParse(_ string, _ *http.Request) bool         { return true }
func (parser *noopParser) Parse(body []byte, _ *http.Request) (any, error) { return body, nil }

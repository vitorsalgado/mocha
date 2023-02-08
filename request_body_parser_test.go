package mocha

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/internal/mimetype"
)

func TestWithRequestBodyParsers_CanParse(t *testing.T) {
	jsonParser := &jsonBodyParser{}
	formParser := &formURLEncodedParser{}
	textParser := &plainTextParser{}
	noop := &noopParser{}

	newReq := func(header map[string]string) *http.Request {
		req, err := http.NewRequest(http.MethodPost, "https://localhost:8080", nil)
		require.NoError(t, err)

		for k, v := range header {
			req.Header.Add(k, v)
		}

		return req
	}

	testCases := []struct {
		name        string
		contentType string
		parser      RequestBodyParser
		req         *http.Request
		expected    bool
	}{
		{"JSON: can parse application/json", mimetype.JSON, jsonParser, newReq(map[string]string{header.ContentType: mimetype.JSON}), true},
		{"JSON: can parse application/json; charset=UTF-8", mimetype.JSONCharsetUTF8, jsonParser, newReq(map[string]string{header.ContentType: mimetype.JSONCharsetUTF8}), true},
		{"JSON: should not parse", mimetype.TextPlain, jsonParser, newReq(map[string]string{header.ContentType: mimetype.TextPlain}), false},

		{"Form: can parse application/form-url-encoded", mimetype.FormURLEncoded, formParser, newReq(map[string]string{header.ContentType: mimetype.FormURLEncoded}), true},
		{"Form: can parse application/form-url-encoded; charset=UTF-8", mimetype.FormURLEncodedCharsetUTF8, formParser, newReq(map[string]string{header.ContentType: mimetype.FormURLEncodedCharsetUTF8}), true},
		{"Form: should not parse", mimetype.JSON, formParser, newReq(map[string]string{header.ContentType: mimetype.JSON}), false},

		{"Text: can parse text/plain", mimetype.TextPlain, textParser, newReq(map[string]string{header.ContentType: mimetype.TextPlain}), true},
		{"Text: can parse text/plain; charset=UTF-8", mimetype.TextPlainCharsetUTF8, textParser, newReq(map[string]string{header.ContentType: mimetype.TextPlainCharsetUTF8}), true},
		{"Text: should not parse", mimetype.JSON, textParser, newReq(map[string]string{header.ContentType: mimetype.JSON}), false},

		{"Noop: should always parse", mimetype.TextPlain, noop, newReq(map[string]string{header.ContentType: mimetype.TextPlain}), true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.parser.CanParse(tc.contentType, tc.req))
		})
	}
}

func TestWithRequestBodyParsers_Parse(t *testing.T) {
	jsonParser := &jsonBodyParser{}
	formParser := &formURLEncodedParser{}
	textParser := &plainTextParser{}
	noop := &noopParser{}

	req, err := http.NewRequest(http.MethodPost, "https://localhost:8080", nil)
	require.NoError(t, err)

	testCases := []struct {
		name     string
		parser   RequestBodyParser
		req      *http.Request
		body     []byte
		expected any
	}{
		{"json", jsonParser, req, []byte(`{"test":"ok"}`), map[string]any{"test": "ok"}},
		{"form", formParser, func(r *http.Request) *http.Request {
			req, err := http.NewRequest(http.MethodPost, "https://localhost:8080", strings.NewReader("test=ok"))
			req.Header.Add(header.ContentType, mimetype.FormURLEncoded)
			require.NoError(t, err)

			return req
		}(req), []byte(`test=ok`), url.Values{"test": []string{"ok"}}},
		{"text", textParser, req, []byte(`hello world`), "hello world"},
		{"noop", noop, req, []byte(`hello world`), []byte(`hello world`)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := tc.parser.Parse(tc.body, tc.req)

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, b)
		})
	}
}

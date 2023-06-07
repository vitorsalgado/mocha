package mhttp

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/misc"
)

func TestWithRequestBodyParsersCanParse(t *testing.T) {
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
		{"Form: can parse application/form-url-encoded", misc.MIMEFormURLEncoded, formParser, newReq(map[string]string{misc.HeaderContentType: misc.MIMEFormURLEncoded}), true},
		{"Form: can parse application/form-url-encoded; charset=UTF-8", misc.MIMEFormURLEncodedCharsetUTF8, formParser, newReq(map[string]string{misc.HeaderContentType: misc.MIMEFormURLEncodedCharsetUTF8}), true},
		{"Form: should not parse", misc.MIMEApplicationJSON, formParser, newReq(map[string]string{misc.HeaderContentType: misc.MIMEApplicationJSON}), false},

		{"Text: can parse text/plain", misc.MIMETextPlain, textParser, newReq(map[string]string{misc.HeaderContentType: misc.MIMETextPlain}), true},
		{"Text: can parse text/plain; charset=UTF-8", misc.MIMETextPlainCharsetUTF8, textParser, newReq(map[string]string{misc.HeaderContentType: misc.MIMETextPlainCharsetUTF8}), true},
		{"Text: should parse JSON as text", misc.MIMEApplicationJSON, textParser, newReq(map[string]string{misc.HeaderContentType: misc.MIMEApplicationJSON}), true},

		{"Noop: should always parse", misc.MIMETextPlain, noop, newReq(map[string]string{misc.HeaderContentType: misc.MIMETextPlain}), true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.parser.CanParse(tc.contentType, tc.req))
		})
	}
}

func TestWithRequestBodyParsersParse(t *testing.T) {
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
		{"form", formParser, func(r *http.Request) *http.Request {
			req, err := http.NewRequest(http.MethodPost, "https://localhost:8080", strings.NewReader("test=ok"))
			req.Header.Add(misc.HeaderContentType, misc.MIMEFormURLEncoded)
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

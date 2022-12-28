package types

import (
	"net/http"
	"net/url"
)

// RequestValues groups HTTP request data, including the parsed body, if any.
type RequestValues struct {
	// RawRequest is the original incoming http.Request.
	RawRequest *http.Request

	// URL is full request url.URL, including scheme, host, port.
	URL *url.URL

	// ParsedBody is the parsed http.Request body.
	ParsedBody any
}

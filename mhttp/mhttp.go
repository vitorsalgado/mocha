package mhttp

import (
	"net/http"
	"net/url"
	"time"
)

// RequestValues groups HTTP request data, including the parsed body.
// It is used by several components during the request matching phase.
type RequestValues struct {
	// StartedAt indicates when the request arrived.
	StartedAt time.Time

	// RawRequest is the original incoming http.Request.
	RawRequest *http.Request

	// URL is the full request url.URL, including scheme, host, port.
	URL *url.URL

	// URLPathSegments stores path segments.
	// Eg.: /test/100 -> []string{"test", "100"}
	URLPathSegments []string

	// RawBody is the HTTP request body bytes.
	RawBody []byte

	// ParsedBody is the parsed http.Request body parsed by a RequestBodyParser instance.
	// It'll be nil if the HTTP request does not contain a body.
	ParsedBody any

	// App exposes application instance associated with the incoming HTTP request.
	App *HTTPMockApp

	// Mock is the matched Mock for the current HTTP request.
	Mock *HTTPMock
}

// CallbackInput represents the arguments that will be passed to every Callback implementation
type CallbackInput struct {
	// RawRequest is the original incoming http.Request.
	RawRequest *http.Request

	// URL is the full request url.URL, including scheme, host, port.
	URL *url.URL

	// ParsedBody is the parsed http.Request body.
	ParsedBody any

	// App exposes the application instance associated with the incoming HTTP request.
	App *HTTPMockApp

	// Mock is the matched Mock for the current HTTP request.
	Mock *HTTPMock

	// Stub is the HTTP response Stub served.
	Stub *Stub
}

// Callback defines the contract for an action that will be executed after serving a mocked HTTP response.
type Callback func(input *CallbackInput) error

// PostActionInput represents the arguments that will be passed to every Callback implementation
type PostActionInput struct {
	// RawRequest is the original incoming http.Request.
	RawRequest *http.Request

	// URL is the full request url.URL, including scheme, host, port.
	URL *url.URL

	// ParsedBody is the parsed http.Request body.
	ParsedBody any

	// App exposes the application instance associated with the incoming HTTP request.
	App *HTTPMockApp

	// Mock is the matched Mock for the current HTTP request.
	Mock *HTTPMock

	// Stub is the HTTP response Stub served.
	Stub *Stub

	// Args allow passing custom arguments to a Callback.
	Args any
}

// PostAction defines the contract for an action that will be executed after serving a mocked HTTP response.
type PostAction interface {
	// Run runs the Callback implementation.
	Run(input *PostActionInput) error
}

type PostActionDef struct {
	Name          string
	RawParameters any
}

// Mapper is the function definition to be used to map a Mock response Stub before serving it.
type Mapper func(requestValues *RequestValues, res *Stub) error

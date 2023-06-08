package httpd

import (
	"net/http"
	"net/url"
	"time"

	"github.com/vitorsalgado/mocha/v3/lib"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

var _ lib.Mock = (*HTTPMock)(nil)

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

// MockFileHandler defines a custom Mock file configuration handler.
// It lets users define custom fields on mock configuration files that could be handled by a MockFileHandler instance.
// It is also possible to change how built-in fields are handled.
type MockFileHandler interface {
	// Handle handles a Mock configuration file.
	//  Parameter fields
	Handle(fields map[string]any, b *HTTPMockBuilder) error
}

// mockFileData represents the data that is passed to Mock files during template parsing.
type mockFileData struct {
	App *templateAppWrapper
}

// HTTPValueSelector defines a function that will be used to extract the value that will be passed to the associated matcher.
type HTTPValueSelector func(r *HTTPValueSelectorInput) any

type HTTPValueSelectorInput struct {
	// RawRequest is the original incoming http.Request.
	RawRequest *http.Request

	// URL is the full request url.URL, including scheme, host, port.
	URL *url.URL

	Query url.Values

	Form url.Values

	// ParsedBody is the parsed http.Request body.
	ParsedBody any
}

type HTTPExpectation = lib.Expectation[*HTTPValueSelectorInput]

// HTTPMock holds metadata and expectations to be matched against HTTP requests in order to serve mocked responses.
// This is the core entity of this project, and most features work based on it.
type HTTPMock struct {
	*lib.BaseMock

	// Callbacks holds a Callback list to be executed after the Mock was matched and served.
	Callbacks []Callback

	PostActions []*PostActionDef

	// Source describes the source of the mock. E.g.: if it was built from a file,
	// it will contain the filename.
	Source string

	Reply Reply

	// Mappers stores response mappers associated with this Mock.
	Mappers []Mapper

	Pipes []lib.Piping

	after        []matcher.OnAfterMockServed
	expectations []*lib.Expectation[*HTTPValueSelectorInput]
}

// newMock returns a new Mock with default values set.
func newMock() *HTTPMock {
	return &HTTPMock{
		BaseMock:     lib.NewMock(),
		Callbacks:    make([]Callback, 0),
		PostActions:  make([]*PostActionDef, 0),
		Mappers:      make([]Mapper, 0),
		Pipes:        make([]lib.Piping, 0),
		expectations: make([]*lib.Expectation[*HTTPValueSelectorInput], 0),
	}
}

func (m *HTTPMock) GetExpectations() []*lib.Expectation[*HTTPValueSelectorInput] {
	return m.expectations
}

// Build allows users to use the Mock itself as a HTTPMockBuilder.
func (m *HTTPMock) Build() (*HTTPMock, error) {
	return m, nil
}

func (m *HTTPMock) Prepare() {
	for _, e := range m.expectations {
		ee, ok := e.Matcher.(matcher.OnAfterMockServed)
		if ok {
			m.after = append(m.after, ee)
		}
	}
}

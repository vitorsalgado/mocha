package mhttp

import (
	"net/http"
	"net/url"

	"github.com/vitorsalgado/mocha/v3/foundation"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

var _ foundation.Mock = (*HTTPMock)(nil)

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

type HTTPExpectation = foundation.Expectation[*HTTPValueSelectorInput]

// HTTPMock holds metadata and expectations to be matched against HTTP requests in order to serve mocked responses.
// This is the core entity of this project, and most features work based on it.
type HTTPMock struct {
	*foundation.BaseMock[Reply]

	// Callbacks holds a Callback list to be executed after the Mock was matched and served.
	Callbacks []Callback

	PostActions []*PostActionDef

	// Source describes the source of the mock. E.g.: if it was built from a file,
	// it will contain the filename.
	Source string

	// Mappers stores response mappers associated with this Mock.
	Mappers []Mapper

	after        []matcher.OnAfterMockServed
	expectations []*foundation.Expectation[*HTTPValueSelectorInput]
}

// newMock returns a new Mock with default values set.
func newMock() *HTTPMock {
	return &HTTPMock{
		BaseMock:     foundation.NewMock[Reply](),
		Callbacks:    make([]Callback, 0),
		PostActions:  make([]*PostActionDef, 0),
		Mappers:      make([]Mapper, 0),
		expectations: make([]*foundation.Expectation[*HTTPValueSelectorInput], 0),
	}
}

func (m *HTTPMock) GetExpectations() []*foundation.Expectation[*HTTPValueSelectorInput] {
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

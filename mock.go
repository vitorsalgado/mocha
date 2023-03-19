package mocha

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/vitorsalgado/mocha/v3/matcher"
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
	App *Mocha

	// Mock is the matched Mock for the current HTTP request.
	Mock *Mock
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
	App *Mocha

	// Mock is the matched Mock for the current HTTP request.
	Mock *Mock

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
	App *Mocha

	// Mock is the matched Mock for the current HTTP request.
	Mock *Mock

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
	Handle(fields map[string]any, b *MockBuilder) error
}

// Extension describes a component that can be registered within the mock server and used lately.
// Only one instance of a specific Extension should be registered, but, depending on the component,
// it could be used many times.
type Extension interface {
	UniqueName() string
}

// valueSelector defines a function that will be used to extract the value that will be passed to the associated matcher.
type valueSelector func(r *valueSelectorInput) any

type valueSelectorInput struct {
	// RawRequest is the original incoming http.Request.
	RawRequest *http.Request

	// URL is the full request url.URL, including scheme, host, port.
	URL *url.URL

	Query url.Values

	Form url.Values

	// ParsedBody is the parsed http.Request body.
	ParsedBody any
}

// expectation holds metadata related to one http.Request Matcher.
type expectation struct {
	// Target is an optional metadata that describes the target of the matcher.
	// Example: the target could have the "header", meaning that the matcher will be applied to one request misc.Header
	Target matchTarget

	Key string

	// Matcher associated with this expectation.
	Matcher matcher.Matcher

	// ValueSelector will extract the http.Request or a specific field of it and feed it to the associated Matcher.
	ValueSelector valueSelector

	// Weight of this expectation.
	Weight weight
}

// matchResult holds information related to a matching operation.
type matchResult struct {
	// Details is the list of non-matches messages.
	Details []mismatchDetail

	// Weight for the matcher. It helps determine the closest match.
	Weight int

	// Pass indicates whether it matched or not.
	Pass bool
}

// mismatchDetail gives more context about why a matcher did not match.
type mismatchDetail struct {
	MatchersName string
	Target       matchTarget
	Key          string
	Result       *matcher.Result
	Err          error
}

// mockFileData represents the data that is passed to Mock files during template parsing.
type mockFileData struct {
	App *templateAppWrapper
}

// weight helps to detect the closest mock match.
type weight int8

// Enum of weight.
const (
	_weightNone weight = iota
	_weightLow  weight = iota * 2
	_weightVeryLow
	_weightRegular
	_weightHigh
)

type matchTarget int8

// matchTarget constants to help debug unmatched requests.
const (
	_targetRequest matchTarget = iota
	_targetScheme
	_targetMethod
	_targetURL
	_targetHeader
	_targetQuery
	_targetBody
	_targetForm
)

func (mt matchTarget) String() string {
	switch mt {
	case _targetRequest:
		return "request"
	case _targetScheme:
		return "scheme"
	case _targetMethod:
		return "method"
	case _targetURL:
		return "url"
	case _targetHeader:
		return "header"
	case _targetQuery:
		return "query"
	case _targetBody:
		return "body"
	case _targetForm:
		return "form"
	default:
		return ""
	}
}

// Builder describes a Mock builder.
type Builder interface {
	Build(app *Mocha) (*Mock, error)
}

// Mock holds metadata and expectations to be matched against HTTP requests in order to serve mocked responses.
// This is the core entity of this project, and most features work based on it.
type Mock struct {
	// ID is the unique identifier of a Mock
	ID string

	// Name describes the mock. Useful for debugging.
	Name string

	// Priority sets the priority of a Mock.
	Priority int

	// Reply is the responder that will be used to serve the HTTP response stub, once matched against an
	// HTTP request.
	Reply Reply

	// Enabled indicates if the Mock is enabled or disabled.
	// Only enabled mocks are considered during the request matching phase.
	Enabled bool

	// Callbacks holds a Callback list to be executed after the Mock was matched and served.
	Callbacks []Callback

	PostActions []*PostActionDef

	// Source describes the source of the mock. E.g.: if it was built from a file,
	// it will contain the filename.
	Source string

	// Delay sets the duration to delay serving the mocked response.
	Delay time.Duration

	// Mappers stores response mappers associated with this Mock.
	Mappers []Mapper

	after        []matcher.OnAfterMockServed
	expectations []*expectation
	mu           sync.RWMutex
	hits         int
}

// newMock returns a new Mock with default values set.
func newMock() *Mock {
	return &Mock{
		ID:           uuid.New().String(),
		Enabled:      true,
		Callbacks:    make([]Callback, 0),
		PostActions:  make([]*PostActionDef, 0),
		expectations: make([]*expectation, 0),
	}
}

// Inc increment one Mock call.
func (m *Mock) Inc() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hits++
}

// Dec reduce one Mock call.
func (m *Mock) Dec() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hits--
}

// Hits returns the amount of time this Mock was matched to a request and served.
func (m *Mock) Hits() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.hits
}

// HasBeenCalled checks if the Mock was called at least once.
func (m *Mock) HasBeenCalled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.hits > 0
}

// Enable enables the Mock.
// The Mock will be eligible to be matched.
func (m *Mock) Enable() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Enabled = true
}

// Disable disables the Mock.
// The Mock will not be eligible to be matched.
func (m *Mock) Disable() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Enabled = false
}

// Build allows users to use the Mock itself as a MockBuilder.
func (m *Mock) Build() (*Mock, error) {
	return m, nil
}

// matchExpectations checks if the current Mock matches against a list of expectations.
// Will iterate through all expectations even if it doesn't match early.
func (m *Mock) matchExpectations(ri *valueSelectorInput, expectations []*expectation) *matchResult {
	w := 0
	ok := true
	details := make([]mismatchDetail, 0)

	for _, exp := range expectations {
		var val any
		if exp.ValueSelector != nil {
			val = exp.ValueSelector(ri)
		}

		result, err := m.matchExpectation(exp, val)

		if err != nil {
			ok = false
			details = append(details, mismatchDetail{
				MatchersName: exp.Matcher.Name(),
				Target:       exp.Target,
				Key:          exp.Key,
				Err:          err,
			})

			continue
		}

		if result.Pass {
			w += int(exp.Weight)
		} else {
			ok = false
			details = append(details, mismatchDetail{
				MatchersName: exp.Matcher.Name(),
				Target:       exp.Target,
				Key:          exp.Key,
				Result:       result,
			})
		}
	}

	return &matchResult{Pass: ok, Weight: w, Details: details}
}

func (m *Mock) matchExpectation(e *expectation, value any) (result *matcher.Result, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: matcher=%s. %v", e.Matcher.Name(), r)
			return
		}
	}()

	result, err = e.Matcher.Match(value)
	if err != nil {
		err = fmt.Errorf("%s: error while matching. %w", e.Matcher.Name(), err)
	}

	return
}

func (m *Mock) prepare() {
	for _, e := range m.expectations {
		ee, ok := e.Matcher.(matcher.OnAfterMockServed)
		if ok {
			m.after = append(m.after, ee)
		}
	}
}

package mocha

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/vitorsalgado/mocha/v3/internal/colorize"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

// Mock holds metadata and expectations to be matched against HTTP requests in order to serve mocked responses.
// This is core entity of this project, mostly features works based on it.
type Mock struct {
	// ID is unique identifier for a Mock
	ID string

	// Name is an optional metadata. It helps to find and debug mocks.
	Name string

	// Priority sets the priority for a Mock.
	Priority int

	// Reply is the responder that will be used to serve the HTTP response stub, once matched against an
	// HTTP request.
	Reply Reply

	// Enabled indicates if the Mock is enabled or disabled. Only enabled mocks are matched.
	Enabled bool

	// PostActions holds PostAction list to be executed after the Mock was matched and served.
	PostActions []PostAction

	// Source describes the source of the mock. E.g.: if it wast built from a file,
	// it will contain the filename.
	Source string

	// Delay sets the duration to delay serving the mocked response.
	Delay time.Duration

	// Mappers stores response mappers associated with this Mock.
	Mappers []Mapper

	expectations []*expectation
	mu           sync.Mutex
	hits         int
}

type Builder interface {
	Build() (*Mock, error)
}

// RequestValues groups HTTP request data, including the parsed body, if any.
type RequestValues struct {
	// RawRequest is the original incoming http.Request.
	RawRequest *http.Request

	// URL is full request url.URL, including scheme, host, port.
	URL *url.URL

	// URLPathSegments stores path segments.
	// Eg.: /test/100 -> []string{"test", "100"}
	URLPathSegments []string

	// ParsedBody is the parsed http.Request body parsed by a RequestBodyParser instance.
	// It'll be nil if the HTTP request does not contain a body.
	ParsedBody any

	// App exposes Mocha instance associated with the incoming HTTP request.
	App *Mocha

	// Mock is the matched Mock for the current HTTP request.
	Mock *Mock
}

// PostActionInput represents the arguments that will be passed to every PostAction implementation
type PostActionInput struct {
	// RawRequest is the original incoming http.Request.
	RawRequest *http.Request

	// URL is full request url.URL, including scheme, host, port.
	URL *url.URL

	// ParsedBody is the parsed http.Request body.
	ParsedBody any

	// App exposes Mocha instance associated with the incoming HTTP request.
	App *Mocha

	// Mock is the matched Mock for the current HTTP request.
	Mock *Mock

	// Stub is the HTTP response Stub served.
	Stub *Stub

	// Args allow passing custom arguments to a PostAction.
	Args any
}

// PostAction defines the contract for an action that will be executed after serving a mocked HTTP response.
type PostAction interface {
	// Run runs the PostAction implementation.
	Run(input *PostActionInput) error
}

// Mapper is the function definition to be used to map Mock Stub before serving it.
// Mapper doesn't work with reply.From or Proxy.
type Mapper func(requestValues *RequestValues, res *Stub) error

// Extension describes a component that can registered within the mock server, and used lately.
// Only one instance of Extension should be registered, but, depending on the component, it could be used many times.
type Extension interface {
	UniqueName() string
}

type (
	// valueSelector defines a function that will be used to extract the value that will be passed to the associated matcher.
	valueSelector func(r *valueSelectorInput) any

	valueSelectorInput struct {
		// RawRequest is the original incoming http.Request.
		RawRequest *http.Request

		// URL is full request url.URL, including scheme, host, port.
		URL *url.URL

		// ParsedBody is the parsed http.Request body.
		ParsedBody any
	}

	// expectation holds metadata related to one http.Request Matcher.
	expectation struct {
		// Target is an optional metadata that describes the target of the matcher.
		// Example: the target could have the "header", meaning that the matcher will be applied to one request header.
		Target matchTarget

		Key string

		// Matcher associated with this expectation.
		Matcher matcher.Matcher

		// ValueSelector will extract the http.Request or a portion of it and feed it to the associated Matcher.
		ValueSelector valueSelector

		// Weight of this expectation.
		Weight weight
	}

	// matchResult holds information related to a matching operation.
	matchResult struct {
		// Details is the list of non matches messages.
		Details []mismatchDetail

		// Weight for the matcher. It helps determine the closest match.
		Weight int

		// Pass indicates whether it matched or not.
		Pass bool
	}

	// mismatchDetail gives more ctx about why a matcher did not match.
	mismatchDetail struct {
		MatchersName string
		Target       matchTarget
		Desc         string
		Err          error
	}
)

// weight helps to detect the closest mock match.
type weight int

// Enum of weight.
const (
	_weightNone weight = iota
	_weightLow  weight = iota * 2
	_weightVeryLow
	_weightRegular
	_weightHigh
)

type matchTarget int

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

// newMock returns a new Mock with default values set.
func newMock() *Mock {
	return &Mock{
		ID:           uuid.New().String(),
		Enabled:      true,
		expectations: make([]*expectation, 0),
		PostActions:  make([]PostAction, 0),

		mu: sync.Mutex{},
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
	return m.hits
}

// HasBeenCalled checks if the Mock was called at least once.
func (m *Mock) HasBeenCalled() bool {
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

// Build allow users use Mock as a Builder.
func (m *Mock) Build() (*Mock, error) {
	return m, nil
}

// requestMatches checks if current Mock matches against a list of expectations.
// Will iterate through all expectations even if it doesn't match early.
func (m *Mock) requestMatches(ri *valueSelectorInput, expectations []*expectation) *matchResult {
	w := 0
	ok := true
	details := make([]mismatchDetail, 0)

	for _, exp := range expectations {
		var val any
		if exp.ValueSelector != nil {
			val = exp.ValueSelector(ri)
		}

		result, err := doMatches(exp, val)

		if err != nil {
			ok = false
			details = append(details, mismatchDetail{
				MatchersName: exp.Matcher.Name(),
				Target:       exp.Target,
				Desc: fmt.Sprintf(
					"%s => Error: %s",
					colorize.Bold(exp.Matcher.Name()),
					colorize.Red(err.Error()),
				),
				Err: err,
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
				Desc:         result.Message,
			})
		}
	}

	return &matchResult{Pass: ok, Weight: w, Details: details}
}

func doMatches(e *expectation, value any) (result *matcher.Result, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("matcher %s panicked. reason=%v", e.Matcher.Name(), r)
			return
		}
	}()

	result, err = e.Matcher.Match(value)

	return
}

package mock

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vitorsalgado/mocha/internal/params"
	"github.com/vitorsalgado/mocha/matcher"
)

type (
	// Mock holds metadata and expectations to be matched against HTTP requests in order to serve mocked responses.
	// This is core entity of this project, mostly features works based on it.
	Mock struct {
		// ID is unique identifier for a Mock
		ID int

		// Name is an optional metadata. It helps to find and debug mocks.
		Name string

		// Priority sets the priority for a Mock.
		Priority int

		// Expectations are a list of Expectation. These will run on every request to find the correct Mock.
		Expectations []any

		// AfterExpectations are a list of Expectation. They will be executed after the request was matched to a Mock.
		// This allows stateful matchers whose state data should not be evaluated every match attempt.
		AfterExpectations []any

		// Reply is the responder that will be used to serve the HTTP response stub, once matched against an
		// HTTP request.
		Reply Reply

		// Hits holds the amount of time this Mock was called and served.
		Hits int

		// Enabled indicates if the Mock is enabled or disabled. Only enabled mocks are matched.
		Enabled bool

		// PostActions holds PostAction list to be executed after the Mock was matched and served.
		PostActions []PostAction

		mu *sync.Mutex
	}

	// PostActionArgs represents the arguments that will be passed to every PostAction implementation
	PostActionArgs struct {
		Request  *http.Request
		Response *Response
		Mock     *Mock
		Params   params.Params
	}

	// PostAction defines the contract for an action that will be executed after serving a mocked HTTP response.
	PostAction interface {
		// Run runs the PostAction implementation.
		Run(args PostActionArgs) error
	}

	// Response defines the HTTP response that will be served once a Mock is matched for an HTTP Request.
	Response struct {
		Status  int
		Header  http.Header
		Cookies []*http.Cookie
		Body    io.Reader
		Delay   time.Duration
		Err     error
		Mappers []ResponseMapper
	}

	// Reply defines the contract to configure an HTTP responder.
	Reply interface {
		// Build returns a Response stub to be served.
		Build(*http.Request, *Mock, params.Params) (*Response, error)
	}

	// ResponseMapperArgs represents the expected arguments for every ResponseMapper.
	ResponseMapperArgs struct {
		Request    *http.Request
		Parameters params.Params
	}

	// ResponseMapper is the function definition to be used to map Mock Response before serving it.
	ResponseMapper func(res *Response, args ResponseMapperArgs) error

	// Storage is the definition for Mock repository.
	Storage interface {
		// Save saves the Mock.
		Save(mock *Mock)

		// FetchEligible returns mocks that can be matched against requests.
		FetchEligible() []*Mock

		// FetchAll returns all stored Mock instances.
		FetchAll() []*Mock

		// Delete removes a Mock by its ID.
		Delete(id int)

		// Flush removes all stored mocks.
		Flush()
	}

	// ExpectationValuePicker is the function used to extract a specific value from http.Request.
	ExpectationValuePicker[V any] func(r *matcher.RequestInfo) V

	// Expectation holds metadata related to one http.Request matcher.Matcher.
	Expectation[V any] struct {
		// Name is an optional metadata to help debugging request expectations.
		Name string

		// Matcher associated with this Expectation.
		Matcher matcher.Matcher[V]

		// ValuePicker will extract the http.Request or a portion of it and feed it to the associated Matcher.
		ValuePicker ExpectationValuePicker[V]

		// Weight of this Expectation.
		Weight int
	}

	// MatchResult holds information related to a matching operation.
	MatchResult struct {
		// NonMatched is the list of Mock they were matched.
		NonMatched []string

		// Weight for the matcher.Matcher
		Weight int

		// IsMatch indicates whether it matched or not.
		IsMatch bool
	}
)

type autoID struct {
	id int32
}

func (i *autoID) next() int {
	atomic.AddInt32(&i.id, 1)
	return int(i.id)
}

var id = autoID{}

// New returns a new Mock with default values set.
func New() *Mock {
	return &Mock{
		ID:                id.next(),
		Enabled:           true,
		Expectations:      make([]any, 0),
		AfterExpectations: make([]any, 0),
		PostActions:       make([]PostAction, 0),

		mu: &sync.Mutex{},
	}
}

// Hit notify that the Mock was called.
func (m *Mock) Hit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Hits++
}

// Dec reduce one Mock call.
func (m *Mock) Dec() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Hits--
}

// Called checks if the Mock was called at least once.
func (m *Mock) Called() bool {
	return m.Hits > 0
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

// Matches checks if current Mock matches against a list of expectations.
func (m *Mock) Matches(params matcher.Args, expectations []any) (MatchResult, error) {
	weight := 0
	for _, expect := range expectations {
		var matched bool
		var err error
		var w int

		switch e := expect.(type) {
		default:
			return MatchResult{IsMatch: false, Weight: weight},
				fmt.Errorf("unhandled matcher type %s", reflect.TypeOf(e))

		case Expectation[any]:
			matched, w, err = matches(e, params)
		case Expectation[string]:
			matched, w, err = matches(e, params)
		case Expectation[float64]:
			matched, w, err = matches(e, params)
		case Expectation[bool]:
			matched, w, err = matches(e, params)
		case Expectation[map[string]any]:
			matched, w, err = matches(e, params)
		case Expectation[map[string]string]:
			matched, w, err = matches(e, params)
		case Expectation[map[string][]string]:
			matched, w, err = matches(e, params)
		case Expectation[[]any]:
			matched, w, err = matches(e, params)
		case Expectation[url.URL]:
			matched, w, err = matches(e, params)
		case Expectation[*http.Request]:
			matched, w, err = matches(e, params)
		case Expectation[url.Values]:
			matched, w, err = matches(e, params)
		}

		if err != nil || !matched {
			return MatchResult{IsMatch: false, Weight: weight}, err
		}

		weight += w
	}

	return MatchResult{IsMatch: true}, nil
}

func matches[V any](e Expectation[V], params matcher.Args) (bool, int, error) {
	res, err := e.Matcher(e.ValuePicker(params.RequestInfo), params)
	return res, e.Weight, err
}

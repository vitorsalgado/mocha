package core

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/vitorsalgado/mocha/expect"
	"github.com/vitorsalgado/mocha/internal/colorize"
	"github.com/vitorsalgado/mocha/internal/parameters"
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
		Expectations []Expectation

		// PostExpectations are a list of Expectation. They will be executed after the request was matched to a Mock.
		// This allows stateful matchers whose state data should not be evaluated every match attempt.
		PostExpectations []Expectation

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
		Params   parameters.Params
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
		Mappers []ResponseMapper
	}

	// Reply defines the contract to configure an HTTP responder.
	Reply interface {
		// Build returns a Response stub to be served.
		Build(*http.Request, *Mock, parameters.Params) (*Response, error)
	}

	// ResponseMapperArgs represents the expected arguments for every ResponseMapper.
	ResponseMapperArgs struct {
		Request    *http.Request
		Parameters parameters.Params
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

	// Weight helps to detect the closest mock match.
	Weight int

	// Expectation holds metadata related to one http.Request Matcher.
	Expectation struct {
		// Name is an optional metadata to help debugging request expectations.
		Name string

		// Matcher associated with this Expectation.
		Matcher expect.Matcher

		// ValueSelector will extract the http.Request or a portion of it and feed it to the associated Matcher.
		ValueSelector expect.ValueSelector

		// Weight of this Expectation.
		Weight Weight
	}

	// MatchResult holds information related to a matching operation.
	MatchResult struct {
		// NonMatched is the list of Mock they were matched.
		NonMatched []string

		// Weight for the Matcher. It helps determine the closest match.
		Weight int

		// IsMatch indicates whether it matched or not.
		IsMatch bool
	}

	// T is based on testing.T and allow mocha components to log information and errors.
	T interface {
		Cleanup(func())
		Helper()
		Logf(string, ...any)
		Errorf(string, ...any)
		Fatalf(string, ...any)
	}
)

// Enums of Weight.
const (
	WeightNone Weight = iota
	WeightVeryLow
	WeightLow
	WeightRegular
	WeightHigh
)

// NewMock returns a new Mock with default values set.
func NewMock() *Mock {
	return &Mock{
		ID:               id.next(),
		Enabled:          true,
		Expectations:     make([]Expectation, 0),
		PostExpectations: make([]Expectation, 0),
		PostActions:      make([]PostAction, 0),

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
// Will iterate through all expectations even if it doesn't match early.
func (m *Mock) Matches(params expect.Args, expectations []Expectation, t T) (MatchResult, error) {
	t.Helper()

	weight := 0
	finalMatched := true

	for _, exp := range expectations {
		var matched bool
		var err error

		matched, err = matches(exp, params, t)

		// fail fast if an error occurs
		if err != nil {
			msg := fmt.Sprintf("matcher %s returned an error: %v", exp.Name, err)
			t.Logf(msg)

			return MatchResult{IsMatch: false, Weight: weight}, fmt.Errorf(msg)
		}

		if !matched {
			finalMatched = matched
		}

		weight += int(exp.Weight)
	}

	return MatchResult{IsMatch: finalMatched, Weight: weight}, nil
}

func matches(e Expectation, params expect.Args, t T) (bool, error) {
	t.Helper()

	val := e.ValueSelector(params.RequestInfo)
	res, err := e.Matcher.Matches(val, params)

	if err != nil {
		t.Logf("\n%s\nError: %s\nContext: %s",
			colorize.Yellow(fmt.Sprintf("Matcher %s returned an error", colorize.Bold(e.Matcher.Name))),
			colorize.Yellow(err.Error()),
			e.Name)

		return false, err
	}

	if !res {
		desc := ""
		if e.Matcher.DescribeMismatch != nil {
			desc = e.Matcher.DescribeMismatch(e.Name, val)
		}

		t.Logf("\n%s\n%s\n%s",
			colorize.Yellow(fmt.Sprintf("Matcher %s dit not match.", colorize.Bold(e.Matcher.Name))),
			"applied to: "+e.Name,
			desc)
	}

	return res, err
}

package mocha

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/vitorsalgado/mocha/v3/expect"
	"github.com/vitorsalgado/mocha/v3/internal/autoid"
	"github.com/vitorsalgado/mocha/v3/params"
	"github.com/vitorsalgado/mocha/v3/reply"
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

		// Reply is the responder that will be used to serve the HTTP response stub, once matched against an
		// HTTP request.
		Reply reply.Reply

		// Enabled indicates if the Mock is enabled or disabled. Only enabled mocks are matched.
		Enabled bool

		// PostActions holds PostAction list to be executed after the Mock was matched and served.
		PostActions []PostAction

		mu   *sync.Mutex
		hits int
	}

	// PostActionArgs represents the arguments that will be passed to every PostAction implementation
	PostActionArgs struct {
		Request  *http.Request
		Response *reply.Response
		Mock     *Mock
		Params   params.P
	}

	// PostAction defines the contract for an action that will be executed after serving a mocked HTTP response.
	PostAction interface {
		// Run runs the PostAction implementation.
		Run(args PostActionArgs) error
	}

	// ValueSelector defines a function that will be used to extract RequestInfo value and provide it to Matcher instances.
	ValueSelector func(r *expect.RequestInfo) any

	// Expectation holds metadata related to one http.Request Matcher.
	Expectation struct {
		// Target is an optional metadata that describes the target of the matcher.
		// Example: the target could have the "header", meaning that the matcher will be applied to one request header.
		Target string

		// Matcher associated with this Expectation.
		Matcher expect.Matcher

		// ValueSelector will extract the http.Request or a portion of it and feed it to the associated Matcher.
		ValueSelector ValueSelector

		// Weight of this Expectation.
		Weight weight
	}
)

type (
	// weight helps to detect the closest mock match.
	weight int

	// matchResult holds information related to a matching operation.
	matchResult struct {
		// MismatchDetails is the list of non matches messages.
		MismatchDetails []mismatchDetail

		// Weight for the Matcher. It helps determine the closest match.
		Weight int

		// IsMatch indicates whether it matched or not.
		IsMatch bool
	}

	// mismatchDetail gives more context about why a matcher did not match.
	mismatchDetail struct {
		Name        string
		Target      string
		Description string
	}
)

// Enums of weight.
const (
	_weightNone weight = iota
	_weightVeryLow
	_weightLow
	_weightRegular
	_weightHigh
)

// newMock returns a new Mock with default values set.
func newMock() *Mock {
	return &Mock{
		ID:           autoid.Next(),
		Enabled:      true,
		Expectations: make([]Expectation, 0),
		PostActions:  make([]PostAction, 0),

		mu: &sync.Mutex{},
	}
}

// Hit notify that the Mock was called.
func (m *Mock) Hit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hits++
}

// Hits returns the amount of time this Mock was matched to a request and served.
func (m *Mock) Hits() int {
	return m.hits
}

// Dec reduce one Mock call.
func (m *Mock) Dec() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hits--
}

// Called checks if the Mock was called at least once.
func (m *Mock) Called() bool {
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

// matches checks if current Mock matches against a list of expectations.
// Will iterate through all expectations even if it doesn't match early.
func (m *Mock) matches(ri *expect.RequestInfo, expectations []Expectation) (matchResult, error) {
	w := 0
	hasMatched := true
	details := make([]mismatchDetail, 0)

	for _, exp := range expectations {
		matched, detail, err := matches(exp, ri)

		// fail fast if an error occurs
		if err != nil {
			return matchResult{IsMatch: false, Weight: w},
				fmt.Errorf("matcher %s returned an error=%v", exp.Target, err)
		}

		if !matched {
			details = append(details, detail)
			hasMatched = matched
		}

		w += int(exp.Weight)
	}

	return matchResult{IsMatch: hasMatched, Weight: w, MismatchDetails: details}, nil
}

func matches(e Expectation, params *expect.RequestInfo) (bool, mismatchDetail, error) {
	val := e.ValueSelector(params)
	res, err := e.Matcher.Match(val)

	if err != nil {
		return false,
			mismatchDetail{Name: e.Matcher.Name(), Target: e.Target},
			err
	}

	if !res {
		return res, mismatchDetail{
			Name:        e.Matcher.Name(),
			Target:      e.Target,
			Description: e.Matcher.DescribeFailure(val),
		}, err
	}

	return res, mismatchDetail{}, err
}

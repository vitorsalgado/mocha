package mocha

import (
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
		Run(args *PostActionArgs) error
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

		// OK indicates whether it matched or not.
		OK bool
	}

	// mismatchDetail gives more context about why a matcher did not match.
	mismatchDetail struct {
		Name   string
		Target string
		Desc   string
		Err    error
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

// requestMatches checks if current Mock matches against a list of expectations.
// Will iterate through all expectations even if it doesn't match early.
func (m *Mock) requestMatches(ri *expect.RequestInfo, expectations []Expectation) *matchResult {
	w := 0
	ok := true
	details := make([]mismatchDetail, 0)

	for _, exp := range expectations {
		val := exp.ValueSelector(ri)
		result, err := matches(exp, val)

		if err != nil {
			ok = false
			details = append(details, mismatchDetail{
				Name:   exp.Matcher.Name(),
				Target: exp.Target,
				Desc:   err.Error(),
				Err:    err,
			})

			continue
		}

		if result.OK {
			w += int(exp.Weight)
		} else {
			ok = false
			details = append(details, mismatchDetail{
				Name:   exp.Matcher.Name(),
				Target: exp.Target,
				Desc:   result.DescribeFailure(),
			})
		}
	}

	return &matchResult{OK: ok, Weight: w, MismatchDetails: details}
}

func matches(e Expectation, value any) (expect.Result, error) {
	res, err := e.Matcher.Match(value)

	if err != nil {
		return expect.Result{}, err
	}

	return res, nil
}

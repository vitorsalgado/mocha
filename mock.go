package mocha

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/vitorsalgado/mocha/v3/internal/autoid"
	"github.com/vitorsalgado/mocha/v3/internal/colorize"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/reply"
)

// Mock holds metadata and expectations to be matched against HTTP requests in order to serve mocked responses.
// This is core entity of this project, mostly features works based on it.
type Mock struct {
	// ID is unique identifier for a Mock
	ID int

	// Name is an optional metadata. It helps to find and debug mocks.
	Name string

	// Priority sets the priority for a Mock.
	Priority int

	// Reply is the responder that will be used to serve the HTTP response stub, once matched against an
	// HTTP request.
	Reply reply.Reply

	// Enabled indicates if the Mock is enabled or disabled. Only enabled mocks are matched.
	Enabled bool

	// PostActions holds PostAction list to be executed after the Mock was matched and served.
	PostActions []PostAction

	// Source describes the source of the mock. E.g.: if it wast built from a file,
	// it will contain the filename.
	Source string

	// Delay sets the duration to delay sending the mocked response.
	Delay time.Duration

	expectations []*expectation
	mu           *sync.Mutex
	hits         int
}

type Builder interface {
	Build() (*Mock, error)
}

// PostActionArgs represents the arguments that will be passed to every PostAction implementation
type PostActionArgs struct {
	Request  *http.Request
	Response *reply.Response
	Mock     *Mock
	Params   reply.Params
}

// PostAction defines the contract for an action that will be executed after serving a mocked HTTP response.
type PostAction interface {
	// Run runs the PostAction implementation.
	Run(args *PostActionArgs) error
}

// target constants to help debug unmatched requests.
const (
	_targetRequest = "request(no specific field)"
	_targetMethod  = "method"
	_targetURL     = "url"
	_targetHeader  = "header"
	_targetQuery   = "query"
	_targetBody    = "body"
	_targetForm    = "form"
)

type (
	// valueSelector defines a function that will be used to extract RequestInfo value and provide it to matcher instances.
	valueSelector func(r *matcher.RequestInfo) any

	// expectation holds metadata related to one http.Request Matcher.
	expectation struct {
		// Target is an optional metadata that describes the target of the matcher.
		// Example: the target could have the "header", meaning that the matcher will be applied to one request header.
		Target string

		// Matcher associated with this expectation.
		Matcher matcher.Matcher

		// ValueSelector will extract the http.Request or a portion of it and feed it to the associated Matcher.
		ValueSelector valueSelector

		// Weight of this expectation.
		Weight weight
	}

	// matchResult holds information related to a matching operation.
	matchResult struct {
		// MismatchDetails is the list of non matches messages.
		MismatchDetails []mismatchDetail

		// Weight for the matcher. It helps determine the closest match.
		Weight int

		// OK indicates whether it matched or not.
		OK bool
	}

	// mismatchDetail gives more ctx about why a matcher did not match.
	mismatchDetail struct {
		Name   string
		Target string
		Desc   string
		Err    error
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

// newMock returns a new Mock with default values set.
func newMock() *Mock {
	return &Mock{
		ID:           autoid.Next(),
		Enabled:      true,
		expectations: make([]*expectation, 0),
		PostActions:  make([]PostAction, 0),

		mu: &sync.Mutex{},
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
func (m *Mock) requestMatches(ri *matcher.RequestInfo, expectations []*expectation) *matchResult {
	w := 0
	ok := true
	details := make([]mismatchDetail, 0)

	for _, exp := range expectations {
		val := exp.ValueSelector(ri)
		result, err := doMatches(exp, val)

		if err != nil {
			ok = false
			details = append(details, mismatchDetail{
				Name:   exp.Matcher.Name(),
				Target: exp.Target,
				Desc: fmt.Sprintf(
					"%s => Error: %s",
					colorize.Bold(exp.Matcher.Name()),
					colorize.Red(err.Error()),
				),
				Err: err,
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

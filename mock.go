package mocha

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/vitorsalgado/mocha/v3/internal/colorize"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/reply"
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
	Reply reply.Reply

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

// PostActionIn represents the arguments that will be passed to every PostAction implementation
type PostActionIn struct {
	Request  *http.Request
	Response *reply.ResponseStub
	Params   reply.Params
}

// PostAction defines the contract for an action that will be executed after serving a mocked HTTP response.
type PostAction interface {
	// Run runs the PostAction implementation.
	Run(args *PostActionIn) error
}

// Mapper is the function definition to be used to map Mock ResponseStub before serving it.
// Mapper doesn't work with reply.Forward or Proxy.
type Mapper func(res *reply.ResponseStub, args *MapperIn) error

// MapperIn represents the expected arguments for every Mapper.
type MapperIn struct {
	Request    *http.Request
	Parameters reply.Params
}

type (
	// values groups HTTP request information to be passed to each Matcher.
	values struct {
		// Request is the actual http.Request.
		Request *http.Request

		// ParsedBody is http.Request parsed body.
		// Value of parsed body can vary depending on the mocha.RequestBodyParser that parsed the request.
		ParsedBody any
	}

	// valueSelector defines a function that will be used to extract the value that will be passed to the associated matcher.
	valueSelector func(r *values) any

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

// MarshalJSON marshal Mock to a JSON that can be loaded later by this mock server.
func (m *Mock) MarshalJSON() ([]byte, error) {
	ext := make(map[string]any)

	fields := make(map[string]any)
	headers := make(map[string]any)
	queries := make(map[string]any)
	body := make([]any, 0)

	for _, e := range m.expectations {
		switch e.Target {
		case _targetMethod:
			ext["method"] = e.Matcher.Spec()
		case _targetURL:
			ext[e.Key] = e.Matcher.Spec()
		case _targetQuery:
			queries[e.Key] = e.Matcher.Spec()
		case _targetHeader:
			headers[e.Key] = e.Matcher.Spec()
		case _targetForm:
			fields[e.Key] = e.Matcher.Spec()
		case _targetBody:
			body = append(body, e.Matcher.Spec())
		case _targetRequest:
			mm := e.Matcher.Spec()
			if mm == nil {
				continue
			}

			arr, ok := mm.([]any)
			if !ok {
				return nil, fmt.Errorf("must be an array")
			}

			f := arr[0].(string)
			v := arr[1]

			ext[f] = v
		}
	}

	ext["body"] = body
	ext["header"] = headers
	ext["query"] = queries
	ext["fields"] = fields

	res := m.Reply.Spec()
	f := res[0].(string)
	v := res[1]

	ext[f] = v

	return json.MarshalIndent(ext, "", " ")
}

// requestMatches checks if current Mock matches against a list of expectations.
// Will iterate through all expectations even if it doesn't match early.
func (m *Mock) requestMatches(ri *values, expectations []*expectation) *matchResult {
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
				Desc:         result.Message(),
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

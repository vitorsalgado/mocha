package mock

import (
	"io"
	"net/http"
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

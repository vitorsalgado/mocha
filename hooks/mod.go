package hooks

import (
	"fmt"
	"net/http"
	"time"
)

// OnRequest event is triggered every time a request arrives at the mock handler.
type OnRequest struct {
	Request   Request
	StartedAt time.Time
}

// OnRequestMatch event is triggered when a mock is found for a request.
type OnRequestMatch struct {
	Request            Request
	ResponseDefinition Response
	Mock               Mock
	Elapsed            time.Duration
	Body               any
}

// OnRequestNotMatched event is triggered when no mocks are found for a request.
type OnRequestNotMatched struct {
	Request Request
	Result  Result
}

// OnError event is triggered when an error occurs during request matching.
type OnError struct {
	Request Request
	Err     error
}

// hook data transfer objects
type (
	// Request defines information from http.Request to be logged.
	Request struct {
		Method     string
		Path       string
		RequestURI string
		Host       string
		Header     http.Header
		Body       any
	}

	// Response defines HTTP Response information to be logged.
	Response struct {
		Status int
		Header http.Header
	}

	// Mock defines core.Mock information to be logged.
	Mock struct {
		ID   int
		Name string
	}

	// ResultDetail defines matching result details to be logged.
	ResultDetail struct {
		Name        string
		Target      string
		Description string
	}

	// Result defines matching result to be logged.
	Result struct {
		HasClosestMatch bool
		ClosestMatch    Mock
		Details         []ResultDetail
	}
)

func (r *Request) FullURL() string {
	return fmt.Sprintf("%s%s", r.Host, r.RequestURI)
}

package reply

import (
	"io"
	"net/http"
	"time"
)

// M implements mock data that should be available on reply build functions.
type M interface {
	// Hits return mock total hits.
	Hits() int
}

// Reply defines the contract to configure an HTTP responder.
type Reply interface {
	// Build returns a Response stub to be served.
	Build(*http.Request, M, Params) (*Response, error)
}

// Response defines the HTTP response that will be served once a Mock is matched for an HTTP Request.
type Response struct {
	Status  int
	Header  http.Header
	Cookies []*http.Cookie
	Body    io.Reader
	Mappers []Mapper

	Delay time.Duration
}

// Mapper is the function definition to be used to map Mock Response before serving it.
type Mapper func(res *Response, args *MapperArgs) error

// MapperArgs represents the expected arguments for every Mapper.
type MapperArgs struct {
	Request    *http.Request
	Parameters Params
}

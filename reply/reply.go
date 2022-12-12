package reply

import (
	"net/http"
)

// MockInfo implements mock data that should be available on reply build functions.
type MockInfo struct {
	// Hits return mock total hits.
	Hits int
}

// Reply defines the contract to configure an HTTP responder.
type Reply interface {
	// Prepare runs once during mock building.
	// Useful for pre-configurations or validations that needs to be executed once.
	Prepare() error

	Spec() []any

	// Build returns a Response stub to be served.
	// Return Response nil if the HTTP response was rendered inside the Build function.
	Build(w http.ResponseWriter, r *http.Request) (*Response, error)
}

// Arg groups extra parameters to build a Reply.
type Arg struct {
	MockInfo MockInfo
	Params   Params
}

// Response defines the HTTP response that will be served once a Mock is matched for an HTTP Request.
type Response struct {
	Status  int
	Header  http.Header
	Cookies []*http.Cookie
	Body    []byte
}

// Sent checks if response was already sent by the Reply implementation.
func (r *Response) Sent() bool {
	return r != nil
}

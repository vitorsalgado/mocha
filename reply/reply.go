package reply

import (
	"net/http"

	"github.com/vitorsalgado/mocha/v3/types"
)

// Reply defines the contract to configure an HTTP responder.
type Reply interface {

	// Build returns an HTTP response Stub to be served.
	// Return Stub nil if the HTTP response was rendered inside the Build function.
	Build(w http.ResponseWriter, r *types.RequestValues) (*Stub, error)
}

// Pre describes a Reply that has preparations steps to run on mocking building.
type Pre interface {
	// Pre runs once during mock building.
	// Useful for pre-configurations or validations that needs to be executed once.
	Pre() error
}

// Stub defines the HTTP response that will be served once a Mock is matched for an HTTP Request.
type Stub struct {
	StatusCode int
	Header     http.Header
	Cookies    []*http.Cookie
	Body       []byte
}

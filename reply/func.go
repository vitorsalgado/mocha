package reply

import (
	"net/http"
)

// FunctionReply represents a reply that will be built using the given function.
type FunctionReply struct {
	fn func(*http.Request, M, Params) (*Response, error)
}

// Function returns a FunctionReply that builds a response stub using the given function.
func Function(fn func(*http.Request, M, Params) (*Response, error)) *FunctionReply {
	return &FunctionReply{fn: fn}
}

// Build builds a response function using previously provided function.
func (f *FunctionReply) Build(r *http.Request, m M, p Params) (*Response, error) {
	return f.fn(r, m, p)
}

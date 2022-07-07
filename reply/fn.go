package reply

import (
	"net/http"

	"github.com/vitorsalgado/mocha/internal/params"
	"github.com/vitorsalgado/mocha/mock"
)

// FunctionReply represents a reply that will be built using the given function.
type FunctionReply struct {
	fn func(*http.Request, *mock.Mock, params.Params) (*mock.Response, error)
}

// Function returns a FunctionReply that builds a response stub using the given function.
func Function(fn func(*http.Request, *mock.Mock, params.Params) (*mock.Response, error)) *FunctionReply {
	return &FunctionReply{fn: fn}
}

// Build builds a response function using previously provided function.
func (f *FunctionReply) Build(r *http.Request, m *mock.Mock, p params.Params) (*mock.Response, error) {
	return f.fn(r, m, p)
}

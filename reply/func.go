package reply

import (
	"net/http"
)

var _ Reply = (*FunctionReply)(nil)

// FunctionReply represents a reply that will be built using the given function.
type FunctionReply struct {
	fn func(http.ResponseWriter, *http.Request) (*ResponseStub, error)
}

// Function returns a FunctionReply that builds a response stub using the given function.
func Function(fn func(http.ResponseWriter, *http.Request) (*ResponseStub, error)) *FunctionReply {
	return &FunctionReply{fn: fn}
}

func (f *FunctionReply) Prepare() error { return nil }

func (f *FunctionReply) Spec() []any {
	return []any{}
}

// Build builds a response function using previously provided function.
func (f *FunctionReply) Build(w http.ResponseWriter, r *http.Request) (*ResponseStub, error) {
	return f.fn(w, r)
}

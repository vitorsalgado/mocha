package httpd

import (
	"net/http"
)

// -- Function Reply

var _ Reply = (*FunctionReply)(nil)

// FunctionReply represents a reply that will be built using the given function.
type FunctionReply struct {
	fn func(http.ResponseWriter, *RequestValues) (*Stub, error)
}

// Function returns a FunctionReply that builds a response stub using the given function.
func Function(fn func(http.ResponseWriter, *RequestValues) (*Stub, error)) *FunctionReply {
	return &FunctionReply{fn: fn}
}

// Build builds a response function using the previously provided function.
func (f *FunctionReply) Build(w http.ResponseWriter, r *RequestValues) (*Stub, error) {
	return f.fn(w, r)
}

// -- HTTP Handler Reply

var _ Reply = (*HandlerReply)(nil)

type HandlerReply struct {
	h http.HandlerFunc
}

func Handler(h http.HandlerFunc) *HandlerReply {
	return &HandlerReply{h: h}
}

func (h *HandlerReply) Build(w http.ResponseWriter, r *RequestValues) (*Stub, error) {
	h.h(w, r.RawRequest)
	return nil, nil
}

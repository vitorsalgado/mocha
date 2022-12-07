package reply

import (
	"net/http"
)

var _ Reply = (*HandlerReply)(nil)

type HandlerReply struct {
	h http.HandlerFunc
}

func Handler(h http.HandlerFunc) *HandlerReply {
	return &HandlerReply{h: h}
}

func (h *HandlerReply) Prepare() error { return nil }

func (h *HandlerReply) Build(w http.ResponseWriter, r *http.Request) (*Response, error) {
	h.h(w, r)
	return nil, nil
}

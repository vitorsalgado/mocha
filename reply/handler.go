package reply

import (
	"net/http"

	"github.com/vitorsalgado/mocha/v3/types"
)

var _ Reply = (*HandlerReply)(nil)

type HandlerReply struct {
	h http.HandlerFunc
}

func Handler(h http.HandlerFunc) *HandlerReply {
	return &HandlerReply{h: h}
}

func (h *HandlerReply) Build(w http.ResponseWriter, r *types.RequestValues) (*Stub, error) {
	h.h(w, r.RawRequest)
	return nil, nil
}

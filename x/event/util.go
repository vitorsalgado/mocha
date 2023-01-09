package event

import (
	"net/http"
)

func FromRequest(r *http.Request) *EvtReq {
	return &EvtReq{
		Method: r.Method,
		Path:   r.URL.Path,
		Host:   r.Host,
		Header: r.Header.Clone(),
	}
}

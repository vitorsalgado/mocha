package event

import (
	"net/http"
)

func FromRequest(r *http.Request) *EvtReq {
	return &EvtReq{
		Method:     r.Method,
		Path:       r.URL.Path,
		RequestURI: r.RequestURI,
		Host:       r.Host,
		Header:     r.Header.Clone(),
	}
}

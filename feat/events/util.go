package events

import "net/http"

func FromRequest(r *http.Request) Request {
	return Request{
		Method:     r.Method,
		Path:       r.URL.Path,
		RequestURI: r.RequestURI,
		Host:       r.Host,
		Header:     r.Header.Clone(),
	}
}

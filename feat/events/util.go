package events

import (
	"fmt"
	"net/http"
)

// FromRequest is a helper function that creates a new Request from a http.Request.
func FromRequest(r *http.Request) Request {
	return Request{
		Method:     r.Method,
		Path:       r.URL.Path,
		RequestURI: r.RequestURI,
		Host:       r.Host,
		Header:     r.Header.Clone(),
	}
}

func fullURL(host, uri string) string {
	return fmt.Sprintf("%s%s", host, uri)
}

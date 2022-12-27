// Package testutil contains internal test utilities.
package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type RequestIntent struct {
	Request *http.Request
}

func NewRequest(method, url string, body io.Reader) *RequestIntent {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Panic(err)
	}

	return &RequestIntent{Request: req}
}

func Get(url string) *RequestIntent {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Panic(err)
	}

	return &RequestIntent{Request: req}
}

func Post(url string, body io.Reader) *RequestIntent {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		log.Panic(err)
	}

	return &RequestIntent{Request: req}
}

func PostJSON(url string, body any) *RequestIntent {
	b, err := json.Marshal(body)
	if err != nil {
		log.Panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		log.Panic(err)
	}

	req.Header.Add("content-type", "application/json")

	return &RequestIntent{Request: req}
}

func (req *RequestIntent) Header(key, value string) *RequestIntent {
	req.Request.Header.Add(key, value)
	return req
}

func (req *RequestIntent) Do() (*http.Response, error) {
	return http.DefaultClient.Do(req.Request)
}

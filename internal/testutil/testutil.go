// Package testutil contains internal test utilities.
package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type RequestValues struct {
	Request *http.Request
}

func NewRequest(method, url string, body io.Reader) *RequestValues {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatal(err)
	}

	return &RequestValues{Request: req}
}

func Get(url string) *RequestValues {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &RequestValues{Request: req}
}

func Post(url string, body io.Reader) *RequestValues {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		log.Fatal(err)
	}

	return &RequestValues{Request: req}
}

func PostJSON(url string, body any) *RequestValues {
	b, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("content-type", "application/json")

	return &RequestValues{Request: req}
}

func (req *RequestValues) Header(key, value string) *RequestValues {
	req.Request.Header.Add(key, value)
	return req
}

func (req *RequestValues) Do() (*http.Response, error) {
	return http.DefaultClient.Do(req.Request)
}

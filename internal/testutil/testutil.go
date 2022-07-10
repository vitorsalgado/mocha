// Package testutil contains internal test utilities.
package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Req struct {
	Request *http.Request
}

func Request(method, url string) *Req {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &Req{Request: req}
}

func Get(url string) *Req {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &Req{Request: req}
}

func Post(url string, body io.Reader) *Req {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		log.Fatal(err)
	}

	return &Req{Request: req}
}

func PostJSON(url string, body any) *Req {
	b, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("content-type", "application/json")

	return &Req{Request: req}
}

func (req *Req) Header(key, value string) *Req {
	req.Request.Header.Add(key, value)
	return req
}

func (req *Req) Do() (*http.Response, error) {
	return http.DefaultClient.Do(req.Request)
}

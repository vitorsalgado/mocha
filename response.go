package mocha

import "net/http"

type (
	Response struct {
		Status  int
		Headers map[string]string
		Body    []byte
		Delay   int
	}

	ResponseDelegate func(r *http.Request, mock *Mock) (Response, error)
)

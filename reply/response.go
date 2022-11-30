package reply

import (
	"io"
	"net/http"
	"time"
)

type (
	// Response defines the HTTP response that will be served once a Mock is matched for an HTTP Request.
	Response struct {
		Status  int
		Header  http.Header
		Cookies []*http.Cookie
		Body    io.Reader
		Delay   time.Duration
		Mappers []Mapper
	}

	// MapperArgs represents the expected arguments for every Mapper.
	MapperArgs struct {
		Request    *http.Request
		Parameters Params
	}

	// Mapper is the function definition to be used to map Mock Response before serving it.
	Mapper func(res *Response, args MapperArgs) error
)

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
		Mappers []ResponseMapper
	}

	// ResponseMapperArgs represents the expected arguments for every ResponseMapper.
	ResponseMapperArgs struct {
		Request    *http.Request
		Parameters Params
	}

	// ResponseMapper is the function definition to be used to map Mock Response before serving it.
	ResponseMapper func(res *Response, args ResponseMapperArgs) error
)

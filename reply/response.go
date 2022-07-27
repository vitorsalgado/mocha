package reply

import (
	"io"
	"net/http"
	"time"

	"github.com/vitorsalgado/mocha/params"
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
		Parameters params.P
	}

	// ResponseMapper is the function definition to be used to map Mock Response before serving it.
	ResponseMapper func(res *Response, args ResponseMapperArgs) error
)

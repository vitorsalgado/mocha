package mocha

import (
	"io"
	"net/http"
	"strings"
)

type (
	Response struct {
		Status  int
		Headers map[string]string
		Body    io.Reader
		Delay   int
	}

	Responder func(r *http.Request, mock *Mock) (*Response, error)

	Reply interface {
		Build() Responder
	}

	StdReply struct {
		response *Response
	}
)

func NewReply() *StdReply {
	return &StdReply{response: &Response{}}
}

func OK() *StdReply      { return NewReply().Status(http.StatusOK) }
func Created() *StdReply { return NewReply().Status(http.StatusCreated) }

func (reply *StdReply) Status(status int) *StdReply {
	reply.response.Status = status
	return reply
}

func (reply *StdReply) Header(key, value string) *StdReply {
	reply.response.Headers[key] = value
	return reply
}

func (reply *StdReply) BodyStr(value string) *StdReply {
	reply.response.Body = strings.NewReader(value)
	return reply
}

func (reply *StdReply) Build() Responder {
	return func(_ *http.Request, _ *Mock) (*Response, error) { return reply.response, nil }
}

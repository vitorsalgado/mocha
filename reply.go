package mocha

import (
	"fmt"
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

	ResponseDelegate func(r *http.Request, mock *Mock) (*Response, error)

	Reply interface {
		Build() ResponseDelegate
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

func (reply *StdReply) Fault() *StdReply {
	reply.response.Body = &F{}
	return reply
}

func (reply *StdReply) BodyStr(value string) *StdReply {
	reply.response.Body = strings.NewReader(value)
	return reply
}

func (reply *StdReply) Build() ResponseDelegate {
	return func(_ *http.Request, _ *Mock) (*Response, error) { return reply.response, nil }
}

type F struct {
}

func (f F) Read([]byte) (int, error) {
	return 0, fmt.Errorf("test error")
}

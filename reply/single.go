package reply

import (
	"bytes"
	"encoding/json"
	"github.com/vitorsalgado/mocha/mock"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type SingleReply struct {
	err      error
	response *mock.Response
}

func New() *SingleReply {
	return &SingleReply{
		response: &mock.Response{
			Cookies: make([]*http.Cookie, 0),
			Header:  make(http.Header)}}
}

func OK() *SingleReply                  { return New().Status(http.StatusOK) }
func Created() *SingleReply             { return New().Status(http.StatusCreated) }
func Accepted() *SingleReply            { return New().Status(http.StatusAccepted) }
func NoContent() *SingleReply           { return New().Status(http.StatusNoContent) }
func PartialContent() *SingleReply      { return New().Status(http.StatusPartialContent) }
func MovedPermanently() *SingleReply    { return New().Status(http.StatusMovedPermanently) }
func NotModified() *SingleReply         { return New().Status(http.StatusNotModified) }
func BadRequest() *SingleReply          { return New().Status(http.StatusBadRequest) }
func Unauthorized() *SingleReply        { return New().Status(http.StatusUnauthorized) }
func Forbidden() *SingleReply           { return New().Status(http.StatusForbidden) }
func NotFound() *SingleReply            { return New().Status(http.StatusNotFound) }
func MethodNotAllowed() *SingleReply    { return New().Status(http.StatusMethodNotAllowed) }
func UnprocessableEntity() *SingleReply { return New().Status(http.StatusUnprocessableEntity) }
func MultipleChoices() *SingleReply     { return New().Status(http.StatusMultipleChoices) }
func InternalServerError() *SingleReply { return New().Status(http.StatusInternalServerError) }
func NotImplemented() *SingleReply      { return New().Status(http.StatusNotImplemented) }
func BadGateway() *SingleReply          { return New().Status(http.StatusBadGateway) }
func ServiceUnavailable() *SingleReply  { return New().Status(http.StatusServiceUnavailable) }
func GatewayTimeout() *SingleReply      { return New().Status(http.StatusGatewayTimeout) }

func (reply *SingleReply) Status(status int) *SingleReply {
	reply.response.Status = status
	return reply
}

func (reply *SingleReply) Header(key, value string) *SingleReply {
	reply.response.Header.Add(key, value)
	return reply
}

func (reply *SingleReply) Cookie(cookie http.Cookie) *SingleReply {
	reply.response.Cookies = append(reply.response.Cookies, &cookie)
	return reply
}

func (reply *SingleReply) RemoveCookie(cookie http.Cookie) *SingleReply {
	cookie.MaxAge = -1
	reply.response.Cookies = append(reply.response.Cookies, &cookie)
	return reply
}

func (reply *SingleReply) Body(value []byte) *SingleReply {
	reply.response.Body = value
	return reply
}

func (reply *SingleReply) BodyString(value string) *SingleReply {
	reply.response.Body = []byte(value)
	return reply
}

func (reply *SingleReply) BodyJSON(data any) *SingleReply {
	buf := &bytes.Buffer{}
	reply.response.Err = json.NewEncoder(buf).Encode(data)
	return reply
}

func (reply *SingleReply) BodyReader(reader io.Reader) *SingleReply {
	b, err := ioutil.ReadAll(reader)

	reply.response.Body = b
	reply.err = err

	return reply
}

func (reply *SingleReply) Delay(duration time.Duration) *SingleReply {
	reply.response.Delay = duration
	return reply
}

func (reply *SingleReply) Err() error {
	return reply.err
}

func (reply *SingleReply) Build(_ *http.Request, _ *mock.Mock) (*mock.Response, error) {
	return reply.response, reply.err
}

package reply

import (
	"bytes"
	"encoding/json"
	"github.com/vitorsalgado/mocha/mock"
	"github.com/vitorsalgado/mocha/templating"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type (
	SingleReply struct {
		err      error
		response *mock.Response
		bodyType BodyType
		template string
	}

	BodyType int
)

const (
	BodyTemplate BodyType = iota
)

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

func (r *SingleReply) Status(status int) *SingleReply {
	r.response.Status = status
	return r
}

func (r *SingleReply) Header(key, value string) *SingleReply {
	r.response.Header.Add(key, value)
	return r
}

func (r *SingleReply) Cookie(cookie http.Cookie) *SingleReply {
	r.response.Cookies = append(r.response.Cookies, &cookie)
	return r
}

func (r *SingleReply) RemoveCookie(cookie http.Cookie) *SingleReply {
	cookie.MaxAge = -1
	r.response.Cookies = append(r.response.Cookies, &cookie)
	return r
}

func (r *SingleReply) Body(value []byte) *SingleReply {
	r.response.Body = value
	return r
}

func (r *SingleReply) BodyString(value string) *SingleReply {
	r.response.Body = []byte(value)
	return r
}

func (r *SingleReply) BodyTemplate(tmpl string) *SingleReply {
	return r
}

func (r *SingleReply) BodyJSON(data any) *SingleReply {
	buf := &bytes.Buffer{}
	r.response.Err = json.NewEncoder(buf).Encode(data)
	return r
}

func (r *SingleReply) BodyReader(reader io.Reader) *SingleReply {
	b, err := ioutil.ReadAll(reader)

	r.response.Body = b
	r.err = err

	return r
}

func (r *SingleReply) Delay(duration time.Duration) *SingleReply {
	r.response.Delay = duration
	return r
}

func (r *SingleReply) Err() error {
	return r.err
}

func (r *SingleReply) Build(_ *http.Request, _ *mock.Mock) (*mock.Response, error) {
	if r.err != nil {
		return nil, r.err
	}

	switch r.bodyType {
	case BodyTemplate:
		tmpl := templating.New()
		tmpl.Template(r.template).Name("")
	}

	return r.response, r.err
}

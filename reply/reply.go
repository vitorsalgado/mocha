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
		template templating.Template
		model    any
	}

	BodyType int
)

const (
	BodyDefault BodyType = iota
	BodyTemplate
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

func (rpl *SingleReply) Status(status int) *SingleReply {
	rpl.response.Status = status
	return rpl
}

func (rpl *SingleReply) Header(key, value string) *SingleReply {
	rpl.response.Header.Add(key, value)
	return rpl
}

func (rpl *SingleReply) Cookie(cookie http.Cookie) *SingleReply {
	rpl.response.Cookies = append(rpl.response.Cookies, &cookie)
	return rpl
}

func (rpl *SingleReply) RemoveCookie(cookie http.Cookie) *SingleReply {
	cookie.MaxAge = -1
	rpl.response.Cookies = append(rpl.response.Cookies, &cookie)
	return rpl
}

func (rpl *SingleReply) Body(value []byte) *SingleReply {
	rpl.response.Body = value
	return rpl
}

func (rpl *SingleReply) BodyString(value string) *SingleReply {
	rpl.response.Body = []byte(value)
	return rpl
}

func (rpl *SingleReply) BodyJSON(data any) *SingleReply {
	buf := &bytes.Buffer{}
	rpl.response.Err = json.NewEncoder(buf).Encode(data)
	return rpl
}

func (rpl *SingleReply) BodyReader(reader io.Reader) *SingleReply {
	b, err := ioutil.ReadAll(reader)

	rpl.response.Body = b
	rpl.err = err

	return rpl
}

func (rpl *SingleReply) BodyTemplate(template any) *SingleReply {
	switch e := template.(type) {
	case string:
		rpl.template = templating.New().Template(e)
	case templating.Template:
		err := e.Compile()
		rpl.template = e
		rpl.err = err
	case *templating.Template:
		rpl.template = *e
	default:
		panic(".BodyTemplate() parameter must be: string | templating.Template")
	}

	return rpl
}

func (rpl *SingleReply) Mode(model any) *SingleReply {
	rpl.model = model
	return rpl
}

func (rpl *SingleReply) Delay(duration time.Duration) *SingleReply {
	rpl.response.Delay = duration
	return rpl
}

func (rpl *SingleReply) Err() error {
	return rpl.err
}

func (rpl *SingleReply) Build(_ *http.Request, _ *mock.Mock) (*mock.Response, error) {
	if rpl.err != nil {
		return nil, rpl.err
	}

	switch rpl.bodyType {
	case BodyTemplate:
		b, err := rpl.template.Parse(rpl.model)
		rpl.response.Body = b
		rpl.err = err
	}

	return rpl.response, rpl.err
}

package reply

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/vitorsalgado/mocha/internal/params"
	"github.com/vitorsalgado/mocha/mock"
	"github.com/vitorsalgado/mocha/templating"
)

type (
	// StdReply holds the configuration on how the mock.Response should be built.
	StdReply struct {
		err      error
		response *mock.Response
		bodyType bodyType
		template templating.Template
		model    any
	}

	bodyType int
)

const (
	bodyDefault bodyType = iota
	bodyTemplate
)

// New creates a new StdReply. Prefer to use factory functions for each status code.
func New() *StdReply {
	return &StdReply{
		response: &mock.Response{
			Cookies: make([]*http.Cookie, 0),
			Header:  make(http.Header),
			Mappers: make([]mock.ResponseMapper, 0),
		},
	}
}

// OK creates a new mock.Reply with http.StatusOK already.
func OK() *StdReply { return New().Status(http.StatusOK) }

// Created creates a new mock.Reply with http.StatusCreated already.
func Created() *StdReply { return New().Status(http.StatusCreated) }

// Accepted creates a new mock.Reply with http.StatusAccepted already.
func Accepted() *StdReply { return New().Status(http.StatusAccepted) }

// NoContent creates a new mock.Reply with http.StatusNoContent already.
func NoContent() *StdReply { return New().Status(http.StatusNoContent) }

// PartialContent creates a new mock.Reply with http.StatusPartialContent already.
func PartialContent() *StdReply { return New().Status(http.StatusPartialContent) }

// MovedPermanently creates a new mock.Reply with http.StatusMovedPermanently already.
func MovedPermanently() *StdReply { return New().Status(http.StatusMovedPermanently) }

// NotModified creates a new mock.Reply with http.StatusNotModified already.
func NotModified() *StdReply { return New().Status(http.StatusNotModified) }

// BadRequest creates a new mock.Reply with http.StatusBadRequest already.
func BadRequest() *StdReply { return New().Status(http.StatusBadRequest) }

// Unauthorized creates a new mock.Reply with http.StatusUnauthorized already.
func Unauthorized() *StdReply { return New().Status(http.StatusUnauthorized) }

// Forbidden creates a new mock.Reply with http.StatusForbidden already.
func Forbidden() *StdReply { return New().Status(http.StatusForbidden) }

// NotFound creates a new mock.Reply with http.StatusNotFound already.
func NotFound() *StdReply { return New().Status(http.StatusNotFound) }

// MethodNotAllowed creates a new mock.Reply with http.StatusMethodNotAllowed already.
func MethodNotAllowed() *StdReply { return New().Status(http.StatusMethodNotAllowed) }

// UnprocessableEntity creates a new mock.Reply with http.StatusUnprocessableEntity already.
func UnprocessableEntity() *StdReply { return New().Status(http.StatusUnprocessableEntity) }

// MultipleChoices creates a new mock.Reply with http.StatusMultipleChoices already.
func MultipleChoices() *StdReply { return New().Status(http.StatusMultipleChoices) }

// InternalServerError creates a new mock.Reply with http.StatusInternalServerError already.
func InternalServerError() *StdReply { return New().Status(http.StatusInternalServerError) }

// NotImplemented creates a new mock.Reply with http.StatusNotImplemented already.
func NotImplemented() *StdReply { return New().Status(http.StatusNotImplemented) }

// BadGateway creates a new mock.Reply with http.StatusBadGateway already.
func BadGateway() *StdReply { return New().Status(http.StatusBadGateway) }

// ServiceUnavailable creates a new mock.Reply with http.StatusServiceUnavailable already.
func ServiceUnavailable() *StdReply { return New().Status(http.StatusServiceUnavailable) }

// GatewayTimeout creates a new mock.Reply with http.StatusGatewayTimeout already.
func GatewayTimeout() *StdReply { return New().Status(http.StatusGatewayTimeout) }

// Status sets the HTTP status code for the mock.Response.
func (rpl *StdReply) Status(status int) *StdReply {
	rpl.response.Status = status
	return rpl
}

// Header adds a header to the mock.Response.
func (rpl *StdReply) Header(key, value string) *StdReply {
	rpl.response.Header.Add(key, value)
	return rpl
}

// Cookie adds a http.Cookie to the mock.Response.
func (rpl *StdReply) Cookie(cookie http.Cookie) *StdReply {
	rpl.response.Cookies = append(rpl.response.Cookies, &cookie)
	return rpl
}

// RemoveCookie expires a cookie.
func (rpl *StdReply) RemoveCookie(cookie http.Cookie) *StdReply {
	cookie.MaxAge = -1
	rpl.response.Cookies = append(rpl.response.Cookies, &cookie)
	return rpl
}

// Body defines the response body using a []byte,
func (rpl *StdReply) Body(value []byte) *StdReply {
	rpl.response.Body = bytes.NewReader(value)
	return rpl
}

// BodyString defines the response body using a string.
func (rpl *StdReply) BodyString(value string) *StdReply {
	rpl.response.Body = strings.NewReader(value)
	return rpl
}

// BodyJSON defines the response body encoding the given value using json.Encoder.
func (rpl *StdReply) BodyJSON(data any) *StdReply {
	buf := &bytes.Buffer{}
	rpl.response.Err = json.NewEncoder(buf).Encode(data)
	return rpl
}

// BodyReader defines the response body using the given io.Reader.
func (rpl *StdReply) BodyReader(reader io.Reader) *StdReply {
	rpl.response.Body = reader
	return rpl
}

// BodyTemplate defines the response body using a template.
// It accepts a string or a templating.Template implementation. If a different type is provided, it panics.
func (rpl *StdReply) BodyTemplate(template any) *StdReply {
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
		panic(".bodyTemplate() parameter must be: string | templating.Template")
	}

	return rpl
}

// Model sets the template data to be used.
func (rpl *StdReply) Model(model any) *StdReply {
	rpl.model = model
	return rpl
}

// Delay sets a delay time before serving the stub mock.Response.
func (rpl *StdReply) Delay(duration time.Duration) *StdReply {
	rpl.response.Delay = duration
	return rpl
}

// Map adds mock.ResponseMapper that will be executed after the mock.Response was built.
func (rpl *StdReply) Map(mapper mock.ResponseMapper) *StdReply {
	rpl.response.Mappers = append(rpl.response.Mappers, mapper)
	return rpl
}

// Build builds a mock.Response based on StdReply definition.
func (rpl *StdReply) Build(_ *http.Request, _ *mock.Mock, _ params.Params) (*mock.Response, error) {
	if rpl.err != nil {
		return nil, rpl.err
	}

	switch rpl.bodyType {
	case bodyTemplate:
		buf := &bytes.Buffer{}
		err := rpl.template.Parse(buf, rpl.model)

		rpl.response.Body = buf
		rpl.err = err
	}

	return rpl.response, rpl.err
}

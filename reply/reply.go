package reply

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/vitorsalgado/mocha/v3/params"
)

type (
	// M implements mock data that should be available on reply build functions.
	M interface {
		// Hits return mock total hits.
		Hits() int
	}

	// Reply defines the contract to configure an HTTP responder.
	Reply interface {
		// Build returns a Response stub to be served.
		Build(*http.Request, M, params.P) (*Response, error)
	}

	// StdReply holds the configuration on how the Response should be built.
	StdReply struct {
		response *Response
		bodyType bodyType
		template Template
		model    any
		err      error
	}

	bodyType int
)

const (
	_bodyDefault bodyType = iota
	_bodyTemplate
)

// New creates a new StdReply. Prefer to use factory functions for each status code.
func New() *StdReply {
	return &StdReply{
		response: &Response{
			Cookies: make([]*http.Cookie, 0),
			Header:  make(http.Header),
			Mappers: make([]ResponseMapper, 0),
		},
		bodyType: _bodyDefault,
	}
}

// Status creates a new Reply with the given HTTP status code.
func Status(status int) *StdReply { return New().Status(status) }

// OK creates a new Reply with http.StatusOK already.
func OK() *StdReply { return New().Status(http.StatusOK) }

// Created creates a new Reply with http.StatusCreated already.
func Created() *StdReply { return New().Status(http.StatusCreated) }

// Accepted creates a new Reply with http.StatusAccepted already.
func Accepted() *StdReply { return New().Status(http.StatusAccepted) }

// NoContent creates a new Reply with http.StatusNoContent already.
func NoContent() *StdReply { return New().Status(http.StatusNoContent) }

// PartialContent creates a new Reply with http.StatusPartialContent already.
func PartialContent() *StdReply { return New().Status(http.StatusPartialContent) }

// MovedPermanently creates a new Reply with http.StatusMovedPermanently already.
func MovedPermanently() *StdReply { return New().Status(http.StatusMovedPermanently) }

// NotModified creates a new Reply with http.StatusNotModified already.
func NotModified() *StdReply { return New().Status(http.StatusNotModified) }

// BadRequest creates a new Reply with http.StatusBadRequest already.
func BadRequest() *StdReply { return New().Status(http.StatusBadRequest) }

// Unauthorized creates a new Reply with http.StatusUnauthorized already.
func Unauthorized() *StdReply { return New().Status(http.StatusUnauthorized) }

// Forbidden creates a new Reply with http.StatusForbidden already.
func Forbidden() *StdReply { return New().Status(http.StatusForbidden) }

// NotFound creates a new Reply with http.StatusNotFound already.
func NotFound() *StdReply { return New().Status(http.StatusNotFound) }

// MethodNotAllowed creates a new Reply with http.StatusMethodNotAllowed already.
func MethodNotAllowed() *StdReply { return New().Status(http.StatusMethodNotAllowed) }

// UnprocessableEntity creates a new Reply with http.StatusUnprocessableEntity already.
func UnprocessableEntity() *StdReply { return New().Status(http.StatusUnprocessableEntity) }

// MultipleChoices creates a new Reply with http.StatusMultipleChoices already.
func MultipleChoices() *StdReply { return New().Status(http.StatusMultipleChoices) }

// InternalServerError creates a new Reply with http.StatusInternalServerError already.
func InternalServerError() *StdReply { return New().Status(http.StatusInternalServerError) }

// NotImplemented creates a new Reply with http.StatusNotImplemented already.
func NotImplemented() *StdReply { return New().Status(http.StatusNotImplemented) }

// BadGateway creates a new Reply with http.StatusBadGateway already.
func BadGateway() *StdReply { return New().Status(http.StatusBadGateway) }

// ServiceUnavailable creates a new Reply with http.StatusServiceUnavailable already.
func ServiceUnavailable() *StdReply { return New().Status(http.StatusServiceUnavailable) }

// GatewayTimeout creates a new Reply with http.StatusGatewayTimeout already.
func GatewayTimeout() *StdReply { return New().Status(http.StatusGatewayTimeout) }

// Status sets the HTTP status code for the Response.
func (rpl *StdReply) Status(status int) *StdReply {
	rpl.response.Status = status
	return rpl
}

// Header adds a header to the Response.
func (rpl *StdReply) Header(key, value string) *StdReply {
	rpl.response.Header.Add(key, value)
	return rpl
}

// Cookie adds a http.Cookie to the Response.
func (rpl *StdReply) Cookie(cookie http.Cookie) *StdReply {
	rpl.response.Cookies = append(rpl.response.Cookies, &cookie)
	return rpl
}

// ExpireCookie expires a cookie.
func (rpl *StdReply) ExpireCookie(cookie http.Cookie) *StdReply {
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
	err := json.NewEncoder(buf).Encode(data)
	if err != nil {
		rpl.err = err
		return rpl
	}

	rpl.response.Body = buf

	return rpl
}

// BodyReader defines the response body using the given io.Reader.
func (rpl *StdReply) BodyReader(reader io.Reader) *StdReply {
	rpl.response.Body = reader
	return rpl
}

// BodyTemplate defines the response body using a template.
// It accepts a string or a reply.Template implementation. If a different type is provided, it panics.
func (rpl *StdReply) BodyTemplate(template any) *StdReply {
	switch e := template.(type) {
	case string:
		rpl.template = NewTextTemplate().Template(e)
	case Template:
		rpl.err = e.Compile()
		rpl.template = e

	default:
		panic(".BodyTemplate() parameter must be: string | reply.Template")
	}

	rpl.bodyType = _bodyTemplate

	return rpl
}

// Model sets the template data to be used.
func (rpl *StdReply) Model(model any) *StdReply {
	rpl.model = model
	return rpl
}

// Delay sets a delay time before serving the stub Response.
func (rpl *StdReply) Delay(duration time.Duration) *StdReply {
	rpl.response.Delay = duration
	return rpl
}

// Map adds ResponseMapper that will be executed after the Response was built.
func (rpl *StdReply) Map(mapper ResponseMapper) *StdReply {
	rpl.response.Mappers = append(rpl.response.Mappers, mapper)
	return rpl
}

// Build builds a Response based on StdReply definition.
func (rpl *StdReply) Build(r *http.Request, _ M, _ params.P) (*Response, error) {
	if rpl.err != nil {
		return nil, rpl.err
	}

	switch rpl.bodyType {
	case _bodyTemplate:
		buf := &bytes.Buffer{}
		model := &TemplateData{Request: r, Data: rpl.model}
		err := rpl.template.Parse(buf, model)
		if err != nil {
			return nil, err
		}

		rpl.response.Body = buf
	}

	return rpl.response, nil
}

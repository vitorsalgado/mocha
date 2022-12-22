package reply

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/internal/mimetype"
	"github.com/vitorsalgado/mocha/v3/types"
)

var _ Reply = (*StdReply)(nil)

// StdReply holds the configuration on how the Stub should be built.
type StdReply struct {
	response       *Stub
	bodyType       bodyType
	template       Template
	templateExtras any
	err            error
}

type bodyType int

const (
	_bodyDefault bodyType = iota
	_bodyTemplate
)

// New creates a new StdReply. Prefer to use factory functions for each status code.
func New() *StdReply {
	return &StdReply{
		response: &Stub{Cookies: make([]*http.Cookie, 0), Header: make(http.Header)},
		bodyType: _bodyDefault}
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

// Status sets the HTTP status code for the Stub.
func (rep *StdReply) Status(status int) *StdReply {
	rep.response.StatusCode = status
	return rep
}

// Header adds a header to the Stub.
func (rep *StdReply) Header(key, value string) *StdReply {
	rep.response.Header.Add(key, value)
	return rep
}

// Cookie adds a http.Cookie to the Stub.
func (rep *StdReply) Cookie(cookie *http.Cookie) *StdReply {
	rep.response.Cookies = append(rep.response.Cookies, cookie)
	return rep
}

// ExpireCookie expires a cookie.
func (rep *StdReply) ExpireCookie(cookie http.Cookie) *StdReply {
	cookie.MaxAge = -1
	rep.response.Cookies = append(rep.response.Cookies, &cookie)
	return rep
}

// Body defines the response body using a []byte,
func (rep *StdReply) Body(value []byte) *StdReply {
	rep.response.Body = value
	return rep
}

// BodyJSON defines the response body encoding the given value using json.Encoder.
func (rep *StdReply) BodyJSON(data any) *StdReply {
	b, err := json.Marshal(data)
	if err != nil {
		rep.err = err
		return rep
	}

	rep.response.Body = b
	rep.Header(header.ContentType, mimetype.JSON)

	return rep
}

// BodyReader defines the response body using the given io.Reader.
func (rep *StdReply) BodyReader(reader io.Reader) *StdReply {
	b, err := io.ReadAll(reader)
	if err != nil {
		rep.err = err
		return rep
	}

	rep.response.Body = b

	return rep
}

// BodyTemplate defines the response body using a template.
// It accepts a string or a reply.Template implementation. If a different type is provided, it panics.
func (rep *StdReply) BodyTemplate(tpl any, extras any) *StdReply {
	switch e := tpl.(type) {
	case string:
		rep.template = NewTextTemplate().Template(e)
	case Template:
		rep.err = e.Compile()
		rep.template = e

	default:
		panic(".BodyTemplate() parameter must be: string | reply.Template")
	}

	rep.bodyType = _bodyTemplate
	rep.templateExtras = extras

	return rep
}

// JSON sets the response to application/json.
func (rep *StdReply) JSON() *StdReply {
	rep.Header(header.ContentType, mimetype.JSON)
	return rep
}

// PlainText defines a text/plain response with the given text body.
func (rep *StdReply) PlainText(value string) *StdReply {
	rep.response.Body = []byte(value)
	rep.Header(header.ContentType, mimetype.TextPlain)
	return rep
}

func (rep *StdReply) Prepare() error {
	if rep.err != nil {
		return rep.err
	}

	return nil
}

func (rep *StdReply) Spec() []any {
	return []any{"response", map[string]any{
		"status":  rep.response.StatusCode,
		"header":  rep.response.Header,
		"body":    string(rep.response.Body),
		"cookies": fmt.Sprintf("%v", rep.response.Cookies),
	}}
}

// Build builds a Stub based on StdReply definition.
func (rep *StdReply) Build(_ http.ResponseWriter, r *types.RequestValues) (*Stub, error) {
	if rep.err != nil {
		return nil, rep.err
	}

	switch rep.bodyType {
	case _bodyTemplate:
		buf := &bytes.Buffer{}
		reqExtra := templateRequest{r.RawRequest.Method, *r.URL, r.RawRequest.Header.Clone(), r.Body}
		model := &templateData{Request: reqExtra, Extras: rep.templateExtras}

		err := rep.template.Render(buf, model)
		if err != nil {
			return nil, err
		}

		rep.response.Body = buf.Bytes()
	}

	return rep.response, nil
}

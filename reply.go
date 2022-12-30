package mocha

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/internal/mimetype"
)

// Reply defines the contract to configure an HTTP responder.
type Reply interface {
	// Build returns an HTTP response Stub to be served.
	// Return Stub nil if the HTTP response was rendered inside the Build function.
	Build(w http.ResponseWriter, r *RequestValues) (*Stub, error)
}

// Pre describes a Reply that has preparations steps to run on mocking building.
type Pre interface {
	// Pre runs once during mock building.
	// Useful for pre-configurations or validations that needs to be executed once.
	Pre() error
}

// Stub defines the HTTP response that will be served once a Mock is matched for an HTTP Request.
type Stub struct {
	StatusCode int
	Header     http.Header
	Cookies    []*http.Cookie
	Body       io.Reader
	Trailer    http.Header
}

func (s *Stub) bodyBytes() ([]byte, error) {
	return io.ReadAll(s.Body)
}

// -- Standard Reply

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
	_bodyGZIP
)

// NewReply creates a new StdReply. Prefer to use factory functions for each status code.
func NewReply() *StdReply {
	return &StdReply{
		response: &Stub{Cookies: make([]*http.Cookie, 0), Header: make(http.Header), Trailer: make(http.Header)},
		bodyType: _bodyDefault}
}

// Status creates a new Reply with the given HTTP status code.
func Status(status int) *StdReply { return NewReply().Status(status) }

// OK creates a new Reply with http.StatusOK already.
func OK() *StdReply { return NewReply().Status(http.StatusOK) }

// Created creates a new Reply with http.StatusCreated already.
func Created() *StdReply { return NewReply().Status(http.StatusCreated) }

// Accepted creates a new Reply with http.StatusAccepted already.
func Accepted() *StdReply { return NewReply().Status(http.StatusAccepted) }

// NoContent creates a new Reply with http.StatusNoContent already.
func NoContent() *StdReply { return NewReply().Status(http.StatusNoContent) }

// PartialContent creates a new Reply with http.StatusPartialContent already.
func PartialContent() *StdReply { return NewReply().Status(http.StatusPartialContent) }

// MovedPermanently creates a new Reply with http.StatusMovedPermanently already.
func MovedPermanently() *StdReply { return NewReply().Status(http.StatusMovedPermanently) }

// NotModified creates a new Reply with http.StatusNotModified already.
func NotModified() *StdReply { return NewReply().Status(http.StatusNotModified) }

// BadRequest creates a new Reply with http.StatusBadRequest already.
func BadRequest() *StdReply { return NewReply().Status(http.StatusBadRequest) }

// Unauthorized creates a new Reply with http.StatusUnauthorized already.
func Unauthorized() *StdReply { return NewReply().Status(http.StatusUnauthorized) }

// Forbidden creates a new Reply with http.StatusForbidden already.
func Forbidden() *StdReply { return NewReply().Status(http.StatusForbidden) }

// NotFound creates a new Reply with http.StatusNotFound already.
func NotFound() *StdReply { return NewReply().Status(http.StatusNotFound) }

// MethodNotAllowed creates a new Reply with http.StatusMethodNotAllowed already.
func MethodNotAllowed() *StdReply { return NewReply().Status(http.StatusMethodNotAllowed) }

// UnprocessableEntity creates a new Reply with http.StatusUnprocessableEntity already.
func UnprocessableEntity() *StdReply { return NewReply().Status(http.StatusUnprocessableEntity) }

// MultipleChoices creates a new Reply with http.StatusMultipleChoices already.
func MultipleChoices() *StdReply { return NewReply().Status(http.StatusMultipleChoices) }

// InternalServerError creates a new Reply with http.StatusInternalServerError already.
func InternalServerError() *StdReply { return NewReply().Status(http.StatusInternalServerError) }

// NotImplemented creates a new Reply with http.StatusNotImplemented already.
func NotImplemented() *StdReply { return NewReply().Status(http.StatusNotImplemented) }

// BadGateway creates a new Reply with http.StatusBadGateway already.
func BadGateway() *StdReply { return NewReply().Status(http.StatusBadGateway) }

// ServiceUnavailable creates a new Reply with http.StatusServiceUnavailable already.
func ServiceUnavailable() *StdReply { return NewReply().Status(http.StatusServiceUnavailable) }

// GatewayTimeout creates a new Reply with http.StatusGatewayTimeout already.
func GatewayTimeout() *StdReply { return NewReply().Status(http.StatusGatewayTimeout) }

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

// Trailer adds a trailer header to the response Stub.
func (rep *StdReply) Trailer(key, value string) *StdReply {
	rep.response.Trailer.Add(key, value)
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

// Body defines the response body using a []byte.
func (rep *StdReply) Body(value []byte) *StdReply {
	rep.response.Body = bytes.NewReader(value)
	return rep
}

// BodyText defines the response body using a string.
func (rep *StdReply) BodyText(text string) *StdReply {
	rep.response.Body = strings.NewReader(text)
	return rep
}

// BodyJSON defines the response body encoding the given value using json.Encoder.
func (rep *StdReply) BodyJSON(data any) *StdReply {
	b, err := json.Marshal(data)
	if err != nil {
		rep.err = err
		return rep
	}

	rep.response.Body = bytes.NewReader(b)

	return rep
}

// BodyReader defines the response body using the given io.Reader.
func (rep *StdReply) BodyReader(reader io.Reader) *StdReply {
	rep.response.Body = reader
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
func (rep *StdReply) JSON(payload any) *StdReply {
	rep.Header(header.ContentType, mimetype.JSON)
	return rep.BodyJSON(payload)
}

// PlainText defines a text/plain response with the given text body.
func (rep *StdReply) PlainText(text string) *StdReply {
	rep.Header(header.ContentType, mimetype.TextPlain)
	return rep.BodyText(text)
}

// Gzip indicates that the response should be gzip encoded.
func (rep *StdReply) Gzip() *StdReply {
	rep.bodyType = _bodyGZIP
	return rep
}

func (rep *StdReply) Pre() error {
	if rep.err != nil {
		return rep.err
	}

	switch rep.bodyType {
	case _bodyGZIP:
		b, err := io.ReadAll(rep.response.Body)
		if err != nil {
			return err
		}

		buf := new(bytes.Buffer)
		gz := gzip.NewWriter(buf)

		_, err = gz.Write(b)
		if err != nil {
			return err
		}

		err = gz.Close()
		if err != nil {
			return err
		}

		rep.response.Body = buf
	}

	return nil
}

// Build builds a Stub based on StdReply definition.
func (rep *StdReply) Build(_ http.ResponseWriter, r *RequestValues) (*Stub, error) {
	if rep.err != nil {
		return nil, rep.err
	}

	switch rep.bodyType {
	case _bodyTemplate:
		buf := &bytes.Buffer{}
		reqExtra := templateRequest{r.RawRequest.Method, *r.URL, r.RawRequest.Header.Clone(), r.ParsedBody}
		model := &templateData{Request: reqExtra, Extras: rep.templateExtras}

		err := rep.template.Render(buf, model)
		if err != nil {
			return nil, err
		}

		rep.response.Body = buf
	}

	return rep.response, nil
}

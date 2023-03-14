package mocha

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/textproto"
	"os"

	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/internal/mimetype"
)

// Reply defines the contract to set up an HTTP replier.
type Reply interface {
	// Build returns an HTTP response Stub to be served; after the HTTP request was matched.
	// Return a nil Stub if the HTTP response was rendered inside the Build function.
	Build(w http.ResponseWriter, r *RequestValues) (*Stub, error)
}

type replyOnBeforeBuild interface {
	beforeBuild(app *Mocha) error
}

var _ Reply = (*StdReply)(nil)

// StdReply holds the configuration on how the Stub should be built.
type StdReply struct {
	response            *Stub
	bodyType            bodyType
	bodyEncoding        bodyEncoding
	bodyFilename        string
	bodyTeRender        TemplateRenderer
	bodyFnTeRender      TemplateRenderer
	headerTeRender      TemplateRenderer
	teType              teType
	teHeader            http.Header
	bodyTemplateContent string
	teData              any
	err                 error
	encoded             bool
}

type bodyType int

const (
	_bodyDefault bodyType = iota
	_bodyTemplate
	_bodyFile
)

type bodyEncoding int

const (
	_bodyEncodingNone bodyEncoding = iota
	_bodyEncodingGZIP
)

type teType byte

const (
	_teBody teType = 1 << iota
	_teBodyFilename
	_teHeader
)

// NewReply creates a new StdReply. Prefer to use factory functions for each status code.
func NewReply() *StdReply {
	return &StdReply{
		response:     newStub(),
		bodyType:     _bodyDefault,
		bodyEncoding: _bodyEncodingNone,
		teHeader:     make(http.Header)}
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
func MovedPermanently(location string) *StdReply {
	return NewReply().Status(http.StatusMovedPermanently).Header(header.Location, location)
}

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

// HeaderTemplate adds a header that can have its value changed using templates.
func (rep *StdReply) HeaderTemplate(key, value string) *StdReply {
	rep.teHeader.Add(key, value)
	rep.teType += _teHeader
	return rep
}

// ContentType sets the response content-type header.
func (rep *StdReply) ContentType(mime string) *StdReply {
	rep.Header(header.ContentType, mime)
	return rep
}

// Trailer adds a trailer header to the response.
func (rep *StdReply) Trailer(key, value string) *StdReply {
	rep.response.Trailer.Add(key, value)
	return rep
}

// Cookie adds a http.Cookie to the response.
func (rep *StdReply) Cookie(cookie *http.Cookie) *StdReply {
	rep.response.Cookies = append(rep.response.Cookies, cookie)
	return rep
}

// ExpireCookie expires a cookie.
func (rep *StdReply) ExpireCookie(cookie *http.Cookie) *StdReply {
	cookie.MaxAge = -1
	rep.response.Cookies = append(rep.response.Cookies, cookie)
	return rep
}

// Body defines the response body using a []byte.
func (rep *StdReply) Body(value []byte) *StdReply {
	rep.response.Body = value
	return rep
}

// BodyText defines the response body using a string.
func (rep *StdReply) BodyText(text string) *StdReply {
	rep.response.Body = []byte(text)
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

	return rep
}

// BodyReader defines the response body using the given io.Reader.
func (rep *StdReply) BodyReader(reader io.Reader) *StdReply {
	b, err := io.ReadAll(reader)
	if err != nil {
		rep.err = err
	}

	rep.response.Body = b

	return rep
}

// BodyFile loads the body content from the given filename.
func (rep *StdReply) BodyFile(filename string) *StdReply {
	rep.bodyFilename = filename
	rep.bodyType = _bodyFile
	rep.teType += _teBodyFilename

	return rep
}

// BodyTemplate defines the response body using a template.
// The parameter content must be the actual template content that should be parsed and rendered.
// Use BodyFileWithTemplate function to load the template from a file.
func (rep *StdReply) BodyTemplate(content string) *StdReply {
	rep.bodyTemplateContent = content
	rep.bodyType = _bodyTemplate
	rep.teType += _teBody

	return rep
}

// BodyFileWithTemplate loads the body content from the given filename.
func (rep *StdReply) BodyFileWithTemplate(filename string) *StdReply {
	rep.bodyFilename = filename
	rep.bodyType = _bodyFile
	rep.teType = _teBodyFilename | _teBody

	return rep
}

// SetTemplateData sets the data model to be used by templates during response building.
func (rep *StdReply) SetTemplateData(data any) *StdReply {
	rep.teData = data
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
	rep.bodyEncoding = _bodyEncodingGZIP
	return rep
}

func (rep *StdReply) beforeBuild(app *Mocha) error {
	if rep.err != nil {
		return rep.err
	}

	if rep.teType&_teBodyFilename == _teBodyFilename {
		r, err := app.te.Parse(rep.bodyFilename)
		if err != nil {
			return err
		}

		rep.bodyFnTeRender = r
	}

	if rep.teType&_teBody == _teBody {
		r, err := app.te.Parse(rep.bodyTemplateContent)
		if err != nil {
			return err
		}

		rep.bodyTeRender = r
	}

	if rep.teType&_teHeader == _teHeader {
		buf := &bytes.Buffer{}
		err := rep.teHeader.Write(buf)
		if err != nil {
			return err
		}

		r, err := app.te.Parse(buf.String() + "\r\n")
		if err != nil {
			return err
		}

		rep.headerTeRender = r
	}

	err := rep.encodeBody()
	if err != nil {
		return err
	}

	return nil
}

// Build builds a Stub based on StdReply definition.
func (rep *StdReply) Build(w http.ResponseWriter, r *RequestValues) (stub *Stub, err error) {
	stub, err = rep.build(w, r)
	if err != nil {
		return nil, fmt.Errorf("reply: %w", err)
	}

	return stub, nil
}

func (rep *StdReply) build(_ http.ResponseWriter, r *RequestValues) (stub *Stub, err error) {
	if rep.err != nil {
		return nil, rep.err
	}

	switch rep.bodyType {
	case _bodyTemplate:
		defer func() {
			if recovered := recover(); recovered != nil {
				err = fmt.Errorf(
					"panic parsing body template. reason=%v",
					recovered,
				)
			}
		}()

		buf := &bytes.Buffer{}
		err = rep.bodyTeRender.Render(buf, rep.buildTemplateData(r, rep.teData))
		if err != nil {
			return nil, err
		}

		rep.response.Body = buf.Bytes()

	case _bodyFile:
		defer func() {
			if recovered := recover(); recovered != nil {
				err = fmt.Errorf(
					"panic loading body from %s. reason=%v",
					rep.bodyFilename,
					recovered,
				)
			}
		}()

		buf := &bytes.Buffer{}
		err = rep.bodyFnTeRender.Render(buf, rep.buildTemplateData(r, rep.teData))
		if err != nil {
			return nil, err
		}

		filename := buf.String()
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}

		defer file.Close()

		b, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}

		if rep.teType&_teBody == _teBody {
			buf.Reset()
			t, err := r.App.TemplateEngine().Parse(string(b))
			if err != nil {
				return nil, err
			}

			err = t.Render(buf, rep.buildTemplateData(r, rep.teData))
			if err != nil {
				return nil, err
			}

			rep.response.Body = buf.Bytes()
		} else {
			rep.response.Body = b
		}
	}

	if rep.teType&_teHeader == _teHeader {
		buf := &bytes.Buffer{}
		err = rep.headerTeRender.Render(buf, rep.buildTemplateData(r, rep.teData))
		if err != nil {
			return nil, err
		}

		tp := textproto.NewReader(bufio.NewReader(buf))
		mimeHeader, err := tp.ReadMIMEHeader()
		if err != nil {
			return nil, err
		}

		h := http.Header(mimeHeader)
		for k, v := range h {
			for _, vv := range v {
				rep.Header(k, vv)
			}
		}
	}

	err = rep.encodeBody()
	if err != nil {
		return nil, err
	}

	return rep.response, nil
}

func (rep *StdReply) encodeBody() error {
	if rep.encoded || len(rep.response.Body) == 0 {
		return nil
	}

	switch rep.bodyEncoding {
	case _bodyEncodingGZIP:
		buf := new(bytes.Buffer)
		gz := gzip.NewWriter(buf)

		_, err := gz.Write(rep.response.Body)
		if err != nil {
			return err
		}

		err = gz.Close()
		if err != nil {
			return err
		}

		rep.response.Body = buf.Bytes()
		rep.response.Header.Add(header.ContentEncoding, "gzip")
		rep.response.Encoding = "gzip"
		rep.encoded = true
	}

	return nil
}

func (rep *StdReply) buildTemplateData(r *RequestValues, ext any) *templateData {
	reqExtra := templateRequest{
		Method:          r.RawRequest.Method,
		URL:             r.URL,
		URLPathSegments: r.URLPathSegments,
		Header:          r.RawRequest.Header.Clone(),
		Cookies:         r.RawRequest.Cookies(),
		Body:            r.ParsedBody,
	}

	return &templateData{Request: reqExtra, App: r.App, Ext: ext}
}

package dzhttp

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
	"sync"

	"github.com/vitorsalgado/mocha/v3/dzhttp/httpval"
	"github.com/vitorsalgado/mocha/v3/dzstd"
)

var _ Reply = (*StdReply)(nil)

const maxBufSize = 1 << 16

var bufPool = &sync.Pool{
	New: func() any { return bytes.NewBuffer(make([]byte, 500)) },
}

func putBuf(b *bytes.Buffer) {
	if b.Cap() > maxBufSize {
		return
	}

	bufPool.Put(b)
}

// Reply defines the contract to set up an HTTP replier.
type Reply interface {
	// Build returns an HTTP response MockedResponse to be served; after the HTTP request was matched.
	// Return a nil MockedResponse if the HTTP response was rendered inside the Build function.
	Build(w http.ResponseWriter, r *RequestValues) (*MockedResponse, error)
}

type replyOnBeforeBuild interface {
	beforeBuild(app *HTTPMockApp) error
}

// StdReply holds the configuration on how the MockedResponse should be built.
type StdReply struct {
	response                 MockedResponse
	bodyType                 bodyType
	bodyEncoding             bodyEncoding
	bodyFilename             string
	bodyTemplateRenderer     dzstd.TemplateRenderer
	bodyFuncTemplateRenderer dzstd.TemplateRenderer
	headerTemplateRender     dzstd.TemplateRenderer
	teType                   teType
	templateHeader           http.Header
	bodyTemplateContent      string
	teData                   any
	delayedErr               error
	encoded                  bool
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

var gzipper = &sync.Pool{New: func() any { return gzip.NewWriter(nil) }}

// NewReply creates a new StdReply. Prefer to use factory functions for each status code.
func NewReply() *StdReply {
	return &StdReply{
		response:       *newResponse(),
		bodyType:       _bodyDefault,
		bodyEncoding:   _bodyEncodingNone,
		templateHeader: make(http.Header),
	}
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
	return NewReply().Status(http.StatusMovedPermanently).Header(httpval.HeaderLocation, location)
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

// Status sets the HTTP status code for the MockedResponse.
func (rep *StdReply) Status(status int) *StdReply {
	rep.response.StatusCode = status
	return rep
}

// Header adds a header to the MockedResponse.
func (rep *StdReply) Header(key, value string) *StdReply {
	rep.response.Header.Add(key, value)
	return rep
}

// HeaderArr adds a header with multiple values to the MockedResponse.
func (rep *StdReply) HeaderArr(key string, values ...string) *StdReply {
	for _, value := range values {
		rep.Header(key, value)
	}

	return rep
}

// HeaderTemplate adds a header that can have its value changed using templates.
func (rep *StdReply) HeaderTemplate(key, value string) *StdReply {
	rep.templateHeader.Add(key, value)
	rep.teType += _teHeader
	return rep
}

// HeaderArrTemplate adds a header with multiple values that can have its value changed using templates.
func (rep *StdReply) HeaderArrTemplate(key string, values ...string) *StdReply {
	for _, value := range values {
		rep.HeaderTemplate(key, value)
	}
	rep.teType += _teHeader
	return rep
}

// ContentType sets the response content-type misc.Header
func (rep *StdReply) ContentType(mime string) *StdReply {
	rep.Header(httpval.HeaderContentType, mime)
	return rep
}

// Trailer adds a trailer header to the response.
func (rep *StdReply) Trailer(key, value string) *StdReply {
	rep.response.Trailer.Add(key, value)
	return rep
}

// TrailerArr adds a trailer with multiple values to the response.
func (rep *StdReply) TrailerArr(key string, values ...string) *StdReply {
	rep.response.Trailer[key] = values
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
		rep.delayedErr = err
		return rep
	}

	rep.response.Body = b

	return rep
}

// BodyReader defines the response body using the given io.Reader.
func (rep *StdReply) BodyReader(reader io.Reader) *StdReply {
	b, err := io.ReadAll(reader)
	if err != nil {
		rep.delayedErr = err
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
	rep.Header(httpval.HeaderContentType, httpval.MIMEApplicationJSON)
	return rep.BodyJSON(payload)
}

// PlainText defines a text/plain response with the given text body.
func (rep *StdReply) PlainText(text string) *StdReply {
	rep.Header(httpval.HeaderContentType, httpval.MIMETextPlain)
	return rep.BodyText(text)
}

// Gzip indicates that the response should be gzip encoded.
func (rep *StdReply) Gzip() *StdReply {
	rep.bodyEncoding = _bodyEncodingGZIP
	return rep
}

func (rep *StdReply) beforeBuild(app *HTTPMockApp) error {
	if rep.delayedErr != nil {
		return rep.delayedErr
	}

	if rep.teType&_teBodyFilename == _teBodyFilename {
		r, err := app.templateEngine.Parse(rep.bodyFilename)
		if err != nil {
			return err
		}

		rep.bodyFuncTemplateRenderer = r
	}

	if rep.teType&_teBody == _teBody {
		r, err := app.templateEngine.Parse(rep.bodyTemplateContent)
		if err != nil {
			return err
		}

		rep.bodyTemplateRenderer = r
	}

	if len(rep.templateHeader) > 0 {
		buf := bufPool.Get().(*bytes.Buffer)
		buf.Reset()

		err := rep.templateHeader.Write(buf)
		if err != nil {
			return err
		}

		r, err := app.templateEngine.Parse(buf.String() + "\r\n")
		if err != nil {
			return err
		}

		rep.headerTemplateRender = r

		putBuf(buf)
	}

	if rep.encoded {
		return nil
	}

	switch rep.bodyEncoding {
	case _bodyEncodingGZIP:
		rep.response.Header.Add(httpval.HeaderContentEncoding, "gzip")
		rep.response.Encoding = "gzip"
		rep.encoded = true
	}

	return nil
}

// Build builds a MockedResponse based on StdReply definition.
func (rep *StdReply) Build(w http.ResponseWriter, r *RequestValues) (stub *MockedResponse, err error) {
	stub, err = rep.build(w, r)
	if err != nil {
		return nil, fmt.Errorf("reply: %w", err)
	}

	return stub, nil
}

func (rep *StdReply) build(_ http.ResponseWriter, r *RequestValues) (stub *MockedResponse, err error) {
	if rep.delayedErr != nil {
		return nil, rep.delayedErr
	}

	if len(rep.templateHeader) > 0 {
		buf := bufPool.Get().(*bytes.Buffer)
		buf.Reset()

		err = rep.headerTemplateRender.Render(buf, rep.buildTemplateData(r, rep.teData))
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
			rep.HeaderArr(k, v...)
		}

		putBuf(buf)
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

		buf := bufPool.Get().(*bytes.Buffer)
		buf.Reset()

		defer putBuf(buf)

		err = rep.bodyTemplateRenderer.Render(buf, rep.buildTemplateData(r, rep.teData))
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
		buf := bufPool.Get().(*bytes.Buffer)
		buf.Reset()

		defer putBuf(buf)

		filename := rep.bodyFilename
		if rep.teType&_teBodyFilename == _teBodyFilename {
			err = rep.bodyFuncTemplateRenderer.Render(buf, rep.buildTemplateData(r, rep.teData))
			if err != nil {
				return nil, err
			}

			filename = buf.String()
		}

		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}

		if rep.teType&_teBody == _teBody {
			defer file.Close()

			buf.Reset()

			b, err := io.ReadAll(file)
			if err != nil {
				return nil, err
			}

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
			rep.response.BodyCloser = file
		}

		rep.response.BodyFilename = filename
	}

	return &rep.response, nil
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

	return &templateData{Request: reqExtra, App: &templateAppWrapper{r.App}, Ext: ext}
}

func (rep *StdReply) Describe() any {
	headers := rep.response.Header.Clone()
	for _, cookie := range rep.response.Cookies {
		headers.Add("Set-Cookie", cookie.String())
	}

	response := map[string]any{"status": rep.response.StatusCode}

	if len(headers) > 0 {
		response["header"] = headers
	}

	if len(rep.response.Trailer) > 0 {
		response["trailers"] = rep.response.Trailer.Clone()
	}

	if len(rep.response.Encoding) > 0 {
		response["encoding"] = rep.response.Encoding
	}

	if len(rep.bodyFilename) > 0 {
		response["body_file"] = rep.bodyFilename
	} else if rep.response.Body != nil {
		response["body"] = string(rep.response.Body)
	}

	switch rep.bodyType {
	case _bodyTemplate:
		tmpl := map[string]any{"enabled": true}

		if rep.teData != nil {
			tmpl["data"] = rep.teData
		}

		response["template"] = tmpl
	}

	return map[string]any{"response": response}
}

package httprec

import (
	"io"
	"net/http"
	"net/http/httptest"
)

var _ http.ResponseWriter = (*HTTPRec)(nil)

type HTTPRec struct {
	wrapped  http.ResponseWriter
	recorder *httptest.ResponseRecorder
	w        io.Writer
}

func Wrap(w http.ResponseWriter) *HTTPRec {
	rec := httptest.NewRecorder()
	return &HTTPRec{wrapped: w, recorder: rec, w: io.MultiWriter(rec, w)}
}

func (r *HTTPRec) Header() http.Header {
	return r.wrapped.Header()
}

func (r *HTTPRec) Write(buf []byte) (int, error) {
	return r.w.Write(buf)
}

func (r *HTTPRec) WriteHeader(statusCode int) {
	for k, v := range r.wrapped.Header() {
		for _, vv := range v {
			r.recorder.Header().Add(k, vv)
		}
	}

	r.recorder.WriteHeader(statusCode)
	r.wrapped.WriteHeader(statusCode)
}

func (r *HTTPRec) Result() *http.Response {
	return r.recorder.Result()
}

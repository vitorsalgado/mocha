package httprec

import (
	"net/http"
	"net/http/httptest"
)

var _ http.ResponseWriter = (*HTTPRec)(nil)

type HTTPRec struct {
	Wrapped  http.ResponseWriter
	recorder *httptest.ResponseRecorder
}

func Wrap(w http.ResponseWriter) *HTTPRec {
	return &HTTPRec{Wrapped: w, recorder: httptest.NewRecorder()}
}

func (r *HTTPRec) Header() http.Header {
	return r.Wrapped.Header()
}

func (r *HTTPRec) Write(buf []byte) (int, error) {
	_, err := r.recorder.Write(buf)
	if err != nil {
		return 0, err
	}

	return r.Wrapped.Write(buf)
}

func (r *HTTPRec) WriteHeader(statusCode int) {
	for k, v := range r.Wrapped.Header() {
		for _, vv := range v {
			r.recorder.Header().Add(k, vv)
		}
	}

	r.recorder.WriteHeader(statusCode)
	r.Wrapped.WriteHeader(statusCode)
}

func (r *HTTPRec) Result() *http.Response {
	return r.recorder.Result()
}

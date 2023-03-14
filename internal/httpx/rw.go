package httpx

import (
	"net/http"
	"net/http/httptest"
)

var _ http.ResponseWriter = (*RRec)(nil)

type RRec struct {
	Wrapped  http.ResponseWriter
	recorder *httptest.ResponseRecorder
}

func Wrap(w http.ResponseWriter) *RRec {
	return &RRec{Wrapped: w, recorder: httptest.NewRecorder()}
}

func (r *RRec) Header() http.Header {
	return r.Wrapped.Header()
}

func (r *RRec) Write(buf []byte) (int, error) {
	_, err := r.recorder.Write(buf)
	if err != nil {
		return 0, err
	}

	return r.Wrapped.Write(buf)
}

func (r *RRec) WriteHeader(statusCode int) {
	for k, v := range r.Wrapped.Header() {
		for _, vv := range v {
			r.recorder.Header().Add(k, vv)
		}
	}

	r.recorder.WriteHeader(statusCode)
	r.Wrapped.WriteHeader(statusCode)
}

func (r *RRec) Result() *http.Response {
	return r.recorder.Result()
}

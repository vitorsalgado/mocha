package httpx

import (
	"net/http"
	"net/http/httptest"
)

var _ http.ResponseWriter = (*Rw)(nil)

type Rw struct {
	rr *httptest.ResponseRecorder
	w  http.ResponseWriter
}

func DecorateWriter(w http.ResponseWriter) *Rw {
	return &Rw{w: w, rr: httptest.NewRecorder()}
}

func (rw *Rw) Header() http.Header {
	return rw.w.Header()
}

func (rw *Rw) Write(buf []byte) (int, error) {
	_, err := rw.rr.Write(buf)
	if err != nil {
		return 0, err
	}

	return rw.w.Write(buf)
}

func (rw *Rw) WriteHeader(statusCode int) {
	rw.rr.WriteHeader(statusCode)
	rw.w.WriteHeader(statusCode)
}

func (rw *Rw) Result() *http.Response {
	return rw.rr.Result()
}

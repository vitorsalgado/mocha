package httpx

import (
	"net/http"
	"net/http/httptest"
)

var _ http.ResponseWriter = (*Rw)(nil)

type Rw struct {
	recorder     *httptest.ResponseRecorder
	actualWriter http.ResponseWriter
}

func Wrap(w http.ResponseWriter) *Rw {
	return &Rw{actualWriter: w, recorder: httptest.NewRecorder()}
}

func (rw *Rw) Header() http.Header {
	return rw.actualWriter.Header()
}

func (rw *Rw) Write(buf []byte) (int, error) {
	_, err := rw.recorder.Write(buf)
	if err != nil {
		return 0, err
	}

	return rw.actualWriter.Write(buf)
}

func (rw *Rw) WriteHeader(statusCode int) {
	for k, v := range rw.actualWriter.Header() {
		for _, vv := range v {
			rw.recorder.Header().Add(k, vv)
		}
	}

	rw.recorder.WriteHeader(statusCode)
	rw.actualWriter.WriteHeader(statusCode)
}

func (rw *Rw) Result() *http.Response {
	return rw.recorder.Result()
}

package mhttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctionReply(t *testing.T) {
	fn := func(http.ResponseWriter, *RequestValues) (*Stub, error) {
		return &Stub{StatusCode: http.StatusAccepted}, nil
	}

	r := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	rv := &RequestValues{RawRequest: r, URL: r.URL}
	replier := Function(fn)
	res, err := replier.Build(nil, rv)

	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusAccepted)
}

func TestHandlerReply(t *testing.T) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}

	r := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	rv := &RequestValues{RawRequest: r, URL: r.URL}
	w := httptest.NewRecorder()

	replier := Handler(fn)
	res, err := replier.Build(w, rv)

	assert.NoError(t, err)
	assert.Nil(t, res)
	assert.Equal(t, http.StatusCreated, w.Code)
}

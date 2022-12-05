package reply

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerReply(t *testing.T) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}

	r := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	w := httptest.NewRecorder()

	replier := Handler(fn)
	res, err := replier.Build(w, r)

	assert.NoError(t, err)
	assert.Nil(t, res)
	assert.Equal(t, http.StatusCreated, w.Code)
}

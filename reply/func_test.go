package reply

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctionReply(t *testing.T) {
	fn := func(http.ResponseWriter, *http.Request) (*Response, error) {
		return &Response{Status: http.StatusAccepted}, nil
	}

	r := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	replier := Function(fn)
	res, err := replier.Build(nil, r)

	assert.Nil(t, err)
	assert.Equal(t, res.Status, http.StatusAccepted)
}

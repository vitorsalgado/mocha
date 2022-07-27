package reply

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/internal/parameters"
)

func TestFunctionReply(t *testing.T) {
	fn := func(*http.Request, M, parameters.Params) (*Response, error) {
		return &Response{Status: http.StatusAccepted}, nil
	}

	r := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	replier := Function(fn)
	m := &mmock{}
	m.On("Hits").Return(0)
	res, err := replier.Build(r, m, nil)

	assert.Nil(t, err)
	assert.Equal(t, res.Status, http.StatusAccepted)
}

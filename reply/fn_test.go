package reply

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorsalgado/mocha/internal/params"
	"github.com/vitorsalgado/mocha/mock"
)

func TestFunctionReply(t *testing.T) {
	fn := func(*http.Request, *mock.Mock, params.Params) (*mock.Response, error) {
		return &mock.Response{Status: http.StatusAccepted}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	replier := Function(fn)
	m := mock.New()
	res, err := replier.Build(req, m, nil)

	assert.Nil(t, err)
	assert.Equal(t, res.Status, http.StatusAccepted)
}

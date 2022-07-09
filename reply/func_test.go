package reply

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/core"
	"github.com/vitorsalgado/mocha/internal/parameters"
)

func TestFunctionReply(t *testing.T) {
	fn := func(*http.Request, *core.Mock, parameters.Params) (*core.Response, error) {
		return &core.Response{Status: http.StatusAccepted}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	replier := Function(fn)
	m := core.NewMock()
	res, err := replier.Build(req, m, nil)

	assert.Nil(t, err)
	assert.Equal(t, res.Status, http.StatusAccepted)
}

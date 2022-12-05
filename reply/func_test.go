package reply

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctionReply(t *testing.T) {
	fn := func(*http.Request, M, Params) (*Response, error) {
		return &Response{Status: http.StatusAccepted}, nil
	}

	r := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	replier := Function(fn)
	m := &mMock{}
	m.On("Hits").Return(0)
	res, err := replier.Build(r, m, nil)

	assert.Nil(t, err)
	assert.Equal(t, res.Status, http.StatusAccepted)
}

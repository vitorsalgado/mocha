package reply

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3/types"
)

func TestFunctionReply(t *testing.T) {
	fn := func(http.ResponseWriter, *types.RequestValues) (*Stub, error) {
		return &Stub{StatusCode: http.StatusAccepted}, nil
	}

	r := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	replier := Function(fn)
	res, err := replier.Build(nil, newReqValues(r))

	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusAccepted)
}

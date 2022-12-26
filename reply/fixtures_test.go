package reply

import (
	"net/http"

	"github.com/vitorsalgado/mocha/v3/types"
)

var _req, _ = http.NewRequest(http.MethodGet, "http://localhost:8080", nil)

func newReqValues(req *http.Request) *types.RequestValues {
	return &types.RequestValues{RawRequest: req, URL: req.URL}
}

package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/reply"
)

func TestCORS(t *testing.T) {
	m := mocha.New(t, mocha.Configure().CORS().Build())
	m.Start()

	defer m.Close()

	m.AddMocks(mocha.Get(matcher.URLPath("/test")).
		Reply(reply.OK()))

	corsReq := testutil.NewRequest(http.MethodOptions, m.URL()+"/test", nil)
	res, err := corsReq.Do()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	req := testutil.NewRequest(http.MethodGet, m.URL()+"/test", nil)
	res, err = req.Do()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

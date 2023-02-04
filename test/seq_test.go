package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/expect"
	"github.com/vitorsalgado/mocha/v3/reply"
)

func TestSeqReply(t *testing.T) {
	m := mocha.New(t)
	m.Start()

	m.AddMocks(mocha.Get(expect.URLPath("/test")).Reply(reply.Seq().Add(reply.Unauthorized(), reply.OK())))

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)

	res, err := http.DefaultClient.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, 401, res.StatusCode)

	res, err = http.DefaultClient.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)

	res, err = http.DefaultClient.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, 418, res.StatusCode)
}

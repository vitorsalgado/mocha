package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3"
	. "github.com/vitorsalgado/mocha/v3/matcher"
)

func TestSequenceReplies(t *testing.T) {
	m := mocha.New()
	m.MustStart()
	m.MustMock(mocha.Get(URLPath("/test")).Reply(mocha.Seq().Add(mocha.Unauthorized(), mocha.OK())))

	defer m.Close()

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)

	res, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	res, err = http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, err = http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusTeapot, res.StatusCode)
}

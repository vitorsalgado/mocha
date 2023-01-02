package test

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

func TestTLS(t *testing.T) {
	m := mocha.New()
	m.MustStartTLS()

	defer m.Close()

	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}

	scoped := m.MustMock(mocha.Get(matcher.URLPath("/test")).
		Header("test", matcher.StrictEqual("hello")).
		Reply(mocha.OK()))

	req, err := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)
	require.NoError(t, err)

	req.Header.Add("test", "hello")

	res, err := client.Do(req)

	assert.NoError(t, err)
	assert.NoError(t, res.Body.Close())
	assert.True(t, scoped.HasBeenCalled())
}

package test

import (
	"testing"
)

func TestTLS(t *testing.T) {
	// m := mocha.New(t)
	// m.MustStartTLS()
	//
	// defer m.Close()
	//
	// // allow insecure https request
	// http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	//
	// scoped := m.MustMock(mocha.Get(matcher.URLPath("/test")).
	// 	Header("test", matcher.Equal("hello")).
	// 	Reply(reply.OK()))
	//
	// req := testutil.Get(m.URL() + "/test")
	// req.Header("test", "hello")
	//
	// res, err := req.Do()
	//
	// assert.NoError(t, err)
	// assert.NoError(t, res.Body.Close())
	// assert.True(t, scoped.Called())
}

package test

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha"
	"github.com/vitorsalgado/mocha/expect"
	"github.com/vitorsalgado/mocha/reply"
	"github.com/vitorsalgado/mocha/util/headers"
	"github.com/vitorsalgado/mocha/util/mimetypes"
)

func TestForward(t *testing.T) {
	dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "ok", r.Header.Get("x-test"))
		assert.Equal(t, "", r.Header.Get("x-del"))
		assert.Equal(t, mimetypes.TextPlain, r.Header.Get(headers.ContentType))

		w.WriteHeader(http.StatusOK)

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		w.Write(b)
	}))

	defer dest.Close()

	m := mocha.ForTest(t)
	m.Start()

	t.Run("should forward request and respond using proxied response and mock definition", func(t *testing.T) {
		scoped := m.Mock(mocha.Post(expect.URLPath("/test")).
			Reply(reply.
				From(dest.URL).
				ProxyHeader("x-test", "ok").
				Header("x-res", "example").
				RemoveProxyHeader("x-del")))

		data := strings.NewReader("hello world")
		req, _ := http.NewRequest(http.MethodPost, m.URL()+"/test", data)
		req.Header.Add("x-del", "to-delete")
		req.Header.Add(headers.ContentType, mimetypes.TextPlain)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		defer res.Body.Close()

		scoped.MustHaveBeenCalled(t)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "example", res.Header.Get("x-res"))
		assert.Equal(t, "hello world", string(b))
	})
}

package httpclient

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient(t *testing.T) {
	client := New(Options{
		Timeout:             3 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		KeepAlive:           30 * time.Second,
		DialTimeout:         1 * time.Second,
	})

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	}))

	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "hello world", string(body))
}

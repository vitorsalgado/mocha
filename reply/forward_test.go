package reply

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForward(t *testing.T) {
	t.Run("should forward and respond basic GET", func(t *testing.T) {
		dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/path/test/example", r.URL.Path)
			assert.Equal(t, "all", r.URL.Query().Get("filter"))

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("hello world"))
		}))

		defer dest.Close()

		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/path/test/example?filter=all", nil)

		forward := From(dest.URL)
		res, err := forward.Build(req, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		assert.NotNil(t, res)
		assert.Equal(t, http.StatusCreated, res.Status)
		assert.Equal(t, "hello world", string(b))
	})

	t.Run("should forward and respond POST with body", func(t *testing.T) {
		dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}

			w.Write(b)
		}))

		defer dest.Close()

		expected := "test text"
		body := strings.NewReader(expected)
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", body)

		forward := ProxiedFrom(dest.URL)
		res, err := forward.Build(req, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		assert.NotNil(t, res)
		assert.Equal(t, http.StatusOK, res.Status)
		assert.Equal(t, expected, string(b))
	})

	t.Run("should forward and respond a No Content", func(t *testing.T) {
		dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))

		defer dest.Close()

		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)

		forward := From(dest.URL)
		res, err := forward.Build(req, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		assert.NotNil(t, res)
		assert.Equal(t, http.StatusNoContent, res.Status)
		assert.Equal(t, "", string(b))
	})
}

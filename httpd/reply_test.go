package httpd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _req, _ = http.NewRequest(http.MethodGet, "http://localhost:8080", nil)

func TestReplyFactories(t *testing.T) {
	assert.Equal(t, http.StatusOK, Status(http.StatusOK).response.StatusCode)
	assert.Equal(t, http.StatusOK, OK().response.StatusCode)
	assert.Equal(t, http.StatusCreated, Created().response.StatusCode)
	assert.Equal(t, http.StatusAccepted, Accepted().response.StatusCode)
	assert.Equal(t, http.StatusNoContent, NoContent().response.StatusCode)
	assert.Equal(t, http.StatusPartialContent, PartialContent().response.StatusCode)
	assert.Equal(t, http.StatusMovedPermanently, MovedPermanently("https://nowhere.com.br").response.StatusCode)
	assert.Equal(t, http.StatusNotModified, NotModified().response.StatusCode)
	assert.Equal(t, http.StatusBadRequest, BadRequest().response.StatusCode)
	assert.Equal(t, http.StatusUnauthorized, Unauthorized().response.StatusCode)
	assert.Equal(t, http.StatusForbidden, Forbidden().response.StatusCode)
	assert.Equal(t, http.StatusNotFound, NotFound().response.StatusCode)
	assert.Equal(t, http.StatusMethodNotAllowed, MethodNotAllowed().response.StatusCode)
	assert.Equal(t, http.StatusUnprocessableEntity, UnprocessableEntity().response.StatusCode)
	assert.Equal(t, http.StatusMultipleChoices, MultipleChoices().response.StatusCode)
	assert.Equal(t, http.StatusInternalServerError, InternalServerError().response.StatusCode)
	assert.Equal(t, http.StatusNotImplemented, NotImplemented().response.StatusCode)
	assert.Equal(t, http.StatusBadGateway, BadGateway().response.StatusCode)
	assert.Equal(t, http.StatusServiceUnavailable, ServiceUnavailable().response.StatusCode)
	assert.Equal(t, http.StatusGatewayTimeout, GatewayTimeout().response.StatusCode)
}

func TestReply(t *testing.T) {
	rv := &RequestValues{RawRequest: _req, URL: _req.URL}
	res, err := NewReply().
		Status(http.StatusCreated).
		Header("test", "dev").
		Header("test", "qa").
		Header("hello", "world").
		Cookie(&http.Cookie{Name: "cookie_test"}).
		ExpireCookie(&http.Cookie{Name: "cookie_test_remove"}).
		Body([]byte("hi")).
		Build(nil, rv)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, []string{"dev", "qa"}, res.Header.Values("test"))
	assert.Equal(t, "world", res.Header.Get("hello"))
	assert.Equal(t, 2, len(res.Cookies))
	assert.Equal(t, "cookie_test", res.Cookies[0].Name)
	assert.Equal(t, "cookie_test_remove", res.Cookies[1].Name)
	assert.Equal(t, -1, res.Cookies[1].MaxAge)
	assert.Equal(t, "hi", string(res.Body))
}

func TestStdReplyBodyString(t *testing.T) {
	rv := &RequestValues{RawRequest: _req, URL: _req.URL}
	res, err := NewReply().
		Status(http.StatusCreated).
		PlainText("text").
		Build(nil, rv)

	require.NoError(t, err)
	require.Equal(t, "text", string(res.Body))
}

func TestStdReplyBodyJSON(t *testing.T) {
	type jsonData struct {
		Name   string `json:"name"`
		Job    string `json:"job"`
		Active bool   `json:"active"`
	}

	rv := &RequestValues{RawRequest: _req, URL: _req.URL}

	t.Run("should convert struct to json", func(t *testing.T) {
		model := jsonData{
			Name:   "the name",
			Job:    "dev",
			Active: true,
		}

		res, err := NewReply().
			Status(http.StatusCreated).
			BodyJSON(model).
			Build(nil, rv)

		assert.NoError(t, err)

		b := jsonData{}
		err = json.Unmarshal(res.Body, &b)

		assert.Nil(t, err)
		assert.Equal(t, model, b)
	})

	t.Run("should report conversion error", func(t *testing.T) {
		res, err := NewReply().
			Status(http.StatusCreated).
			BodyJSON(make(chan int)).
			Build(nil, rv)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestStdReplyBodyReader(t *testing.T) {
	buf := new(bytes.Buffer)
	buf.WriteString("hello\n")
	buf.WriteString("world\n")

	rv := &RequestValues{RawRequest: _req, URL: _req.URL}
	res, err := NewReply().
		Status(http.StatusCreated).
		BodyReader(buf).
		Build(nil, rv)

	assert.NoError(t, err)
	assert.Equal(t, "hello\nworld\n", string(res.Body))
}

func BenchmarkReply(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		NewReply().
			Status(http.StatusCreated).
			Header("test", "dev").
			Header("test", "qa").
			Header("hello", "world").
			Cookie(&http.Cookie{Name: "cookie_test"}).
			ExpireCookie(&http.Cookie{Name: "cookie_test_remove"}).
			Body([]byte("hi")).
			Build(nil, &RequestValues{RawRequest: _req, URL: _req.URL})
	}
}

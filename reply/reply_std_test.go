package reply

import (
	"encoding/json"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplyFactories(t *testing.T) {
	assert.Equal(t, http.StatusOK, Status(http.StatusOK).response.Status)
	assert.Equal(t, http.StatusOK, OK().response.Status)
	assert.Equal(t, http.StatusCreated, Created().response.Status)
	assert.Equal(t, http.StatusAccepted, Accepted().response.Status)
	assert.Equal(t, http.StatusNoContent, NoContent().response.Status)
	assert.Equal(t, http.StatusPartialContent, PartialContent().response.Status)
	assert.Equal(t, http.StatusMovedPermanently, MovedPermanently().response.Status)
	assert.Equal(t, http.StatusNotModified, NotModified().response.Status)
	assert.Equal(t, http.StatusBadRequest, BadRequest().response.Status)
	assert.Equal(t, http.StatusUnauthorized, Unauthorized().response.Status)
	assert.Equal(t, http.StatusForbidden, Forbidden().response.Status)
	assert.Equal(t, http.StatusNotFound, NotFound().response.Status)
	assert.Equal(t, http.StatusMethodNotAllowed, MethodNotAllowed().response.Status)
	assert.Equal(t, http.StatusUnprocessableEntity, UnprocessableEntity().response.Status)
	assert.Equal(t, http.StatusMultipleChoices, MultipleChoices().response.Status)
	assert.Equal(t, http.StatusInternalServerError, InternalServerError().response.Status)
	assert.Equal(t, http.StatusNotImplemented, NotImplemented().response.Status)
	assert.Equal(t, http.StatusBadGateway, BadGateway().response.Status)
	assert.Equal(t, http.StatusServiceUnavailable, ServiceUnavailable().response.Status)
	assert.Equal(t, http.StatusGatewayTimeout, GatewayTimeout().response.Status)
}

func TestReply(t *testing.T) {
	res, err := New().
		Status(http.StatusCreated).
		Header("test", "dev").
		Header("test", "qa").
		Header("hello", "world").
		Cookie(&http.Cookie{Name: "cookie_test"}).
		ExpireCookie(http.Cookie{Name: "cookie_test_remove"}).
		Body([]byte("hi")).
		Build(nil, _req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.Status)
	assert.Equal(t, []string{"dev", "qa"}, res.Header.Values("test"))
	assert.Equal(t, "world", res.Header.Get("hello"))
	assert.Equal(t, 2, len(res.Cookies))
	assert.Equal(t, "cookie_test", res.Cookies[0].Name)
	assert.Equal(t, "cookie_test_remove", res.Cookies[1].Name)
	assert.Equal(t, -1, res.Cookies[1].MaxAge)
	assert.Equal(t, "hi", string(res.Body))
}

func TestStdReply_BodyString(t *testing.T) {
	res, err := New().
		Status(http.StatusCreated).
		PlainText("text").
		Build(nil, _req)

	assert.NoError(t, err)
	assert.Equal(t, "text", string(res.Body))
}

func TestStdReply_BodyJSON(t *testing.T) {
	type jsonData struct {
		Name   string `json:"name"`
		Job    string `json:"job"`
		Active bool   `json:"active"`
	}

	t.Run("should convert struct to json", func(t *testing.T) {
		model := jsonData{
			Name:   "the name",
			Job:    "dev",
			Active: true,
		}

		res, err := New().
			Status(http.StatusCreated).
			BodyJSON(model).
			Build(nil, _req)

		assert.Nil(t, err)

		b := jsonData{}
		err = json.Unmarshal(res.Body, &b)

		assert.Nil(t, err)
		assert.Equal(t, model, b)
	})

	t.Run("should report conversion error", func(t *testing.T) {
		res, err := New().
			Status(http.StatusCreated).
			BodyJSON(make(chan int)).
			Build(nil, _req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}

func TestStdReply_BodyReader(t *testing.T) {
	wd, _ := os.Getwd()
	f, err := os.Open(path.Join(wd, "testdata", "data.txt"))
	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	res, err := New().
		Status(http.StatusCreated).
		BodyReader(f).
		Build(nil, _req)

	assert.NoError(t, err)
	assert.Equal(t, "hello\nworld\n", string(res.Body))
}

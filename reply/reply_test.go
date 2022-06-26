package reply

import (
	"github.com/stretchr/testify/assert"
	"github.com/vitorsalgado/mocha/mock"
	"net/http"
	"testing"
	"time"
)

func TestSingleFactories(t *testing.T) {
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

func TestSingleReplies(t *testing.T) {
	m := mock.Mock{Name: "mock_test"}
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)

	res, err := New().
		Status(http.StatusCreated).
		Header("test", "dev").
		Header("test", "qa").
		Header("hello", "world").
		Cookie(http.Cookie{Name: "cookie_test"}).
		RemoveCookie(http.Cookie{Name: "cookie_test_remove"}).
		Body([]byte("hi")).
		Delay(5*time.Second).
		Build(req, &m)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.Status)
	assert.Equal(t, []string{"dev", "qa"}, res.Header.Values("test"))
	assert.Equal(t, "world", res.Header.Get("hello"))
	assert.Equal(t, 2, len(res.Cookies))
	assert.Equal(t, "cookie_test", res.Cookies[0].Name)
	assert.Equal(t, "cookie_test_remove", res.Cookies[1].Name)
	assert.Equal(t, -1, res.Cookies[1].MaxAge)
	assert.Equal(t, "hi", string(res.Body))
	assert.Equal(t, 5*time.Second, res.Delay)
}

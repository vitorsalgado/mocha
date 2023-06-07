package application

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vitorsalgado/mocha/v3"
	. "github.com/vitorsalgado/mocha/v3/matcher"
	mhttp2 "github.com/vitorsalgado/mocha/v3/mhttp"
)

func TestHandlerGetById(t *testing.T) {
	m := mocha.New().CloseWithT(t)
	m.MustStart()

	id := "super-id"
	customer := Customer{ID: id, Name: "nice-name"}

	m.MustMock(mhttp2.Get(URLPathf("/customers/%s", id)).
		Header(headerAccept, StrictEqual(contentTypeJSON)).
		Header(headerContentType, StrictEqual(contentTypeJSON)).
		Reply(mhttp2.OK().BodyJSON(customer)))

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/customer/%s", id), nil)
	rr := httptest.NewRecorder()
	handler := Handler{api: CustomerAPI{base: m.URL()}}
	h := http.HandlerFunc(handler.GetById)

	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code to be %d. got %d", status, http.StatusOK)
		t.FailNow()
	}

	expected := `{"id":"super-id","name":"nice-name"}`
	if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(expected) {
		t.Errorf("expected body %s. got %s", expected, rr.Body.String())
		t.FailNow()
	}
}

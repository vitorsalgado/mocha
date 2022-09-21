package application

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/expect"
	"github.com/vitorsalgado/mocha/v3/reply"
)

func TestHandler_GetById(t *testing.T) {
	m := mocha.New(t).CloseOnCleanup(t)
	m.Start()

	id := "super-id"
	customer := Customer{ID: id, Name: "nice-name"}

	m.AddMocks(mocha.
		Get(expect.URLPath(fmt.Sprintf("/customers/%s", id))).
		Header(headerAccept, expect.ToEqual(contentTypeJSON)).
		Header(headerContentType, expect.ToEqual(contentTypeJSON)).
		Reply(reply.OK().BodyJSON(customer)))

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

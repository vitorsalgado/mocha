package mocha

import (
	"github.com/vitorsalgado/mocha/mock"
	"log"
	"net/http"

	"github.com/vitorsalgado/mocha/matcher"
)

type Handler struct {
	mocks   mock.Storage
	parsers []BodyParser
	extras  Extras
}

func newHandler(
	mockstore mock.Storage,
	parsers []BodyParser,
	extras Extras,
) *Handler {
	return &Handler{mocks: mockstore, parsers: parsers, extras: extras}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parsedbody, err := ParseRequestBody(r, h.parsers)
	if err != nil {
		respondErr(w, err)
		return
	}

	params := matcher.Params{RequestInfo: &matcher.RequestInfo{Request: r, ParsedBody: parsedbody}, Extras: &h.extras}
	result, err := findMockForRequest(h.mocks, params)

	if err != nil {
		respondErr(w, err)
		return
	}

	if !result.Matches {
		noMatch(w, result)
		return
	}

	m := result.Matched
	res, err := result.Matched.Reply.Build(r, m)
	if err != nil {
		respondErr(w, err)
		return
	}

	m.Hit()

	for k := range res.Header {
		w.Header().Add(k, res.Header.Get(k))
	}

	w.WriteHeader(res.Status)

	if res.Body == nil {
		return
	}

	w.Write(res.Body)
}

func noMatch(w http.ResponseWriter, result *findMockResult) {
	w.WriteHeader(http.StatusTeapot)
	_, _ = w.Write([]byte("Request was not matched."))

	if result.ClosestMatch != nil {
		_, _ = w.Write([]byte("\n"))
		_, _ = w.Write([]byte("\n"))
	}
}

func respondErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusTeapot)
	_, _ = w.Write([]byte("Request was not matched."))
	_, _ = w.Write([]byte(err.Error()))

	log.Printf("Reason: %v", err)
}

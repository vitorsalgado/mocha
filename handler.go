package mocha

import (
	"bufio"
	"log"
	"net/http"

	"github.com/vitorsalgado/mocha/matcher"
	"github.com/vitorsalgado/mocha/mock"
	"github.com/vitorsalgado/mocha/params"
)

type Handler struct {
	mocks   mock.Storage
	parsers []BodyParser
	params  *params.Params
}

func newHandler(
	mockstore mock.Storage,
	parsers []BodyParser,
	params *params.Params,
) *Handler {
	return &Handler{mocks: mockstore, parsers: parsers, params: params}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parsedbody, err := ParseRequestBody(r, h.parsers)
	if err != nil {
		respondErr(w, err)
		return
	}

	params := matcher.Params{RequestInfo: &matcher.RequestInfo{Request: r, ParsedBody: parsedbody}, Extras: h.params}
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
	res, err := result.Matched.Reply.Build(r, m, h.params)
	if err != nil {
		respondErr(w, err)
		return
	}

	m.Hit()

	for k := range res.Header {
		w.Header().Add(k, res.Header.Get(k))
	}

	w.WriteHeader(res.Status)

	if res.Body != nil {
		scanner := bufio.NewScanner(res.Body)
		for scanner.Scan() {
			w.Write(scanner.Bytes())
		}

		if scanner.Err() != nil {
			panic(scanner.Err())
		}
	}

	args := mock.PostActionArgs{Request: r, Response: res, Mock: m, Params: h.params}
	for _, action := range m.PostActions {
		err := action.Run(args)
		if err != nil {
			log.Println(err)
		}
	}
}

func noMatch(w http.ResponseWriter, result *findMockResult) {
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte("Request was not matched."))

	if result.ClosestMatch != nil {
		w.Write([]byte("\n"))
		w.Write([]byte("\n"))
	}
}

func respondErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte("Request was not matched."))
	w.Write([]byte(err.Error()))

	log.Printf("Reason: %v", err)
}

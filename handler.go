package mocha

import (
	"bufio"
	"log"
	"net/http"
	"time"

	"github.com/vitorsalgado/mocha/internal/params"
	"github.com/vitorsalgado/mocha/mock"
	"github.com/vitorsalgado/mocha/to"
)

type mockHandler struct {
	mocks       mock.Storage
	bodyParsers []RequestBodyParser
	params      params.Params
	t           mock.T
}

func newHandler(
	storage mock.Storage,
	bodyParsers []RequestBodyParser,
	params params.Params,
	t mock.T,
) *mockHandler {
	return &mockHandler{mocks: storage, bodyParsers: bodyParsers, params: params, t: t}
}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.t.Helper()

	parser, err := parseRequestBody(r, h.bodyParsers)
	if err != nil {
		respondError(w, err)
		return
	}

	// Match current request with all eligible stored matchers in order to find one mock.
	parameters := to.Args{RequestInfo: &to.RequestInfo{Request: r, ParsedBody: parser}, Params: h.params}
	result, err := mock.FindForRequest(h.mocks, parameters, h.t)
	if err != nil {
		respondError(w, err)
		return
	}

	if !result.Matches {
		respondNonMatched(w, result)
		return
	}

	m := result.Matched
	m.Hit()

	// Run post matchers, after standard ones and after marking the mock as called.
	afterResult, err := m.Matches(parameters, m.PostExpectations, h.t)
	if err != nil {
		respondError(w, err)
		return
	}

	if !afterResult.IsMatch {
		respondNonMatched(w, result)
		return
	}

	// Get the reply for the mock, after running all possible matchers.
	res, err := result.Matched.Reply.Build(r, m, h.params)
	if err != nil {
		h.t.Errorf(err.Error())
		respondError(w, err)
		return
	}

	// Map the response using mock mappers.
	mapperArgs := mock.ResponseMapperArgs{Request: r, Parameters: h.params}
	for _, mapper := range res.Mappers {
		if err = mapper(res, mapperArgs); err != nil {
			respondError(w, err)
			return
		}
	}

	// If a delay is set, it will wait before continuing serving the mocked response.
	if res.Delay > 0 {
		<-time.After(res.Delay)
	}

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

	// Run post actions.
	// Errors that happens here will only be logged.
	args := mock.PostActionArgs{Request: r, Response: res, Mock: m, Params: h.params}
	for _, action := range m.PostActions {
		err = action.Run(args)
		if err != nil {
			log.Println(err)
		}
	}
}

func respondNonMatched(w http.ResponseWriter, result *mock.FindResult) {
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte("Request was not matched."))

	if result.ClosestMatch != nil {
		w.Write([]byte("\n"))
		w.Write([]byte("\n"))
	}
}

func respondError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte("Request was not matched."))
	w.Write([]byte(err.Error()))

	log.Printf("Reason: %v", err)
}

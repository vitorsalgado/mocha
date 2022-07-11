package mocha

import (
	"bufio"
	"fmt"
	"net/http"
	"time"

	"github.com/vitorsalgado/mocha/core"
	"github.com/vitorsalgado/mocha/expect"
	"github.com/vitorsalgado/mocha/internal/parameters"
	"github.com/vitorsalgado/mocha/x/headers"
	"github.com/vitorsalgado/mocha/x/mimetypes"
)

type mockHandler struct {
	mocks       core.Storage
	bodyParsers []RequestBodyParser
	params      parameters.Params
	t           core.T
}

func newHandler(
	storage core.Storage,
	bodyParsers []RequestBodyParser,
	params parameters.Params,
	t core.T,
) *mockHandler {
	return &mockHandler{mocks: storage, bodyParsers: bodyParsers, params: params, t: t}
}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.t.Helper()

	parsedBody, err := parseRequestBody(r, h.bodyParsers)
	if err != nil {
		respondError(w, err)
		return
	}

	// match current request with all eligible stored matchers in order to find one mock.
	args := expect.Args{
		RequestInfo: &expect.RequestInfo{Request: r, ParsedBody: parsedBody},
		Params:      h.params}
	result, err := core.FindMockForRequest(h.mocks, args, h.t)
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

	// run post matchers, after standard ones and after marking the mock as called.
	afterResult, err := m.Matches(args, m.PostExpectations, h.t)
	if err != nil {
		respondError(w, err)
		return
	}

	if !afterResult.IsMatch {
		respondNonMatched(w, result)
		return
	}

	// get the reply for the mock, after running all possible matchers.
	res, err := result.Matched.Reply.Build(r, m, h.params)
	if err != nil {
		h.t.Logf(err.Error())
		respondError(w, err)
		return
	}

	// map the response using mock mappers.
	mapperArgs := core.ResponseMapperArgs{Request: r, Parameters: h.params}
	for _, mapper := range res.Mappers {
		if err = mapper(res, mapperArgs); err != nil {
			respondError(w, err)
			return
		}
	}

	// if a delay is set, it will wait before continuing serving the mocked response.
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
			h.t.Logf("error writing response body: error=%v", scanner.Err())
		}
	}

	// run post actions.
	paArgs := core.PostActionArgs{Request: r, Response: res, Mock: m, Params: h.params}
	for i, action := range m.PostActions {
		err = action.Run(paArgs)
		if err != nil {
			h.t.Logf("\nan error occurred running post action %d. error=%v", i, err)
		}
	}
}

func respondNonMatched(w http.ResponseWriter, result *core.FindResult) {
	w.Header().Add(headers.ContentType, mimetypes.TextPlain)
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte("Request did not match.\n"))

	if result.ClosestMatch != nil {
		w.Write([]byte("Closest Match:\n"))
		w.Write([]byte(fmt.Sprintf("ID: %d\n", result.ClosestMatch.ID)))
		w.Write([]byte(fmt.Sprintf("Name: %s", result.ClosestMatch.Name)))
	}
}

func respondError(w http.ResponseWriter, err error) {
	w.Header().Add(headers.ContentType, mimetypes.TextPlain)
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte("Request did not match. An error occurred.\n"))
	w.Write([]byte(err.Error()))
}

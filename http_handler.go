package mocha

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/vitorsalgado/mocha/v3/expect"
	"github.com/vitorsalgado/mocha/v3/hooks"
	"github.com/vitorsalgado/mocha/v3/internal/headers"
	"github.com/vitorsalgado/mocha/v3/internal/mimetypes"
	"github.com/vitorsalgado/mocha/v3/params"
	"github.com/vitorsalgado/mocha/v3/reply"
)

type mockHandler struct {
	mocks       storage
	scenarios   scenarioStore
	bodyParsers []RequestBodyParser
	params      params.P
	evt         *hooks.Emitter
	t           T
}

func newHandler(
	storage storage,
	scenarios scenarioStore,
	bodyParsers []RequestBodyParser,
	params params.P,
	evt *hooks.Emitter,
	t T,
) *mockHandler {
	return &mockHandler{mocks: storage, scenarios: scenarios, bodyParsers: bodyParsers, params: params, evt: evt, t: t}
}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	er := hooks.FromRequest(r)

	h.evt.Emit(hooks.OnRequest{Request: er, StartedAt: start})

	parsedBody, err := parseRequestBody(r, h.bodyParsers)
	if err != nil {
		respondError(w, r, h.evt, err)
		return
	}

	// match current request with all eligible stored matchers in order to find one mock.
	args := expect.Args{
		RequestInfo: &expect.RequestInfo{Request: r, ParsedBody: parsedBody},
		Params:      h.params}
	result, err := findMockForRequest(h.mocks, args)
	if err != nil {
		respondError(w, r, h.evt, err)
		return
	}

	if !result.Matches {
		respondNonMatched(w, r, result, h.evt)
		return
	}

	mock := result.Matched

	if mock.Repeat > 0 && mock.Hits()+1 > mock.Repeat {
		respondError(w, r, h.evt,
			fmt.Errorf("mock is set to respond only %d times. current hits is %d", mock.Repeat, mock.Hits()))
		return
	}

	if mock.ScenarioName != "" {
		scn, ok := h.scenarios.FetchByName(mock.ScenarioName)

		if !ok && mock.ScenarioRequiredState == _scenarioStateStarted {
			scn = h.scenarios.CreateNewIfNeeded(mock.ScenarioName)
			ok = true
		}

		if ok {
			if scn.State == mock.ScenarioRequiredState {
				if mock.ScenarioNewState != "" {
					scn.State = mock.ScenarioNewState
					h.scenarios.Save(scn)
				}
			} else {
				respondError(w, r, h.evt,
					fmt.Errorf("expected mock id=%d scenario=%s to be %s. got %s",
						mock.ID, mock.ScenarioName, mock.ScenarioRequiredState, scn.State))
				return
			}
		}
	}

	// get the reply for the mock, after running all possible matchers.
	res, err := result.Matched.Reply.Build(r, mock, h.params)
	if err != nil {
		h.t.Logf(err.Error())
		respondError(w, r, h.evt, err)
		return
	}

	// map the response using mock mappers.
	mapperArgs := reply.ResponseMapperArgs{Request: r, Parameters: h.params}
	for _, mapper := range res.Mappers {
		if err = mapper(res, mapperArgs); err != nil {
			respondError(w, r, h.evt, err)
			return
		}
	}

	// success
	mock.Hit()

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
	paArgs := PostActionArgs{Request: r, Response: res, Mock: mock, Params: h.params}
	for i, action := range mock.PostActions {
		err = action.Run(paArgs)
		if err != nil {
			h.t.Logf("\nan error occurred running post action %d. error=%v", i, err)
		}
	}

	h.evt.Emit(hooks.OnRequestMatch{
		Request:            er,
		ResponseDefinition: hooks.Response{Status: res.Status, Header: res.Header.Clone()},
		Mock:               hooks.Mock{ID: mock.ID, Name: mock.Name},
		Elapsed:            time.Since(start)})
}

func respondNonMatched(w http.ResponseWriter, r *http.Request, result *findResult, evt *hooks.Emitter) {
	e := hooks.OnRequestNotMatched{Request: hooks.FromRequest(r), Result: hooks.Result{Details: make([]hooks.ResultDetail, 0)}}

	if result.ClosestMatch != nil {
		e.Result.HasClosestMatch = true
		e.Result.ClosestMatch = hooks.Mock{ID: result.ClosestMatch.ID, Name: result.ClosestMatch.Name}
	}

	for _, detail := range result.MismatchDetails {
		e.Result.Details = append(e.Result.Details,
			hooks.ResultDetail{Name: detail.Name, Description: detail.Description, Target: detail.Target})
	}

	evt.Emit(e)

	builder := strings.Builder{}
	builder.WriteString("REQUEST DID NOT MATCH.\n")

	if result.ClosestMatch != nil {
		builder.WriteString(
			fmt.Sprintf("Closest Match: %d %s\n\n", result.ClosestMatch.ID, result.ClosestMatch.Name))
	}

	builder.WriteString("Mismatches:\n")

	for _, detail := range result.MismatchDetails {
		builder.WriteString(fmt.Sprintf("%s, reason=%s, applied-to=%s\n",
			detail.Name, detail.Description, detail.Target))
	}

	w.Header().Add(headers.ContentType, mimetypes.TextPlain)
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte(builder.String()))
}

func respondError(w http.ResponseWriter, r *http.Request, evt *hooks.Emitter, err error) {
	evt.Emit(hooks.OnError{Request: hooks.FromRequest(r), Err: err})

	w.Header().Add(headers.ContentType, mimetypes.TextPlain)
	w.WriteHeader(http.StatusTeapot)

	w.Write([]byte("Request did not match. An error occurred.\n"))
	w.Write([]byte(fmt.Sprintf("%v", err)))
}

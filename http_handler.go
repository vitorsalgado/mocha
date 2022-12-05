package mocha

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/internal/mimetype"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/mod"
	"github.com/vitorsalgado/mocha/v3/reply"
)

type mockHandler struct {
	mocks       mockStore
	bodyParsers []RequestBodyParser
	params      reply.Params
	evt         *eventListener
	t           TestingT
}

func newHandler(
	storage mockStore,
	bodyParsers []RequestBodyParser,
	params reply.Params,
	evt *eventListener,
	t TestingT,
) *mockHandler {
	return &mockHandler{mocks: storage, bodyParsers: bodyParsers, params: params, evt: evt, t: t}
}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	evtReq := evtRequest(r)

	parsedBody, err := parseRequestBody(r, h.bodyParsers)
	if err != nil {
		h.evt.Emit(&OnRequest{Request: evtReq, StartedAt: start})
		respondError(w, evtReq, h.evt, fmt.Errorf("error parsing request body. reason=%w", err))
		return
	}

	evtReq.Body = parsedBody
	h.evt.Emit(&OnRequest{Request: evtReq, StartedAt: start})

	// match current request with all eligible stored matchers in order to find one mock.
	info := &matcher.RequestInfo{Request: r, ParsedBody: parsedBody}
	result, err := findMockForRequest(h.mocks, info)
	if err != nil {
		respondError(w, evtReq, h.evt, fmt.Errorf("error trying to find a mock. reason=%w", err))
		return
	}

	if !result.Matches {
		respondNonMatched(w, evtReq, result, h.evt)
		return
	}

	mock := result.Matched

	if mock.Delay > 0 {
		<-time.After(mock.Delay)
	}

	// request with reply vars
	r = r.WithContext(
		context.WithValue(
			r.Context(), reply.KArg, &reply.Arg{M: reply.M{Hits: mock.Hits()}, Params: h.params}))

	// get the reply for the mock, after running all possible matchers.
	res, err := result.Matched.Reply.Build(w, r)
	if err != nil {
		h.t.Logf(err.Error())
		respondError(w, evtReq, h.evt, fmt.Errorf("error building reply. reason=%w", err))
		return
	}

	if res.SendPending() {
		// map the response using mock mappers.
		mapperArgs := &reply.MapperArgs{Request: r, Parameters: h.params}
		for i, mapper := range res.Mappers {
			if err = mapper(res, mapperArgs); err != nil {
				respondError(w, evtReq, h.evt, fmt.Errorf("error with response mapper[%d]. reason=%w", i, err))
				return
			}
		}
	}

	// success
	mock.Inc()

	if res.SendPending() {
		for k, v := range res.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}

		w.WriteHeader(res.Status)

		if res.Body != nil {
			w.Write(res.Body)
		}
	}

	for _, exp := range mock.expectations {
		err = exp.Matcher.OnMockServed()
		if err != nil {
			h.t.Logf("matcher %s .OnMockServed() returned the error=%v", exp.Matcher.Name(), err)
		}
	}

	// run post actions.
	paArgs := &PostActionArgs{Request: r, Response: res, Mock: mock, Params: h.params}
	for i, action := range mock.PostActions {
		err = action.Run(paArgs)
		if err != nil {
			h.t.Logf("\nerror running post action %d. error=%v", i, err)
		}
	}

	h.evt.Emit(&OnRequestMatch{
		Request:            evtReq,
		ResponseDefinition: mod.EvtRes{Status: res.Status, Header: res.Header.Clone()},
		Mock:               mod.EvtMk{ID: mock.ID, Name: mock.Name},
		Elapsed:            time.Since(start)})
}

func respondNonMatched(w http.ResponseWriter, r mod.EvtReq, result *findResult, evt *eventListener) {
	e := &OnRequestNotMatched{Request: r, Result: mod.EvtResult{Details: make([]mod.EvtResultExt, 0)}}

	if result.ClosestMatch != nil {
		e.Result.HasClosestMatch = true
		e.Result.ClosestMatch = mod.EvtMk{ID: result.ClosestMatch.ID, Name: result.ClosestMatch.Name}
	}

	for _, detail := range result.MismatchDetails {
		e.Result.Details = append(e.Result.Details,
			mod.EvtResultExt{Name: detail.Name, Description: detail.Desc, Target: detail.Target})
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
			detail.Name, detail.Desc, detail.Target))
	}

	w.Header().Add(header.ContentType, mimetype.TextPlain)
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte(builder.String()))
}

func respondError(w http.ResponseWriter, r mod.EvtReq, evt *eventListener, err error) {
	evt.Emit(&OnError{Request: r, Err: err})

	w.Header().Add(header.ContentType, mimetype.TextPlain)
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte(fmt.Sprintf("EvtReq did not match. An error occurred.\n%v", err)))
}

func evtRequest(r *http.Request) mod.EvtReq {
	return mod.EvtReq{
		Method:     r.Method,
		Path:       r.URL.Path,
		RequestURI: r.RequestURI,
		Host:       r.Host,
		Header:     r.Header.Clone(),
	}
}

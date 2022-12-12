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
	conf          *Config
	store         mockStore
	bodyParsers   []RequestBodyParser
	params        reply.Params
	proxy         *proxy
	eventListener *eventListener
	rec           *record
	t             TestingT
	d             Debug
}

func newHandler(
	conf *Config,
	store mockStore,
	bodyParsers []RequestBodyParser,
	params reply.Params,
	proxy *proxy,
	evt *eventListener,
	rec *record,
	t TestingT,
	d Debug,
) *mockHandler {
	return &mockHandler{
		conf:          conf,
		store:         store,
		bodyParsers:   bodyParsers,
		params:        params,
		proxy:         proxy,
		eventListener: evt,
		rec:           rec,
		t:             t,
		d:             d,
	}
}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	evtReq := evtRequest(r)

	parsedBody, bb, err := parseRequestBody(r, h.bodyParsers)
	if err != nil {
		h.eventListener.Emit(&OnRequest{Request: evtReq, StartedAt: start})
		h.respondError(w, evtReq, fmt.Errorf("error parsing request body. reason=%w", err))
		return
	}

	evtReq.Body = bb
	h.eventListener.Emit(&OnRequest{Request: evtReq, StartedAt: start})

	// match current request with all eligible stored matchers in order to find one mock.
	info := &matcher.RequestInfo{Request: r, ParsedBody: parsedBody}
	result := findMockForRequest(h.store, info)

	if !result.Matches {
		if h.proxy != nil {
			// proxy non-matched requests.
			h.proxy.Proxy(w, r, bb)
			return
		}

		h.respondNonMatched(w, evtReq, result)
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
		h.respondError(w, evtReq, fmt.Errorf("error building reply. reason=%w", err))
		return
	}

	mock.Inc()

	if res.SendPending() {
		// map the response using mock mappers.
		mapperArgs := &MapperArgs{Request: r, Parameters: h.params}
		for i, mapper := range mock.Mappers {
			if err = mapper(res, mapperArgs); err != nil {
				mock.Dec()
				h.respondError(w, evtReq, fmt.Errorf("error with response mapper[%d]. reason=%w", i, err))
				return
			}
		}

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
			h.maybeDebug(err)
		}
	}

	// run post actions.
	paArgs := &PostActionArgs{Request: r, Response: res, Mock: mock, Params: h.params}
	for i, action := range mock.PostActions {
		err = action.Run(paArgs)
		if err != nil {
			h.t.Logf("\nerror running post action [%d]. error=%v", i, err)
			h.maybeDebug(err)
		}
	}

	h.eventListener.Emit(&OnRequestMatch{
		Request:            evtReq,
		ResponseDefinition: mod.EvtRes{Status: res.Status, Header: res.Header.Clone()},
		Mock:               mod.EvtMk{ID: mock.ID, Name: mock.Name},
		Elapsed:            time.Since(start)})

	if h.conf.Record != nil {
		doRec(h.rec, r, bb, res)
	}
}

func (h *mockHandler) respondNonMatched(w http.ResponseWriter, r *mod.EvtReq, result *findResult) {
	e := &OnRequestNotMatched{Request: r, Result: mod.EvtResult{Details: make([]mod.EvtResultExt, 0)}}

	if result.ClosestMatch != nil {
		e.Result.HasClosestMatch = true
		e.Result.ClosestMatch = mod.EvtMk{ID: result.ClosestMatch.ID, Name: result.ClosestMatch.Name}
	}

	for _, detail := range result.MismatchDetails {
		e.Result.Details = append(e.Result.Details,
			mod.EvtResultExt{Name: detail.Name, Description: detail.Desc, Target: detail.Target})
	}

	h.eventListener.Emit(e)

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
	w.WriteHeader(StatusNoMockFound)
	w.Write([]byte(builder.String()))
}

func (h *mockHandler) respondError(w http.ResponseWriter, r *mod.EvtReq, err error) {
	h.eventListener.Emit(&OnError{Request: r, Err: err})
	h.maybeDebug(err)

	w.Header().Add(header.ContentType, mimetype.TextPlain)
	w.WriteHeader(StatusNoMockFound)
	w.Write([]byte(fmt.Sprintf("EvtReq did not match. An error occurred.\n%v", err)))
}

func (h *mockHandler) maybeDebug(err error) {
	if h.d != nil {
		h.d(err)
	}
}

func doRec(rec *record, r *http.Request, bb []byte, res *reply.Response) {
	rec.record(&recArgs{
		request: recRequest{
			path:   r.URL.Path,
			method: r.Method,
			header: r.Header.Clone(),
			query:  r.URL.Query(),
			body:   bb,
		},
		response: recResponse{
			status: res.Status,
			header: res.Header.Clone(),
			body:   res.Body,
		},
	})
}

func evtRequest(r *http.Request) *mod.EvtReq {
	return &mod.EvtReq{
		Method:     r.Method,
		Path:       r.URL.Path,
		RequestURI: r.RequestURI,
		Host:       r.Host,
		Header:     r.Header.Clone(),
	}
}

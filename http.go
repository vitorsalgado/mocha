package mocha

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/vitorsalgado/mocha/v3/event"
	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/internal/httpx"
	"github.com/vitorsalgado/mocha/v3/internal/mimetype"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/reply"
)

type mockHandler struct {
	app *Mocha
}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	evtReq := h.toEvent(r)

	w = httpx.DecorateWriter(w)

	parsedBody, rawBody, err := parseRequestBody(r, h.app.requestBodyParsers)
	if err != nil {
		h.app.listener.Emit(&event.OnRequest{Request: evtReq, StartedAt: start})
		h.onError(w, evtReq, fmt.Errorf("error parsing request body. reason=%w", err))
		return
	}

	evtReq.Body = rawBody
	h.app.listener.Emit(&event.OnRequest{Request: evtReq, StartedAt: start})

	// match current request with all eligible stored matchers in order to find one mock.
	info := &matcher.RequestInfo{Request: r, ParsedBody: parsedBody}
	result := findMockForRequest(h.app.storage, info)

	if !result.Matches {
		if h.app.proxy != nil {
			// proxy non-matched requests.
			h.app.proxy.ServeHTTP(w, r)
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
			r.Context(), reply.KArg, &reply.Arg{MockInfo: reply.MockInfo{Hits: mock.Hits()}, Params: h.app.params}))

	// get the reply for the mock, after running all possible matchers.
	res, err := result.Matched.Reply.Build(w, r)
	if err != nil {
		h.app.t.Logf(err.Error())
		h.onError(w, evtReq, fmt.Errorf("error building reply. reason=%w", err))
		return
	}

	mock.Inc()

	if res.Sent() {
		// map the response using mock mappers.
		mapperArgs := &MapperArgs{Request: r, Parameters: h.app.params}
		for i, mapper := range mock.Mappers {
			if err = mapper(res, mapperArgs); err != nil {
				mock.Dec()
				h.onError(w, evtReq, fmt.Errorf("error with response mapper[%d]. reason=%w", i, err))
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
			h.app.t.Logf("matcher %s .OnMockServed() returned the error=%v", exp.Matcher.Name(), err)
		}
	}

	// run post actions.
	paArgs := &PostActionArgs{Request: r, Response: res, Mock: mock, Params: h.app.params}
	for i, action := range mock.PostActions {
		err = action.Run(paArgs)
		if err != nil {
			h.app.t.Logf("\nerror running post action [%d]. error=%v", i, err)
		}
	}

	h.app.listener.Emit(&event.OnRequestMatch{
		Request:            evtReq,
		ResponseDefinition: event.EvtRes{Status: res.Status, Header: res.Header.Clone()},
		Mock:               event.EvtMk{ID: mock.ID, Name: mock.Name},
		Elapsed:            time.Since(start)})

	if h.app.rec != nil {
		rw := w.(*httpx.Rw)
		recorded := rw.Result()

		defer recorded.Body.Close()

		body, err := io.ReadAll(recorded.Body)
		if err != nil {
			h.app.t.Errorf(fmt.Sprintf("error reading recorded body. reason=%s", err.Error()))
			h.app.listener.Emit(&event.OnError{Request: evtReq, Err: err})
			return
		}

		h.app.rec.record(r, rawBody, recorded, body)
	}
}

func (h *mockHandler) respondNonMatched(w http.ResponseWriter, r *event.EvtReq, result *findResult) {
	e := &event.OnRequestNotMatched{Request: r, Result: event.EvtResult{Details: make([]event.EvtResultExt, 0)}}

	if result.ClosestMatch != nil {
		e.Result.HasClosestMatch = true
		e.Result.ClosestMatch = event.EvtMk{ID: result.ClosestMatch.ID, Name: result.ClosestMatch.Name}
	}

	for _, detail := range result.MismatchDetails {
		e.Result.Details = append(e.Result.Details, event.EvtResultExt{
			Name:        detail.Name,
			Description: detail.Desc,
			Target:      strconv.FormatInt(int64(detail.Target), 10),
		})
	}

	h.app.listener.Emit(e)

	builder := strings.Builder{}
	builder.WriteString("REQUEST DID NOT MATCH.\n")

	if result.ClosestMatch != nil {
		builder.WriteString(
			fmt.Sprintf("Closest Match: %s %s\n\n", result.ClosestMatch.ID, result.ClosestMatch.Name))
	}

	builder.WriteString("Mismatches:\n")

	for _, detail := range result.MismatchDetails {
		builder.WriteString(fmt.Sprintf("%s, reason=%s, applied-to=%v\n",
			detail.Name, detail.Desc, detail.Target))
	}

	w.Header().Add(header.ContentType, mimetype.TextPlain)
	w.WriteHeader(StatusNoMockFound)
	w.Write([]byte(builder.String()))
}

func (h *mockHandler) onError(w http.ResponseWriter, r *event.EvtReq, err error) {
	h.app.t.Logf(err.Error())
	h.app.listener.Emit(&event.OnError{Request: r, Err: err})

	w.Header().Add(header.ContentType, mimetype.TextPlain)
	w.WriteHeader(StatusNoMockFound)
	w.Write([]byte(fmt.Sprintf("EvtReq did not match. An error occurred.\n%v", err)))
}

func (h *mockHandler) toEvent(r *http.Request) *event.EvtReq {
	return &event.EvtReq{
		Method:     r.Method,
		Path:       r.URL.Path,
		RequestURI: r.RequestURI,
		Host:       r.Host,
		Header:     r.Header.Clone(),
	}
}

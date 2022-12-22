package mocha

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/internal/httpx"
	"github.com/vitorsalgado/mocha/v3/internal/mimetype"
	"github.com/vitorsalgado/mocha/v3/reply"
	"github.com/vitorsalgado/mocha/v3/types"
	"github.com/vitorsalgado/mocha/v3/x/event"
)

type mockHandler struct {
	app *Mocha
}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	evtReq := event.FromRequest(r)

	w = httpx.Wrap(w)

	parsedBody, rawBody, err := parseRequestBody(r, h.app.requestBodyParsers)
	if err != nil {
		h.app.listener.Emit(&event.OnRequest{Request: evtReq, StartedAt: start})
		h.onError(w, evtReq, fmt.Errorf("error parsing request body. reason=%w", err))
		return
	}

	evtReq.Body = rawBody
	h.app.listener.Emit(&event.OnRequest{Request: evtReq, StartedAt: start})

	// match current request with all eligible stored matchers in order to find one mock.
	info := &values{Request: r, ParsedBody: parsedBody}
	result := findMockForRequest(h.app.storage, info)

	if !result.Pass {
		if h.app.proxy != nil {
			// proxy non-matched requests.
			h.app.proxy.ServeHTTP(w, r)
			res, err := h.buildResponseFromWriter(w)
			if err != nil {
				h.app.log.Logf(err.Error())
				return
			}

			if h.app.rec != nil {
				h.app.rec.record(r, rawBody, res)
			}
			return
		} else {
			h.onNoMatches(w, evtReq, result)
			return
		}
	}

	mock := result.Matched

	if mock.Delay > 0 {
		<-time.After(mock.Delay)
	}

	rawURL := r.URL.Path
	if len(r.URL.RawQuery) > 0 {
		rawURL += "?" + r.URL.RawQuery
	}

	u, _ := url.Parse(rawURL)
	reqValues := &types.RequestValues{RawRequest: r, URL: u, Body: parsedBody}

	res, err := result.Matched.Reply.Build(w, reqValues)
	if err != nil {
		h.app.log.Logf(err.Error())
		h.onError(w, evtReq, fmt.Errorf("error building reply. reason=%w", err))
		return
	}

	mock.Inc()

	if res != nil {
		// map the response using mock mappers.
		mapperArgs := &MapperIn{Request: r, Parameters: h.app.params}
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

		for _, cookie := range res.Cookies {
			http.SetCookie(w, cookie)
		}

		w.WriteHeader(res.StatusCode)

		if res.Body != nil {
			w.Write(res.Body)
		}
	} else {
		res, err = h.buildResponseFromWriter(w)
		if err != nil {
			h.app.log.Logf(err.Error())
			return
		}
	}

	for _, exp := range mock.expectations {
		err = exp.Matcher.OnMockServed()
		if err != nil {
			h.app.log.Logf("matcher %s .OnMockServed() returned the error=%v", exp.Matcher.Name(), err)
		}
	}

	input := &PostActionIn{Request: r, Response: res, Params: h.app.params}
	for i, action := range mock.PostActions {
		err = action.Run(input)
		if err != nil {
			h.app.log.Logf("\nerror running post action [%d]. error=%v", i, err)
		}
	}

	h.app.listener.Emit(&event.OnRequestMatch{
		Request:            evtReq,
		ResponseDefinition: event.EvtRes{Status: res.StatusCode, Header: res.Header},
		Mock:               event.EvtMk{ID: mock.ID, Name: mock.Name},
		Elapsed:            time.Since(start)})

	if h.app.rec != nil && h.app.config.Proxy == nil {
		h.app.rec.record(r, rawBody, res)
	}
}

func (h *mockHandler) onNoMatches(w http.ResponseWriter, r *event.EvtReq, result *findResult) {
	e := &event.OnRequestNotMatched{Request: r, Result: event.EvtResult{Details: make([]event.EvtResultExt, 0)}}

	if result.ClosestMatch != nil {
		e.Result.HasClosestMatch = true
		e.Result.ClosestMatch = event.EvtMk{ID: result.ClosestMatch.ID, Name: result.ClosestMatch.Name}
	}

	for _, detail := range result.MismatchDetails {
		e.Result.Details = append(e.Result.Details, event.EvtResultExt{
			Name:        detail.MatchersName,
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
			detail.MatchersName, detail.Desc, detail.Target))
	}

	w.Header().Add(header.ContentType, mimetype.TextPlain)
	w.WriteHeader(StatusRequestDidNotMatch)
	w.Write([]byte(builder.String()))
}

func (h *mockHandler) onError(w http.ResponseWriter, r *event.EvtReq, err error) {
	h.app.log.Logf(err.Error())
	h.app.listener.Emit(&event.OnError{Request: r, Err: err})

	w.Header().Add(header.ContentType, mimetype.TextPlain)
	w.WriteHeader(StatusRequestDidNotMatch)
	w.Write([]byte(fmt.Sprintf("An error occurred.\n%s", err.Error())))
}

func (h *mockHandler) buildResponseFromWriter(w http.ResponseWriter) (*reply.Stub, error) {
	rw := w.(*httpx.Rw)
	result := rw.Result()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	defer result.Body.Close()

	return &reply.Stub{
		StatusCode: result.StatusCode,
		Header:     result.Header,
		Cookies:    result.Cookies(),
		Body:       body,
	}, nil
}

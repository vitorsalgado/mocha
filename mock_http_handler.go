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
	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/x/event"
)

type mockHandler struct {
	app *Mocha
}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	evtReq := event.FromRequest(r)

	reqPath := r.URL.Path
	if len(r.URL.RawQuery) > 0 {
		reqPath += "?" + r.URL.RawQuery
	}

	urlSegments := strings.Split(reqPath, "/")
	parsedURL, _ := url.Parse(h.app.URL() + reqPath)

	if h.app.config.Record != nil {
		w = httpx.Wrap(w)
	}

	parsedBody, rawBody, err := parseRequestBody(r, h.app.requestBodyParsers)

	evtReq.URL = parsedURL.String()
	evtReq.Body = rawBody
	h.app.listener.Emit(&event.OnRequest{Request: evtReq, StartedAt: start})

	if err != nil {
		h.app.listener.Emit(&event.OnError{Request: evtReq, Err: fmt.Errorf("error parsing request body. reason=%w", err)})
	}

	result := findMockForRequest(h.app.storage, &valueSelectorInput{r, parsedURL, parsedBody})

	if !result.Pass {
		if h.app.proxy != nil {
			// proxy non-matched requests.
			h.app.proxy.ServeHTTP(w, r)

			if h.app.rec != nil {
				res, err := h.buildStubFromWriter(w)
				if err != nil {
					h.app.log.Logf(err.Error())
					return
				}

				h.app.rec.dispatch(r, parsedURL, rawBody, res)
			}
			return
		} else {
			h.onNoMatches(w, evtReq, result)
			return
		}
	}

	mock := result.Matched

	if mock.Delay > 0 {
		time.Sleep(mock.Delay)
	}

	reqValues := &RequestValues{r, parsedURL, urlSegments, parsedBody, h.app, mock}
	res, err := result.Matched.Reply.Build(w, reqValues)
	if err != nil {
		h.app.log.Logf(err.Error())
		h.onError(w, evtReq, fmt.Errorf("error building reply. reason=%w", err))
		return
	}

	mock.Inc()

	if res != nil {
		// map the response using mock mappers.
		for i, mapper := range mock.Mappers {
			if err = mapper(reqValues, res); err != nil {
				mock.Dec()
				h.onError(w, evtReq, fmt.Errorf("error with response mapper at index [%d]. reason=%w", i, err))
				return
			}
		}

		for k, v := range res.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}

		for k := range res.Trailer {
			w.Header().Add(header.Trailer, k)
		}

		for _, cookie := range res.Cookies {
			http.SetCookie(w, cookie)
		}

		w.WriteHeader(res.StatusCode)

		if res.Body != nil {
			w.Write(res.Body)
		}

		for k, v := range res.Trailer {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
	} else {
		res, err = h.buildStubFromWriter(w)
		if err != nil {
			h.app.log.Logf(err.Error())
			return
		}
	}

	for _, aa := range mock.after {
		err = aa.AfterMockServed()
		if err != nil {
			h.app.log.Logf("matcher %s .AfterMockServed() returned the error=%v", aa.(matcher.Matcher).Name(), err)
		}
	}

	input := &PostActionInput{r, parsedURL, parsedBody, h.app, mock, res, nil}
	for i, action := range mock.PostActions {
		err = action.Run(input)
		if err != nil {
			h.app.log.Logf("\nerror running post action [%d]. error=%v", i, err)
		}
	}

	h.app.listener.Emit(&event.OnRequestMatch{
		Request: evtReq,
		ResponseDefinition: event.EvtRes{
			Status:  res.StatusCode,
			Header:  res.Header,
			Body:    res.Body,
			Encoded: res.Header.Get(HeaderContentEncoding) != "",
		},
		Mock:    event.EvtMk{ID: mock.ID, Name: mock.Name},
		Elapsed: time.Since(start)})

	if h.app.rec != nil && h.app.config.Proxy == nil {
		h.app.rec.dispatch(r, parsedURL, rawBody, res)
	}
}

func (h *mockHandler) onNoMatches(w http.ResponseWriter, r *event.EvtReq, result *findResult) {
	e := &event.OnRequestNotMatched{Request: r, Result: event.EvtResult{Details: make([]event.EvtResultExt, 0)}}

	if result.ClosestMatch != nil {
		e.Result.HasClosestMatch = true
		e.Result.ClosestMatch = event.EvtMk{ID: result.ClosestMatch.ID, Name: result.ClosestMatch.Name}
	}

	for _, detail := range result.MismatchDetails {
		if detail.Err != nil {
			e.Result.Details = append(e.Result.Details, event.EvtResultExt{
				Name:    detail.MatchersName,
				Message: detail.Err.Error(),
				Target:  strconv.FormatInt(int64(detail.Target), 10),
			})

			continue
		}

		e.Result.Details = append(e.Result.Details, event.EvtResultExt{
			Name:    detail.MatchersName,
			Message: detail.Result.Message,
			Ext:     detail.Result.Ext,
			Target:  strconv.FormatInt(int64(detail.Target), 10),
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
		message := ""
		if detail.Err != nil {
			message = detail.Err.Error()
		} else {
			message = detail.Result.Message
		}

		builder.WriteString(fmt.Sprintf("%s, reason=%s, applied-to=%v\n",
			detail.MatchersName, message, detail.Target))
	}

	w.Header().Add(header.ContentType, mimetype.TextPlain)
	w.WriteHeader(h.app.config.MockNotFoundStatusCode)
	w.Write([]byte(builder.String()))
}

func (h *mockHandler) onError(w http.ResponseWriter, r *event.EvtReq, err error) {
	h.app.log.Logf(err.Error())
	h.app.listener.Emit(&event.OnError{Request: r, Err: err})

	w.Header().Add(header.ContentType, mimetype.TextPlain)
	w.WriteHeader(h.app.config.MockNotFoundStatusCode)
	w.Write([]byte(fmt.Sprintf("An error occurred.\n%s", err.Error())))
}

func (h *mockHandler) buildStubFromWriter(w http.ResponseWriter) (*Stub, error) {
	rw := w.(*httpx.Rw)
	result := rw.Result()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	defer result.Body.Close()

	if len(w.Header()) != len(result.Header) {
		for k, v := range w.Header() {
			if result.Header.Get(k) == "" {
				for _, vv := range v {
					result.Header.Add(k, vv)
				}
			}
		}
	}

	stub := &Stub{
		StatusCode: result.StatusCode,
		Header:     result.Header.Clone(),
		Cookies:    result.Cookies(),
	}

	if len(body) > 0 {
		stub.Body = body
	}

	return stub, nil
}

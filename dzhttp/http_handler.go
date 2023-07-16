package dzhttp

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/vitorsalgado/mocha/v3/dzhttp/httpval"
	"github.com/vitorsalgado/mocha/v3/dzhttp/internal/httprec"
	"github.com/vitorsalgado/mocha/v3/dzstd"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

type mockHandler struct {
	app       *HTTPMockApp
	lifecycle mockHTTPLifecycle
}

func newMockHandler(app *HTTPMockApp, lifecycle mockHTTPLifecycle) *mockHandler {
	return &mockHandler{app, lifecycle}
}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w = httprec.Wrap(w)
	parsedURL, urlSegments := h.parseURL(r)
	parsedBody, rawBody, err := parseRequestBody(r, h.app.requestBodyParsers)
	reqValues := &RequestValues{time.Now(), r, parsedURL, urlSegments, nil, nil, h.app, nil}
	if err != nil {
		h.lifecycle.OnWarning(reqValues, fmt.Errorf("request_body_parser: %w", err))
	}

	reqValues.RawBody = rawBody
	reqValues.ParsedBody = parsedBody

	h.lifecycle.OnRequest(reqValues)

	mocks := h.app.storage.GetEligible()
	description := dzstd.Description{Buf: make([]string, 0, len(mocks))}
	result := dzstd.FindMockForRequest(mocks,
		&HTTPValueSelectorInput{r, parsedURL, r.URL.Query(), r.Form, parsedBody}, &description)

	if !result.Pass {
		if h.app.proxy != nil {
			// proxy non-matched requests.
			h.app.proxy.ServeHTTP(w, r)

			if h.app.rec != nil {
				res := Stub{}
				err := newResponseStub(w, &res)
				if err != nil {
					h.onError(w, reqValues, err)
					return
				}

				h.app.rec.dispatch(r, parsedURL, rawBody, &res)
			}

			return
		}

		h.onNoMatches(w, reqValues, result, &description)
		return
	}

	mock := result.Matched
	reqValues.Mock = mock

	if mock.Delay != nil {
		ctxTimeout, cancel := context.WithTimeout(r.Context(), mock.Delay())
		defer cancel()

		<-ctxTimeout.Done()
	}

	stub, err := result.Matched.Reply.Build(w, reqValues)
	if err != nil {
		h.onError(w, reqValues, fmt.Errorf(
			"http: error building reply for request=%s source=%s. %w",
			reqValues.URL.Path,
			mock.Source,
			err,
		))
		return
	}

	mock.Inc()

	if stub != nil {
		for i, mapper := range mock.Mappers {
			if err = mapper(reqValues, stub); err != nil {
				mock.Dec()
				h.onError(w, reqValues, fmt.Errorf("http: error with response mapper at index %d. %w", i, err))
				return
			}
		}

		copyHeaders(stub.Header, w.Header())

		for k := range stub.Trailer {
			w.Header().Add(httpval.HeaderTrailer, k)
		}

		for _, cookie := range stub.Cookies {
			http.SetCookie(w, cookie)
		}

		w.WriteHeader(stub.StatusCode)

		if stub.Body != nil {
			if len(mock.Pipes) > 0 {
				connector := dzstd.NewConnector(mock.Pipes)
				_, err = connector.Connect(stub.Body, w)
				if err != nil {
					h.lifecycle.OnWarning(reqValues, err)
				}
			} else {
				w.Write(stub.Body)
			}
		}

		for k, v := range stub.Trailer {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
	} else {
		stub := Stub{}
		err = newResponseStub(w, &stub)
		if err != nil {
			h.onError(w, reqValues, err)
			return
		}
	}

	for _, idx := range mock.after {
		err = mock.expectations[idx].Matcher.(matcher.OnAfterMockServed).AfterMockServed()
		if err != nil {
			h.lifecycle.OnWarning(reqValues,
				fmt.Errorf("http: after mock served event: matcher[%d] %w", idx, err))
		}
	}

	callbackInput := CallbackInput{r, parsedURL, parsedBody, h.app, mock, stub}
	for i, callback := range mock.Callbacks {
		err = callback(&callbackInput)
		if err != nil {
			h.lifecycle.OnWarning(reqValues, fmt.Errorf(
				"callback: error with callback %d %T:\n%w",
				i,
				callback,
				err,
			))
		}
	}

	for i, def := range mock.PostActions {
		postAction := h.app.config.PostActions[def.Name]
		input := PostActionInput{r, parsedURL, parsedBody, h.app, mock, stub, def.RawParameters}
		err = postAction.Run(&input)
		if err != nil {
			h.lifecycle.OnWarning(reqValues, fmt.Errorf(
				"post_action: error with post action %s(%d) %T:\n%w",
				def.Name,
				i,
				postAction,
				err,
			))
		}
	}

	h.lifecycle.OnMatch(reqValues, stub)

	// when proxy is enabled, recording happens some steps before
	// so, we need to record here only if proxy is disabled.
	if h.app.rec != nil && h.app.config.Proxy == nil {
		h.app.rec.dispatch(r, parsedURL, rawBody, stub)
	}
}

func (h *mockHandler) onNoMatches(w http.ResponseWriter, r *RequestValues, result *dzstd.FindResult[*HTTPMock], desc *dzstd.Description) {
	defer h.lifecycle.OnNoMatch(r, result, desc)

	buf := bytes.Buffer{}

	buf.WriteString("REQUEST DID NOT MATCH\n")

	if result.ClosestMatch != nil {
		fmt.Fprintf(&buf, "NEAREST: %s %s\n", result.ClosestMatch.GetID(), result.ClosestMatch.GetName())
	}

	if desc.Len() > 0 {
		buf.WriteString("MISSES:\n")
		buf.WriteString(desc.String())
	}

	w.Header().Add(httpval.HeaderContentType, httpval.MIMETextPlainCharsetUTF8)
	w.WriteHeader(h.app.config.RequestWasNotMatchedStatusCode)

	buf.WriteTo(w)
}

func (h *mockHandler) onError(w http.ResponseWriter, r *RequestValues, err error) {
	defer h.lifecycle.OnError(r, err)

	w.Header().Add(httpval.HeaderContentType, httpval.MIMETextPlainCharsetUTF8)
	w.WriteHeader(h.app.config.RequestWasNotMatchedStatusCode)
	w.Write([]byte("ERROR DURING REQUEST MATCHING\n"))
	w.Write([]byte(err.Error()))
}

func (h *mockHandler) parseURL(r *http.Request) (*url.URL, []string) {
	reqPath := r.URL.Path
	if len(r.URL.RawQuery) > 0 {
		reqPath += "?" + r.URL.RawQuery
	}

	parsedURL, _ := url.Parse(h.app.URL() + reqPath)

	return parsedURL, strings.FieldsFunc(reqPath, func(c rune) bool { return c == '/' })
}

func copyHeaders(src, dst http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

package httpd

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/vitorsalgado/mocha/v3/httpd/httpval"
	"github.com/vitorsalgado/mocha/v3/httpd/internal/httprec"
	"github.com/vitorsalgado/mocha/v3/lib"
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
	result := lib.FindMockForRequest(mocks,
		&HTTPValueSelectorInput{r, parsedURL, r.URL.Query(), r.Form, parsedBody})

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

		h.onNoMatches(w, reqValues, result)
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
				connector := lib.NewConnector(mock.Pipes)
				connector.Connect(stub.Body, w)
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

	for _, aa := range mock.after {
		err = aa.AfterMockServed()
		if err != nil {
			h.lifecycle.OnWarning(reqValues,
				fmt.Errorf(
					"http: matcher %s failed while executing the AfterMockServed() function. mapper: %w",
					aa.(matcher.Matcher).Name(),
					err))
		}
	}

	callbackInput := &CallbackInput{r, parsedURL, parsedBody, h.app, mock, stub}
	for i, callback := range mock.Callbacks {
		err = callback(callbackInput)
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
		input := &PostActionInput{r, parsedURL, parsedBody, h.app, mock, stub, def.RawParameters}
		err = postAction.Run(input)
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

func (h *mockHandler) onNoMatches(w http.ResponseWriter, r *RequestValues, result *lib.FindResult[*HTTPMock]) {
	defer h.lifecycle.OnNoMatch(r, result)

	buf := bytes.Buffer{}

	buf.WriteString("REQUEST DID NOT MATCH\n")

	if result.ClosestMatch != nil {
		buf.WriteString(
			fmt.Sprintf("NEAREST: %s %s\n", result.ClosestMatch.GetID(), result.ClosestMatch.GetName()))
	}

	buf.WriteString("MISSES:\n")

	for _, detail := range result.MismatchDetails {
		buf.WriteString("[")
		buf.WriteString(detail.Target.String())

		if detail.Key != "" {
			buf.WriteString("(")
			buf.WriteString(detail.Key)
			buf.WriteString(")")
		}

		buf.WriteString("] ")
		buf.WriteString(detail.MatchersName)

		if detail.Err != nil {
			buf.WriteString(" ")
			buf.WriteString(detail.Err.Error())
			buf.WriteString("\n")
			continue
		}

		buf.WriteString("(")

		if len(detail.Result.Ext) > 0 {
			buf.WriteString(strings.Join(detail.Result.Ext, ", "))
			buf.WriteString(") ")
			buf.WriteString(detail.Result.Message)
		} else {
			buf.WriteString(detail.Result.Message)
			buf.WriteString(")")
		}

		buf.WriteString("\n")
	}

	w.Header().Add(httpval.HeaderContentType, httpval.MIMETextPlainCharsetUTF8)
	w.WriteHeader(h.app.config.RequestWasNotMatchedStatusCode)
	buf.WriteTo(w)
}

func (h *mockHandler) onError(w http.ResponseWriter, r *RequestValues, err error) {
	defer h.lifecycle.OnError(r, err)

	w.Header().Add(httpval.HeaderContentType, httpval.MIMETextPlainCharsetUTF8)
	w.WriteHeader(h.app.config.RequestWasNotMatchedStatusCode)

	fmt.Fprintf(w, "ERROR DURING REQUEST MATCHING\n%v", err)
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

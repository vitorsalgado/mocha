package dzhttp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/vitorsalgado/mocha/v3/dzhttp/httpval"
	"github.com/vitorsalgado/mocha/v3/dzhttp/internal/httprec"
	"github.com/vitorsalgado/mocha/v3/dzstd"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

type topResponseWriter struct {
	chain dzstd.Chain
	rw    http.ResponseWriter
}

func (w *topResponseWriter) Header() http.Header         { return w.rw.Header() }
func (w *topResponseWriter) Write(p []byte) (int, error) { return w.chain.Next(p) }
func (w *topResponseWriter) WriteHeader(statusCode int)  { w.rw.WriteHeader(statusCode) }

type mockHandler struct {
	app       *HTTPMockApp
	lifecycle mockHTTPLifecycle
}

func newMockHandler(app *HTTPMockApp, lifecycle mockHTTPLifecycle) *mockHandler {
	return &mockHandler{app, lifecycle}
}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parsedURL, urlSegments := h.parseURL(r)
	parsedBody, rawBody, err := parseRequestBody(h.app, r)
	reqValues := &RequestValues{time.Now(), r, parsedURL, urlSegments, nil, nil, h.app, nil}
	if err != nil {
		h.lifecycle.OnWarning(reqValues, fmt.Errorf("request_body_parser: %w", err))
	}

	reqValues.RawBody = rawBody
	reqValues.ParsedBody = parsedBody

	h.lifecycle.OnRequest(reqValues)

	results := &dzstd.Results{Buf: make([]string, 0)}
	selectorInput := &HTTPValueSelectorInput{r, parsedURL, r.URL.Query(), r.Form, parsedBody}
	result, err := dzstd.FindMockForRequest(
		r.Context(),
		h.app.storage,
		func(mock *HTTPMock) []*dzstd.Expectation[*HTTPValueSelectorInput] { return mock.expectations },
		selectorInput,
		results,
		&dzstd.FindOptions{FailFast: h.app.config.FailFast},
	)
	if err != nil {
		h.lifecycle.OnError(reqValues, fmt.Errorf("http: error finding eligible mocks for request. %w", err))
		return
	}

	// no mocks matched with the incoming request
	// lets check if we need to proxy and record this request.
	if !result.Pass {
		if h.app.proxy == nil {
			h.onNoMatches(w, reqValues, result, results)
			return
		}

		// checking if we should record writes to later record mocks.
		if h.app.IsRecording() {
			w = httprec.Wrap(w)
		}

		// proxying non-matched requests.
		h.app.proxy.ServeHTTP(w, r)

		if h.app.IsRecording() {
			res, err := responseFromWriter(w)
			if err != nil {
				h.onError(w, reqValues, err)
				return
			}

			h.app.rec.dispatch(r, parsedURL, rawBody, res)
		}

		return
	}

	mock := result.Matched
	reqValues.Mock = mock

	if mock.Delay != nil {
		ctxTimeout, cancel := context.WithTimeout(r.Context(), mock.Delay())
		defer cancel()

		<-ctxTimeout.Done()
	}

	res, err := result.Matched.Reply.Build(w, reqValues)
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

	if res != nil {
		for i, mapper := range mock.Mappers {
			if err = mapper(reqValues, res); err != nil {
				mock.Dec()
				h.onError(w, reqValues, fmt.Errorf("http: error with response mapper at index %d. %w", i, err))
				return
			}
		}

		if len(mock.Interceptors) > 0 {
			w = &topResponseWriter{dzstd.NewChain(append(mock.Interceptors, &dzstd.RootIntereptor{W: w})), w}
		}

		for k, vv := range res.Header {
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}

		for k := range res.Trailer {
			w.Header().Add(httpval.HeaderTrailer, k)
		}

		for _, cookie := range res.Cookies {
			http.SetCookie(w, cookie)
		}

		w.WriteHeader(res.StatusCode)

		if len(res.Body) > 0 {
			if len(res.Encoding) > 0 {
				buf := bufPool.Get().(*bytes.Buffer)
				buf.Reset()

				_, err = res.encode(buf, bytes.NewReader(res.Body))
				if err != nil {
					h.lifecycle.OnError(reqValues, fmt.Errorf("http: error encoding body: %w", err))
					return
				}

				_, err = buf.WriteTo(w)
				if err != nil {
					h.lifecycle.OnError(reqValues, fmt.Errorf("http: error writting body: %w", err))
					return
				}

				bufPool.Put(buf)
			} else {
				w.Write(res.Body)
			}
		} else if res.BodyCloser != nil {
			defer res.BodyCloser.Close()

			if len(res.Encoding) > 0 {
				_, err = res.encode(w, res.BodyCloser)
				if err != nil {
					h.lifecycle.OnError(reqValues, fmt.Errorf("http: error encoding body: %w", err))
					return
				}
			} else {
				_, err = io.Copy(w, res.BodyCloser)
				if err != nil {
					h.lifecycle.OnError(reqValues, fmt.Errorf("http: error writting body: %w", err))
					return
				}
			}
		}

		for k, v := range res.Trailer {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
	}

	// run matcher events after we sent the mock.
	// some matcher might need to keep state.
	for _, idx := range mock.after {
		err = mock.expectations[idx].Matcher.(matcher.OnMockSent).OnMockSent()
		if err != nil {
			h.lifecycle.OnWarning(reqValues,
				fmt.Errorf("http: after mock served event: matcher[%d] %w", idx, err))
		}
	}

	callbackInput := &CallbackInput{r, parsedURL, parsedBody, h.app, mock, res}
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
		input := &PostActionInput{r, parsedURL, parsedBody, h.app, mock, res, def.RawParameters}
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

	h.lifecycle.OnMatch(reqValues, res)

	// when proxy is enabled, recording happens some steps before
	// so, we need to record here only if proxy is disabled.
	if h.app.IsRecording() && h.app.config.Proxy == nil {
		h.app.rec.dispatch(r, parsedURL, rawBody, res)
	}
}

func (h *mockHandler) onNoMatches(w http.ResponseWriter, r *RequestValues, result *dzstd.FindResult[*HTTPMock], desc *dzstd.Results) {
	defer h.lifecycle.OnNoMatch(r, result, desc)

	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	buf.WriteString("REQUEST DID NOT MATCH\n")
	buf.WriteString(r.RawRequest.Method)
	buf.WriteString(" ")
	buf.WriteString(r.URL.String())
	buf.WriteString("\n\n")

	if result.ClosestMatch != nil {
		fmt.Fprintf(buf, "NEAREST: %s %s\n", result.ClosestMatch.GetID(), result.ClosestMatch.GetName())
	}

	if desc.Len() > 0 {
		fmt.Fprintf(buf, "MISSES(%d):\n", result.MismatchesCount)
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

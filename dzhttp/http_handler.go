package dzhttp

import (
	"bytes"
	"compress/gzip"
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

type topResponseWriter struct {
	chain dzstd.Chain
	rw    http.ResponseWriter
}

func (w *topResponseWriter) Header() http.Header {
	return w.rw.Header()
}

func (w *topResponseWriter) Write(p []byte) (int, error) {
	return w.chain.Next(p)
}

func (w *topResponseWriter) WriteHeader(statusCode int) {
	w.rw.WriteHeader(statusCode)
}

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

	mocks := h.app.storage.GetEligible()
	description := &dzstd.Results{Buf: make([]string, 0, len(mocks))}
	result := dzstd.FindMockForRequest(r.Context(), mocks,
		&HTTPValueSelectorInput{r, parsedURL, r.URL.Query(), r.Form, parsedBody}, description)

	// no mocks matched with the incoming request
	// lets check if we need to proxy and record this request.
	if !result.Pass {
		if h.app.proxy == nil {
			h.onNoMatches(w, reqValues, result, description)
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

		if res.Body != nil {
			if len(res.Encoding) > 0 {
				switch res.Encoding {
				case "gzip":
					buf := new(bytes.Buffer)
					gz := gzipper.Get().(*gzip.Writer)
					defer gzipper.Put(gz)

					gz.Reset(buf)

					_, err = gz.Write(res.Body)
					if err != nil {
						h.lifecycle.OnWarning(reqValues, err)
					}

					err = gz.Close()
					if err != nil {
						h.lifecycle.OnWarning(reqValues, err)
					}

					_, err = buf.WriteTo(w)
					if err != nil {
						h.lifecycle.OnWarning(reqValues, err)
					}
				}
			} else {
				w.Write(res.Body)
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

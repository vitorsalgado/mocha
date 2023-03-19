package mocha

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/internal/httprec"
	"github.com/vitorsalgado/mocha/v3/internal/mimetype"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

type mockHandler struct {
	app       *Mocha
	lifecycle mockHTTPLifecycle
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

	result := findMockForRequest(h.app.storage,
		&valueSelectorInput{r, parsedURL, r.URL.Query(), r.Form, parsedBody})

	if !result.Pass {
		if h.app.proxy != nil {
			// proxy non-matched requests.
			h.app.proxy.ServeHTTP(w, r)

			if h.app.rec != nil {
				res, err := makeStub(w)
				if err != nil {
					h.onError(w, reqValues, err)
					return
				}

				h.app.rec.dispatch(r, parsedURL, rawBody, res)
			}
			return
		} else {
			h.onNoMatches(w, reqValues, result)
			return
		}
	}

	mock := result.Matched
	reqValues.Mock = mock

	if mock.Delay > 0 {
		time.Sleep(mock.Delay)
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

		for k, v := range stub.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}

		for k := range stub.Trailer {
			w.Header().Add(header.Trailer, k)
		}

		for _, cookie := range stub.Cookies {
			http.SetCookie(w, cookie)
		}

		w.WriteHeader(stub.StatusCode)

		if stub.Body != nil {
			w.Write(stub.Body)
		}

		for k, v := range stub.Trailer {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
	} else {
		stub, err = makeStub(w)
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

	input := &PostActionInput{r, parsedURL, parsedBody, h.app, mock, stub, nil}
	for i, action := range mock.PostActions {
		err = action.Run(input)
		if err != nil {
			h.lifecycle.OnWarning(reqValues, fmt.Errorf(
				"http: error running post action at index %d and with type %v. post_action: %w",
				i,
				reflect.TypeOf(action),
				err,
			))
		}
	}

	h.lifecycle.OnMatch(reqValues, stub)

	if h.app.rec != nil && h.app.config.Proxy == nil {
		h.app.rec.dispatch(r, parsedURL, rawBody, stub)
	}
}

func (h *mockHandler) onNoMatches(w http.ResponseWriter, r *RequestValues, result *findResult) {
	defer h.lifecycle.OnNoMatch(r, result)

	builder := strings.Builder{}
	builder.WriteString("REQUEST WAS NOT MATCHED\n")

	if result.ClosestMatch != nil {
		builder.WriteString(
			fmt.Sprintf("CLOSEST MATCH: %s %s\n", result.ClosestMatch.ID, result.ClosestMatch.Name))
	}

	builder.WriteString("MISMATCHES:\n")

	for _, detail := range result.MismatchDetails {
		builder.WriteString("[" + detail.Target.String())
		if detail.Key != "" {
			builder.WriteString("(" + detail.Key + ")")
		}
		builder.WriteString("] ")
		builder.WriteString(detail.MatchersName)

		if detail.Err != nil {
			builder.WriteString(" ")
			builder.WriteString(detail.Err.Error())
			builder.WriteString("\n")
			continue
		}

		builder.WriteString("(")
		if len(detail.Result.Ext) > 0 {
			builder.WriteString(strings.Join(detail.Result.Ext, ", "))
			builder.WriteString(") ")
			builder.WriteString(detail.Result.Message)
		} else {
			builder.WriteString(detail.Result.Message)
			builder.WriteString(")")
		}

		builder.WriteString("\n")
	}

	w.Header().Add(header.ContentType, mimetype.TextPlainCharsetUTF8)
	w.WriteHeader(h.app.config.RequestWasNotMatchedStatusCode)
	w.Write([]byte(builder.String()))
}

func (h *mockHandler) onError(w http.ResponseWriter, r *RequestValues, err error) {
	defer h.lifecycle.OnError(r, err)

	w.Header().Add(header.ContentType, mimetype.TextPlainCharsetUTF8)
	w.WriteHeader(h.app.config.RequestWasNotMatchedStatusCode)
	w.Write([]byte(fmt.Sprintf("ERROR DURING REQUEST MATCHING\n%v", err)))
}

func (h *mockHandler) parseURL(r *http.Request) (u *url.URL, pathSegments []string) {
	reqPath := r.URL.Path
	if len(r.URL.RawQuery) > 0 {
		reqPath += "?" + r.URL.RawQuery
	}

	segments := make([]string, 0)
	urlSegments := strings.Split(reqPath, "/")
	for _, segment := range urlSegments {
		if segment != "" {
			segments = append(segments, segment)
		}
	}

	parsedURL, _ := url.Parse(h.app.URL() + reqPath)

	return parsedURL, segments
}

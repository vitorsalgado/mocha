package mocha

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/vitorsalgado/mocha/v3/internal/colorize"
	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/internal/httpx"
	"github.com/vitorsalgado/mocha/v3/internal/mimetype"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

type mockHTTPLifecycle interface {
	OnRequest(*RequestValues)
	OnMatch(*RequestValues, *Stub)
	OnNoMatch(*RequestValues, *findResult)
	OnWarning(*RequestValues, error)
	OnError(*RequestValues, error)
}

type mockHandler struct {
	app       *Mocha
	lifecycle mockHTTPLifecycle
}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w = httpx.Wrap(w)
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
				h.onError(
					w,
					reqValues,
					fmt.Errorf("http: error with response mapper at index %d. %w", i, err),
				)
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
			builder.WriteString("[" + detail.Key + "]")
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

// Standard HTTP Lifecycle

type builtInMockHTTPLifecycle struct {
	app *Mocha
}

func (h *builtInMockHTTPLifecycle) OnRequest(r *RequestValues) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	evt := h.app.log.Info().Str("url", r.URL.String())
	req := zerolog.Dict()

	if h.app.config.LogVerbosity >= LogHeader {
		req.Any("header", r.RawRequest.Header)
	}

	if h.app.config.LogVerbosity >= LogBody && len(r.RawBody) > 0 {
		req.Bytes("body", r.RawBody)
	}

	evt.Dict("request", req)
	evt.Msgf("---> REQUEST RECEIVED %s %s", r.RawRequest.Method, r.URL.Path)
}

func (h *builtInMockHTTPLifecycle) OnMatch(r *RequestValues, s *Stub) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	evt := h.app.log.Info().
		Str("url", r.URL.String()).
		Dur("elapsed", time.Since(r.StartedAt))

	dict := zerolog.Dict().Str("id", r.Mock.ID)
	if r.Mock.Name != "" {
		dict.Str("name", r.Mock.Name)
	}

	res := zerolog.Dict().Int("status", s.StatusCode)

	if h.app.config.LogVerbosity >= LogHeader {
		res.Any("header", s.Header)
		if len(s.Trailer) > 0 {
			res.Any("trailer", s.Trailer)
		}
	}

	bodyLen := int64(len(s.Body))
	if h.app.config.LogVerbosity >= LogBody && bodyLen > 0 &&
		(h.app.config.LogBodyMaxSize == 0 || bodyLen <= h.app.config.LogBodyMaxSize) {
		res.Bytes("body", s.Body)
		if s.Encoding != "" {
			res.Str("encoding", s.Encoding)
		}
	}

	evt.Dict("mock", dict)
	evt.Dict("response", res)

	evt.Msgf("<--- REQUEST MATCHED %s %s", r.RawRequest.Method, r.URL.Path)
}

func (h *builtInMockHTTPLifecycle) OnNoMatch(r *RequestValues, fr *findResult) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	evt := h.app.log.Warn().Str("url", r.URL.String())
	if fr.ClosestMatch != nil {
		md := zerolog.Dict().Str("id", fr.ClosestMatch.ID)
		if fr.ClosestMatch.Name != "" {
			md.Str("name", fr.ClosestMatch.Name)
		}

		evt.Dict("closest_match", md)
	}

	if len(fr.MismatchDetails) > 0 {
		mismatches := zerolog.Arr()
		builder := strings.Builder{}

		for _, detail := range fr.MismatchDetails {
			builder.WriteString("[" + detail.Target.String())
			if detail.Key != "" {
				builder.WriteString("[" + detail.Key + "]")
			}
			builder.WriteString("] ")

			builder.WriteString(detail.MatchersName)

			if detail.Err != nil {
				builder.WriteString(" ")
				builder.WriteString(detail.Err.Error())
				continue
			}

			builder.WriteString("(")
			if len(detail.Result.Ext) > 0 {
				builder.WriteString(strings.Join(detail.Result.Ext, ", "))
				builder.WriteString(") ")
				builder.WriteString(detail.Result.Message)
			} else {
				builder.WriteString(detail.Result.Message)
				builder.WriteString(") ")
			}

			mismatches.Str(builder.String())
			builder.Reset()
		}

		evt.Array("mismatches", mismatches)
	}

	evt.Msgf("<--- REQUEST WAS NOT MATCHED %s %s", r.RawRequest.Method, r.URL.Path)
}

func (h *builtInMockHTTPLifecycle) OnWarning(r *RequestValues, err error) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	h.app.log.Warn().
		Err(err).
		Str("url", r.URL.String()).
		Msgf("<--- WARNING %s %s", r.RawRequest.Method, r.URL.Path)
}

func (h *builtInMockHTTPLifecycle) OnError(r *RequestValues, err error) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	h.app.log.Error().
		Err(err).
		Str("url", r.URL.String()).
		Msgf("<--- ERROR %s %s", r.RawRequest.Method, r.URL.Path)
}

// Descriptive Logger

type builtInDescriptiveMockHTTPLifecycle struct {
	app *Mocha
	cz  *colorize.Colorize
}

func (h *builtInDescriptiveMockHTTPLifecycle) OnRequest(e *RequestValues) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s %s ---> %s %s\n%s %s",
		h.cz.BlueBright(h.cz.Bold("REQUEST RECEIVED")),
		e.StartedAt.Format(time.RFC3339),
		h.cz.Blue(e.RawRequest.Method),
		h.cz.Blue(e.URL.Path),
		e.RawRequest.Method,
		e.URL))

	if h.app.config.LogVerbosity >= LogHeader {
		builder.WriteString("\n")
		builder.WriteString(h.cz.Blue("Headers "))
		builder.WriteString(fmt.Sprintf("%s", e.RawRequest.Header))
	}

	if h.app.config.LogVerbosity >= LogBody && len(e.RawBody) > 0 {
		builder.WriteString("\n")
		builder.WriteString(h.cz.Blue("Body: "))
		builder.WriteString(fmt.Sprintf("%v\n", string(e.RawBody)))
	}

	fmt.Println(builder.String())
}

func (h *builtInDescriptiveMockHTTPLifecycle) OnMatch(e *RequestValues, s *Stub) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s %s <--- %s %s\n%s ",
		h.cz.GreenBright(h.cz.Bold("REQUEST MATCHED")),
		time.Now().Format(time.RFC3339),
		h.cz.Green(e.RawRequest.Method),
		h.cz.Green(e.URL.Path),
		e.RawRequest.Method))
	builder.WriteString(e.URL.String())
	builder.WriteString("\n")
	builder.WriteString(h.cz.Bold("Mock: "))
	builder.WriteString(e.Mock.ID + " " + e.Mock.Name)
	builder.WriteString("\n")
	builder.WriteString(h.cz.Green("Took(ms): "))
	builder.WriteString(strconv.FormatInt(time.Since(e.StartedAt).Milliseconds(), 10))
	builder.WriteString("\n")
	builder.WriteString(h.cz.Green("Status: "))
	builder.WriteString(strconv.FormatInt(int64(s.StatusCode), 10))

	if h.app.config.LogVerbosity >= LogHeader {
		builder.WriteString("\n")
		builder.WriteString(h.cz.Green("Headers: "))
		builder.WriteString(fmt.Sprintf("%s", s.Header))
	}

	bodyLen := int64(len(s.Body))
	if h.app.config.LogVerbosity >= LogBody && bodyLen > 0 &&
		(h.app.config.LogBodyMaxSize == 0 || bodyLen <= h.app.config.LogBodyMaxSize) {
		builder.WriteString("\n")
		builder.WriteString(h.cz.Green("Body: "))
		if s.Encoding == "" {
			builder.WriteString(string(s.Body))
		} else {
			builder.WriteString("<encoded body omitted>")
		}
	}

	builder.WriteString("\n")

	fmt.Println(builder.String())
}

func (h *builtInDescriptiveMockHTTPLifecycle) OnNoMatch(r *RequestValues, fr *findResult) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("\n%s %s <--- %s %s\n%s %s\n",
		h.cz.YellowBright(h.cz.Bold("REQUEST WAS NOT MATCHED")),
		time.Now().Format(time.RFC3339),
		h.cz.Yellow(r.RawRequest.Method),
		h.cz.Yellow(r.URL.Path),
		r.RawRequest.Method,
		r.URL.String()))

	if fr.ClosestMatch != nil {
		builder.WriteString(fmt.Sprintf("%s: %s %s\n\n",
			h.cz.Bold("Closest Match"), fr.ClosestMatch.ID, fr.ClosestMatch.Name))
	}

	if len(fr.MismatchDetails) > 0 {
		if h.app.config.LogVerbosity <= LogHeader {
			builder.WriteString(fmt.Sprintf("%s: %d", h.cz.Bold("Mismatches"), len(fr.MismatchDetails)))
		} else {
			builder.WriteString(fmt.Sprintf("%s(%d):\n", h.cz.Bold("Mismatches"), len(fr.MismatchDetails)))
			for _, detail := range fr.MismatchDetails {
				builder.WriteString("[" + detail.Target.String())
				if detail.Key != "" {
					builder.WriteString("[" + detail.Key + "]")
				}
				builder.WriteString("] ")
				builder.WriteString(detail.MatchersName)
				builder.WriteString("(")

				if detail.Err != nil {
					builder.WriteString(detail.Err.Error())
					builder.WriteString(")\n")
					continue
				}

				if detail.Result == nil {
					builder.WriteString(")\n")
					continue
				}

				if len(detail.Result.Ext) == 0 {
					builder.WriteString(h.cz.Bold(detail.Result.Message))
					builder.WriteString(")")
				} else {
					builder.WriteString(h.cz.Bold(strings.Join(detail.Result.Ext, ", ")))
					builder.WriteString(") ")
					builder.WriteString(detail.Result.Message)
				}

				builder.WriteString("\n")
			}
		}
	}

	fmt.Println(builder.String())
}

func (h *builtInDescriptiveMockHTTPLifecycle) OnWarning(r *RequestValues, err error) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	fmt.Printf("%s %s <--- %s %s\n%s %s\n%s: %v\n",
		h.cz.RedBright(h.cz.Bold("WARNING")),
		time.Now().Format(time.RFC3339),
		h.cz.Red(r.RawRequest.Method),
		h.cz.Red(r.URL.Path),
		r.RawRequest.Method,
		r.URL.String(),
		h.cz.Red(h.cz.Bold("Error: ")),
		err,
	)
}

func (h *builtInDescriptiveMockHTTPLifecycle) OnError(r *RequestValues, err error) {
	if h.app.config.LogLevel == LogLevelDisabled {
		return
	}

	fmt.Printf("%s %s <--- %s %s\n%s %s\n%s: %v\n",
		h.cz.RedBright(h.cz.Bold("ERROR")),
		time.Now().Format(time.RFC3339),
		h.cz.Red(r.RawRequest.Method),
		h.cz.Red(r.URL.Path),
		r.RawRequest.Method,
		r.URL.String(),
		h.cz.Red(h.cz.Bold("Error: ")),
		err,
	)
}

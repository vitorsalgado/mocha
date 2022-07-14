package mocha

import (
	"bufio"
	"fmt"
	"net/http"
	"time"

	"github.com/vitorsalgado/mocha/core"
	"github.com/vitorsalgado/mocha/expect"
	"github.com/vitorsalgado/mocha/internal/colorize"
	"github.com/vitorsalgado/mocha/internal/parameters"
	"github.com/vitorsalgado/mocha/x/headers"
	"github.com/vitorsalgado/mocha/x/mimetypes"
)

type mockHandler struct {
	mocks       core.Storage
	bodyParsers []RequestBodyParser
	params      parameters.Params
	t           core.T
}

func newHandler(
	storage core.Storage,
	bodyParsers []RequestBodyParser,
	params parameters.Params,
	t core.T,
) *mockHandler {
	return &mockHandler{mocks: storage, bodyParsers: bodyParsers, params: params, t: t}
}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.t.Helper()

	start := time.Now()

	h.t.Logf("\n%s %s ---> %s %s\n%s %s\n\n%s:\n %s: %v\n",
		colorize.BlueBright(colorize.Bold("REQUEST RECEIVED")),
		start.Format(time.RFC3339),
		colorize.Blue(r.Method),
		colorize.Blue(r.URL.Path),
		r.Method,
		fullURL(r.Host, r.RequestURI),
		colorize.Blue("Request"),
		colorize.Blue("Headers"),
		r.Header,
	)

	parsedBody, err := parseRequestBody(r, h.bodyParsers)
	if err != nil {
		respondError(w, r, h.t, err)
		return
	}

	// match current request with all eligible stored matchers in order to find one mock.
	args := expect.Args{
		RequestInfo: &expect.RequestInfo{Request: r, ParsedBody: parsedBody},
		Params:      h.params}
	result, err := core.FindMockForRequest(h.mocks, args)
	if err != nil {
		respondError(w, r, h.t, err)
		return
	}

	if !result.Matches {
		respondNonMatched(w, r, result, h.t)
		return
	}

	m := result.Matched
	m.Hit()

	// run post matchers, after standard ones and after marking the mock as called.
	afterResult, err := m.Matches(args, m.PostExpectations)
	if err != nil {
		respondError(w, r, h.t, err)
		return
	}

	if !afterResult.IsMatch {
		respondNonMatched(w, r, result, h.t)
		return
	}

	// get the reply for the mock, after running all possible matchers.
	res, err := result.Matched.Reply.Build(r, m, h.params)
	if err != nil {
		h.t.Logf(err.Error())
		respondError(w, r, h.t, err)
		return
	}

	// map the response using mock mappers.
	mapperArgs := core.ResponseMapperArgs{Request: r, Parameters: h.params}
	for _, mapper := range res.Mappers {
		if err = mapper(res, mapperArgs); err != nil {
			respondError(w, r, h.t, err)
			return
		}
	}

	// if a delay is set, it will wait before continuing serving the mocked response.
	if res.Delay > 0 {
		<-time.After(res.Delay)
	}

	for k := range res.Header {
		w.Header().Add(k, res.Header.Get(k))
	}

	w.WriteHeader(res.Status)

	if res.Body != nil {
		scanner := bufio.NewScanner(res.Body)
		for scanner.Scan() {
			w.Write(scanner.Bytes())
		}

		if scanner.Err() != nil {
			h.t.Logf("error writing response body: error=%v", scanner.Err())
		}
	}

	// run post actions.
	paArgs := core.PostActionArgs{Request: r, Response: res, Mock: m, Params: h.params}
	for i, action := range m.PostActions {
		err = action.Run(paArgs)
		if err != nil {
			h.t.Logf("\nan error occurred running post action %d. error=%v", i, err)
		}
	}

	elapsed := time.Since(start)

	h.t.Logf("\n%s %s <--- %s %s\n%s %s\n\n%s%d - %s\n\n%s: %dms\n%s:\n %s: %d\n %s: %v\n",
		colorize.GreenBright(colorize.Bold("REQUEST DID MATCH")),
		time.Now().Format(time.RFC3339),
		colorize.Green(r.Method),
		colorize.Green(r.URL.Path),
		r.Method,
		fullURL(r.Host, r.RequestURI),
		colorize.Bold("Mock: "),
		m.ID,
		m.Name,
		colorize.Green("Took"),
		elapsed.Milliseconds(),
		colorize.Green("Response Definition"),
		colorize.Green("Status"),
		res.Status,
		colorize.Green("Headers"),
		res.Header,
	)
}

func respondNonMatched(w http.ResponseWriter, r *http.Request, result *core.FindResult, t core.T) {
	t.Logf("\n%s %s <--- %s %s\n%s %s\n\n",
		colorize.YellowBright(colorize.Bold("REQUEST DID NOT MATCH")),
		time.Now().Format(time.RFC3339),
		colorize.Yellow(r.Method),
		colorize.Yellow(r.URL.Path),
		r.Method,
		fullURL(r.Host, r.RequestURI),
	)

	w.Header().Add(headers.ContentType, mimetypes.TextPlain)
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte("REQUEST DID NOT MATCH.\n"))

	if result.ClosestMatch != nil {
		format := "id: %d\nname: %s\n\n"
		w.Write([]byte("Closest Match:\n"))
		w.Write([]byte(fmt.Sprintf(format, result.ClosestMatch.ID, result.ClosestMatch.Name)))

		t.Logf(colorize.Yellow("Closest Match:\n"))
		t.Logf(format, result.ClosestMatch.ID, result.ClosestMatch.Name)
	}

	w.Write([]byte("Mismatches:\n"))
	t.Logf(colorize.Yellow("Mismatches:\n"))

	for _, detail := range result.MismatchDetails {
		w.Write([]byte(
			fmt.Sprintf("Matcher: %s\nTarget: %s\nReason: %s\n\n",
				detail.Name, detail.Target, detail.Description)))

		t.Logf(fmt.Sprintf("%s\n%s\n%s",
			fmt.Sprintf("Matcher \"%s\"", colorize.Bold(detail.Name)),
			"Target: "+detail.Target,
			detail.Description))
	}
}

func respondError(w http.ResponseWriter, r *http.Request, t core.T, err error) {
	t.Logf("\n%s %s <--- %s %s\n%s %s\n\n%s: %v",
		colorize.RedBright(colorize.Bold("REQUEST DID NOT MATCH")),
		time.Now().Format(time.RFC3339),
		colorize.Red(r.Method),
		colorize.Red(r.URL.Path),
		r.Method,
		fullURL(r.Host, r.RequestURI),
		colorize.Red(colorize.Bold("Error: ")),
		err,
	)

	w.Header().Add(headers.ContentType, mimetypes.TextPlain)
	w.WriteHeader(http.StatusTeapot)

	w.Write([]byte("Request did not match. An error occurred.\n"))
	w.Write([]byte(fmt.Sprintf("%v", err)))
}

func fullURL(host, uri string) string {
	return fmt.Sprintf("%s%s", host, uri)
}

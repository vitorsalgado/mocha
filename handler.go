package mocha

import (
	"bufio"
	"log"
	"net/http"
	"time"

	"github.com/vitorsalgado/mocha/internal/params"
	"github.com/vitorsalgado/mocha/matcher"
	"github.com/vitorsalgado/mocha/mock"
)

type Handler struct {
	mocks   mock.Storage
	parsers []BodyParser
	params  params.Params
}

func newHandler(
	storage mock.Storage,
	parsers []BodyParser,
	params params.Params,
) *Handler {
	return &Handler{mocks: storage, parsers: parsers, params: params}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parseRequestBody, err := ParseRequestBody(r, h.parsers)
	if err != nil {
		respondErr(w, err)
		return
	}

	parameters := matcher.Args{RequestInfo: &matcher.RequestInfo{Request: r, ParsedBody: parseRequestBody}, Params: h.params}
	result, err := findMockForRequest(h.mocks, parameters)
	if err != nil {
		respondErr(w, err)
		return
	}

	if !result.Matches {
		noMatch(w, result)
		return
	}

	m := result.Matched
	res, err := result.Matched.Reply.Build(r, m, h.params)
	if err != nil {
		respondErr(w, err)
		return
	}

	mp := mock.ResponseMapperArgs{Request: r, Parameters: h.params}
	for _, mapper := range res.Mappers {
		if err := mapper(res, mp); err != nil {
			respondErr(w, err)
			return
		}
	}

	m.Hit()

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
			panic(scanner.Err())
		}
	}

	args := mock.PostActionArgs{Request: r, Response: res, Mock: m, Params: h.params}
	for _, action := range m.PostActions {
		err := action.Run(args)
		if err != nil {
			log.Println(err)
		}
	}
}

func noMatch(w http.ResponseWriter, result *findMockResult) {
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte("Request was not matched."))

	if result.ClosestMatch != nil {
		w.Write([]byte("\n"))
		w.Write([]byte("\n"))
	}
}

func respondErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte("Request was not matched."))
	w.Write([]byte(err.Error()))

	log.Printf("Reason: %v", err)
}

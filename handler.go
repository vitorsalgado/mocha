package mocha

import (
	"bufio"
	"log"
	"net/http"
)

type Handler struct {
	repo               MockRepository
	scenarioRepository ScenarioRepository
	bodyParsers        []BodyParser
}

func newHandler(
	mockstore MockRepository,
	scenariostore ScenarioRepository, parsers []BodyParser,
) *Handler {
	return &Handler{repo: mockstore, scenarioRepository: scenariostore, bodyParsers: parsers}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := WrapRequest(r, h.bodyParsers)
	if err != nil {
		respondErr(w, err)
		return
	}

	params := MatcherParams{Req: req, Repo: h.repo, ScenarioRepository: h.scenarioRepository}
	result, err := FindMockForRequest(params)

	if err != nil {
		respondErr(w, err)
		return
	}

	if !result.Matches {
		noMatch(w, result)
		return
	}

	mock := result.Matched
	res, err := result.Matched.ResFn(r, mock)
	if err != nil {
		respondErr(w, err)
		return
	}

	mock.Hits++
	h.repo.Save(mock)

	for k, v := range res.Headers {
		w.Header().Add(k, v)
	}

	w.WriteHeader(res.Status)

	if res.Body == nil {
		return
	}

	scanner := bufio.NewScanner(res.Body)

	for scanner.Scan() {
		_, err = w.Write(scanner.Bytes())
		if err != nil {
			respondErr(w, err)
		}
	}

	if scanner.Err() != nil {
		respondErr(w, err)
		return
	}
}

func noMatch(w http.ResponseWriter, result *FindMockResult) {
	w.WriteHeader(http.StatusTeapot)
	_, _ = w.Write([]byte("Request was not matched."))

	if result.ClosestMatch != nil {
		_, _ = w.Write([]byte("\n"))
		_, _ = w.Write([]byte("\n"))
	}
}

func respondErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusTeapot)
	_, _ = w.Write([]byte("Request was not matched."))
	_, _ = w.Write([]byte(err.Error()))

	log.Printf("Reason: %v", err)
}

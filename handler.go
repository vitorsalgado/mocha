package mocha

import (
	"bufio"
	"log"
	"net/http"
)

type Handler struct {
	repo               MockRepository
	scenarioRepository ScenarioRepository
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := MatcherContext{Req: r, Repo: h.repo, ScenarioRepository: h.scenarioRepository}
	result, err := FindMockForRequest(ctx)

	if err != nil {
		respondErr(w, err)
		return
	}

	if result.Matches {
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

		return
	}

	noMatch(w, result)
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

	return
}

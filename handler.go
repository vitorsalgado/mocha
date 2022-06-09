package mocha

import (
	"github.com/vitorsalgado/mocha/base"
	"log"
	"net/http"
	"reflect"

	"github.com/vitorsalgado/mocha/internal/arrays"
)

type Handler struct {
	repo MockRepository
}

type Result struct {
	IsMatch bool
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mocks := h.repo.GetAllSorted()

	for _, mock := range mocks {
		if arrays.All(mock.Expectations, func(e any) bool {
			var err error
			var res bool

			switch e := e.(type) {
			case Expectation[string]:
				res, err = e.Matcher(e.Pick(r), base.MatcherContext{})
			default:
				log.Fatalf("unhandled matcher type %s", reflect.TypeOf(e))
			}

			if err != nil {
				log.Fatal(e)
			}

			return res
		}) {
			res, err := mock.ResFn(r, mock)
			if err != nil {
				log.Fatal(err)
			}

			mock.Hits++

			w.WriteHeader(res.Status)

			_, _ = w.Write(res.Body)

			for k, v := range res.Headers {
				w.Header().Add(k, v)
			}
		}
	}
}

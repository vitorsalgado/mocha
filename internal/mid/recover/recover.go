package recover

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func New(fn func(err error), status int) *Recover {
	return &Recover{fn: fn, status: status}
}

type Recover struct {
	status int
	fn     func(err error)
}

func (h *Recover) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovery := recover(); recovery != nil {
				w.WriteHeader(h.status)
				_, _ = fmt.Fprintf(w, "%d - Unexpected Error!\nPanic: %v", h.status, recovery)

				h.fn(fmt.Errorf(
					"http: panic during request matching. %s %s. %v\n%s",
					r.Method,
					r.URL.String(),
					recovery,
					debug.Stack(),
				))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

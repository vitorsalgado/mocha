package recover

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

type L interface {
	Logf(string, ...any)
}

func New(l L, status int) *Recover {
	return &Recover{l: l, status: status}
}

type Recover struct {
	l      L
	status int
}

func (h *Recover) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovery := recover(); recovery != nil {
				err := fmt.Errorf("panic=%v\n %s", recovery, debug.Stack())

				w.WriteHeader(h.status)
				w.Write([]byte(fmt.Sprintf(
					"%d - Unexpected Error!\n %v",
					h.status,
					fmt.Sprintf("panic=%v", recovery),
				)))

				h.l.Logf(err.Error())
			}
		}()

		next.ServeHTTP(w, r)
	})
}

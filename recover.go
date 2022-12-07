package mocha

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/internal/mimetype"
)

type recoverMid struct {
	d   Debug
	t   TestingT
	evt *eventListener
}

func (h *recoverMid) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovery := recover(); recovery != nil {
				buf := make([]byte, 1024)
				buf = buf[:runtime.Stack(buf, false)]

				err := fmt.Errorf("panic=%v\n%s\n", recovery, buf)

				w.Header().Set(header.ContentType, mimetype.TextPlain)
				w.WriteHeader(http.StatusTeapot)
				w.Write([]byte(fmt.Sprintf(
					"%d (%s) - Unexpected Error!\n\n%v",
					http.StatusTeapot,
					http.StatusText(http.StatusTeapot),
					err,
				)))

				h.evt.Emit(&OnError{Request: evtRequest(r), Err: err})
				h.t.Logf(err.Error())

				if h.d != nil {
					h.d(err)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}

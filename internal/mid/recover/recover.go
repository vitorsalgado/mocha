package recover

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/vitorsalgado/mocha/v3/internal/header"
	"github.com/vitorsalgado/mocha/v3/internal/mimetype"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovery := recover(); recovery != nil {
				buf := make([]byte, 1024)
				buf = buf[:runtime.Stack(buf, false)]

				log.Printf("panic=%v\n%s\n", recovery, buf)

				w.Header().Set(header.ContentType, mimetype.TextPlain)
				w.WriteHeader(http.StatusTeapot)
				w.Write([]byte(fmt.Sprintf(
					"%d (%s) - Panic!\n\n%v",
					http.StatusTeapot,
					http.StatusText(http.StatusTeapot),
					recovery,
				)))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

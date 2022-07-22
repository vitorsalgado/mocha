package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/vitorsalgado/mocha/x/headers"
	"github.com/vitorsalgado/mocha/x/mimetypes"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovery := recover(); recovery != nil {
				buf := make([]byte, 1024)
				buf = buf[:runtime.Stack(buf, false)]

				msg := fmt.Sprintf("panic: %v\n%s\n", recovery, buf)

				log.Printf(msg)

				w.Header().Set(headers.ContentType, mimetypes.TextPlain)
				w.WriteHeader(http.StatusTeapot)
				w.Write([]byte(fmt.Sprintf("%s - %s", http.StatusText(http.StatusTeapot), msg)))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

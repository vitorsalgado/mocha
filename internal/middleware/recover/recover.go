package recover

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/vitorsalgado/mocha/v2/internal/headers"
	"github.com/vitorsalgado/mocha/v2/internal/mimetypes"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovery := recover(); recovery != nil {
				buf := make([]byte, 1024)
				buf = buf[:runtime.Stack(buf, false)]

				log.Printf("panic: %v\n%s\n", recovery, buf)

				w.Header().Set(headers.ContentType, mimetypes.TextPlain)
				w.WriteHeader(http.StatusTeapot)
				w.Write([]byte(fmt.Sprintf("%s - an unexpected error has occurred", http.StatusText(http.StatusTeapot))))
				w.Write([]byte(fmt.Sprintf("%v", recovery)))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

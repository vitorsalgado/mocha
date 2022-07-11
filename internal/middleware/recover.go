package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vitorsalgado/mocha/x/headers"
	"github.com/vitorsalgado/mocha/x/mimetypes"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovery := recover(); recovery != nil {
				msg := fmt.Sprintf("an unexpected error occured. %v", recovery)

				log.Printf("panic: %v", recovery)

				w.Header().Set(headers.ContentType, mimetypes.TextPlain)
				w.WriteHeader(http.StatusTeapot)
				w.Write([]byte(fmt.Sprintf("%s - %s", http.StatusText(http.StatusTeapot), msg)))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

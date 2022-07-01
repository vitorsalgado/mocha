package middleware

import (
	"fmt"
	"log"
	"net/http"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovery := recover(); recovery != nil {
				msg := fmt.Sprintf("an unexpected error occured. %v", recovery)

				log.Printf("panic: %v", recovery)

				w.Header().Set("content-type", "text/plain")
				w.WriteHeader(http.StatusTeapot)
				w.Write([]byte(fmt.Sprintf("%s - %s", http.StatusText(http.StatusTeapot), msg)))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

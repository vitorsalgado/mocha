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
				var msg string
				switch x := recovery.(type) {
				case string:
					msg = x
				case error:
					msg = x.Error()
				default:
					msg = fmt.Sprintf("an unexpected error occured. %v", x)
				}

				log.Printf("panic: %v", recovery)

				w.Header().Set("content-type", "text/plain")
				w.WriteHeader(http.StatusTeapot)
				_, _ = w.Write([]byte(fmt.Sprintf("%s - %s", http.StatusText(http.StatusTeapot), msg)))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

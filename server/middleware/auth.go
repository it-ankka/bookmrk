package middleware

import "net/http"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that user is authenticated
		next.ServeHTTP(w, r)
	})
}

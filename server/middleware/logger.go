package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

func Logger(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
		body, err := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		if err != nil {
			logger.Info(s, "body", string(body))
		} else {
			logger.Info(s)
		}
		next.ServeHTTP(w, r)
	})
}

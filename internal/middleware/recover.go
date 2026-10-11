package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic recovered",
					slog.Any("error", err),
					slog.String("stack", string(debug.Stack())),
				)
				http.Error(w, `{"error":"internal_error","message":"internal server error"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

func requestLogger(ctx context.Context) func(next http.Handler) http.Handler {
	logger := zerolog.Ctx(ctx)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrapped := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(wrapped, r)
			status := wrapped.Status()
			var lvl zerolog.Level
			if status < 300 {
				lvl = zerolog.InfoLevel
			} else if status < 400 {
				lvl = zerolog.WarnLevel
			} else {
				lvl = zerolog.ErrorLevel
			}
			logger.WithLevel(lvl).Int("status", status).Msgf("%s %s", r.Method, r.URL.Path)
		})
	}
}

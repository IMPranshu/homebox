package mid

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

type spy struct {
	http.ResponseWriter
	status int
}

func (s *spy) WriteHeader(status int) {
	s.status = status
	s.ResponseWriter.WriteHeader(status)
}

func Logger(l zerolog.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Context().Value(middleware.RequestIDKey).(string)

			l.Info().Str("method", r.Method).Str("path", r.URL.Path).Str("rid", reqID).Msg("request received")

			s := &spy{ResponseWriter: w}
			h.ServeHTTP(s, r)

			l.Info().Str("method", r.Method).Str("path", r.URL.Path).Int("status", s.status).Str("rid", reqID).Msg("request finished")
		})
	}
}

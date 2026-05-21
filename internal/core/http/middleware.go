package corehttp

import (
	"net/http"

	corelogger "github.com/egotk/golang-advert-app/internal/core/logger"
	"github.com/google/uuid"
)

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(
	h http.Handler,
	m ...Middleware,
) http.Handler {
	if len(m) == 0 {
		return h
	}

	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}

	return h
}

const requestIDHeader = "X-Request-ID"

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log corelogger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			l := log.With(
				corelogger.String("request_id", requestID),
				corelogger.String("url", r.URL.String()),
			)

			ctx := corelogger.ToContext(r.Context(), l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

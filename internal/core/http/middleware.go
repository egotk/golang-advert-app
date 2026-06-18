package corehttp

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corehttpresponse "github.com/egotk/golang-advert-app/internal/core/http/response"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	"github.com/google/uuid"
	"go.uber.org/zap"
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

func Logger(log *corezaplogger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := corezaplogger.ToContext(r.Context(), l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

const authHeader = "Authorization"

type JWTService interface {
	ParseAccessToken(access string) (corejwt.Claims, error)
}

func JWToken(jwtService JWTService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			responseHandler := corehttpresponse.From(ctx, rw)

			auth := r.Header.Get(authHeader)
			if auth == "" {
				responseHandler.ErrorResponse(
					coreerrors.ErrUnauthorized,
					"failed to get Authorization header",
				)

				return
			}

			authParts := strings.Split(auth, " ")
			if len(authParts) != 2 || !strings.EqualFold(authParts[0], "Bearer") {
				responseHandler.ErrorResponse(
					coreerrors.ErrUnauthorized,
					"failed to parse Authorization header",
				)

				return
			}

			tokenString := authParts[1]

			claims, err := jwtService.ParseAccessToken(tokenString)
			if err != nil {
				responseHandler.ErrorResponse(
					coreerrors.ErrUnauthorized,
					"failed to parse JWT",
				)

				return
			}

			ctx = corejwt.ClaimsToContext(ctx, claims)

			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}

func Role(requiredRoles ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			responseHandler := corehttpresponse.From(ctx, rw)

			claims, err := corejwt.ClaimsFromContext(ctx)
			if err != nil {
				responseHandler.ErrorResponse(
					coreerrors.ErrUnauthorized,
					"failed to get claims",
				)

				return
			}

			role := claims.Role

			if slices.Contains(requiredRoles, role) {
				responseHandler.ErrorResponse(
					coreerrors.ErrForbidden,
					fmt.Sprintf("role %s not allowed", role),
				)

				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}

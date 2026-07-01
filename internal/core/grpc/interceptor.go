package coregrpc

import (
	"context"
	"errors"
	"fmt"
	"slices"

	corectx "github.com/egotk/golang-advert-app/internal/core/ctx"
	coreerrors "github.com/egotk/golang-advert-app/internal/core/errors"
	corejwt "github.com/egotk/golang-advert-app/internal/core/jwt"
	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func RequestID() grpc.UnaryServerInterceptor {
	const requestIDHeader = "x-request-id"

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		requestID := GetRequestIDFromMetadata(ctx)

		grpc.SetHeader(ctx, metadata.Pairs(requestIDHeader, requestID))
		ctx = corectx.RequestIDToContext(ctx, requestID)

		return handler(ctx, req)
	}
}

func Logger(log *corezaplogger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		requestID, err := corectx.RequestIDFromContext(ctx)
		if err != nil {
			requestID = "unknown"
		}

		l := log.With(
			zap.String("request_id", requestID),
			zap.String("method", info.FullMethod),
		)

		ctx = corezaplogger.ToContext(ctx, l)

		return handler(ctx, req)
	}
}

type JWTService interface {
	ParseAccessToken(access string) (corejwt.Claims, error)
}

func JWToken(jwtService JWTService, appliedToMethods ...string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		required := slices.Contains(appliedToMethods, info.FullMethod)
		if !required {
			return handler(ctx, req)
		}

		token, err := GetTokenFromMetadata(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get token: %w", coreerrors.ErrUnauthorized)
		}

		claims, err := jwtService.ParseAccessToken(token)
		if err != nil {
			return nil, fmt.Errorf("failed to get claims: %w", coreerrors.ErrUnauthorized)
		}

		ctx = corejwt.ClaimsToContext(ctx, claims)

		return handler(ctx, req)
	}
}

func Role(methodRoles map[string][]string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		requiredRoles, ok := methodRoles[info.FullMethod]
		if !ok {
			return handler(ctx, req)
		}

		claims, err := corejwt.ClaimsFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get claims: %w", coreerrors.ErrUnauthorized)
		}

		role := claims.Role

		if !slices.Contains(requiredRoles, role) {
			return nil, fmt.Errorf("role %s not allowed: %w", role, coreerrors.ErrForbidden)
		}

		return handler(ctx, req)
	}
}

func ErrorHandler() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		resp, err = handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		log := corezaplogger.FromContext(ctx)

		var (
			code    codes.Code
			logFunc func(string, ...zap.Field)
		)

		switch {
		case errors.Is(err, coreerrors.ErrInvalidArgument):
			code = codes.InvalidArgument
			logFunc = log.Warn

		case errors.Is(err, coreerrors.ErrNotFound):
			code = codes.NotFound
			logFunc = log.Debug

		case errors.Is(err, coreerrors.ErrConflict):
			code = codes.AlreadyExists
			logFunc = log.Warn

		case errors.Is(err, coreerrors.ErrUnauthorized):
			code = codes.Unauthenticated
			logFunc = log.Warn

		case errors.Is(err, coreerrors.ErrForbidden):
			code = codes.PermissionDenied
			logFunc = log.Debug

		default:
			code = codes.Internal
			logFunc = log.Error
		}

		logFunc("request failed", zap.Uint32("code", uint32(code)), zap.String("method", info.FullMethod), zap.Error(err))

		return nil, status.Error(code, err.Error())
	}
}

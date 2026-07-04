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

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *wrappedStream) Context() context.Context {
	return s.ctx
}

const requestIDHeader = "x-request-id"

func RequestID() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		ctx = requestIDToContext(ctx)

		return handler(ctx, req)
	}
}

func RequestIDStream() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := requestIDToContext(ss.Context())

		ws := &wrappedStream{
			ServerStream: ss,
			ctx:          ctx,
		}

		return handler(srv, ws)
	}
}

func requestIDToContext(ctx context.Context) context.Context {
	requestID := GetRequestIDFromMetadata(ctx)

	grpc.SetHeader(ctx, metadata.Pairs(requestIDHeader, requestID))
	ctx = corectx.RequestIDToContext(ctx, requestID)

	return ctx
}

func Logger(log *corezaplogger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		ctx = loggerToContext(ctx, log, info.FullMethod)

		return handler(ctx, req)
	}
}

func LoggerStream(log *corezaplogger.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := loggerToContext(ss.Context(), log, info.FullMethod)

		ws := &wrappedStream{
			ServerStream: ss,
			ctx:          ctx,
		}

		return handler(srv, ws)
	}
}

func loggerToContext(ctx context.Context, log *corezaplogger.Logger, fullMethod string) context.Context {
	requestID, err := corectx.RequestIDFromContext(ctx)
	if err != nil {
		requestID = "unknown"
	}

	l := log.With(
		zap.String("request_id", requestID),
		zap.String("method", fullMethod),
	)

	ctx = corezaplogger.ToContext(ctx, l)

	return ctx
}

type JWTService interface {
	ParseAccessToken(access string) (corejwt.Claims, error)
}

func JWToken(jwtService JWTService, appliedToMethods ...string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		ctx, err = claimsToContext(ctx, jwtService, appliedToMethods, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func JWTokenStream(jwtService JWTService, appliedToMethods ...string) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx, err := claimsToContext(ss.Context(), jwtService, appliedToMethods, info.FullMethod)
		if err != nil {
			return err
		}

		ws := &wrappedStream{
			ServerStream: ss,
			ctx:          ctx,
		}

		return handler(srv, ws)
	}
}

func claimsToContext(ctx context.Context, jwtService JWTService, appliedToMethods []string, fullMethod string) (context.Context, error) {
	required := slices.Contains(appliedToMethods, fullMethod)
	if !required {
		return ctx, nil
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

	return ctx, nil
}

func Role(methodRoles map[string][]string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if err := verifyRole(ctx, methodRoles, info.FullMethod); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func RoleStream(methodRoles map[string][]string) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := verifyRole(ss.Context(), methodRoles, info.FullMethod); err != nil {
			return err
		}

		return handler(srv, ss)
	}
}

func verifyRole(ctx context.Context, methodRoles map[string][]string, fullMethod string) error {
	requiredRoles, ok := methodRoles[fullMethod]
	if !ok {
		return nil
	}

	claims, err := corejwt.ClaimsFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get claims: %w", coreerrors.ErrUnauthorized)
	}

	role := claims.Role

	if !slices.Contains(requiredRoles, role) {
		return fmt.Errorf("role %s not allowed: %w", role, coreerrors.ErrForbidden)
	}

	return nil
}

func ErrorHandler() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			return nil, handleError(ctx, err, info.FullMethod)
		}

		return resp, nil
	}
}

func ErrorHandlerStream() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(srv, ss)
		if err != nil {
			return handleError(ss.Context(), err, info.FullMethod)
		}

		return nil
	}
}

func handleError(ctx context.Context, err error, fullMethod string) error {
	var (
		code    codes.Code
		logFunc func(string, ...zap.Field)
	)

	log := corezaplogger.FromContext(ctx)

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

	logFunc(
		"request failed",
		zap.Uint32("code", uint32(code)),
		zap.String("method", fullMethod),
		zap.Error(err),
	)

	return status.Error(code, err.Error())
}

package coregrpc

import (
	"context"
	"fmt"
	"net"
	"time"

	corezaplogger "github.com/egotk/golang-advert-app/internal/core/logger/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	Server *grpc.Server
	config Config
	log    *corezaplogger.Logger
}

func NewServer(
	config Config,
	log *corezaplogger.Logger,
	uInterceptors []grpc.UnaryServerInterceptor,
	sInterceptors []grpc.StreamServerInterceptor,
) *Server {

	return &Server{
		Server: grpc.NewServer(
			grpc.ChainUnaryInterceptor(uInterceptors...),
			grpc.ChainStreamInterceptor(sInterceptors...),
		),
		config: config,
		log:    log,
	}
}

func (s *Server) GetRegistrar() grpc.ServiceRegistrar {
	return s.Server
}

func (s *Server) Start(ctx context.Context) error {
	if s.config.ShouldUseReflection {
		reflection.Register(s.Server)
	}

	lis, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		return fmt.Errorf("listen gRPC: %w", err)
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Warn("start gRPC server", zap.String("addr", s.config.Addr))

		if err := s.Server.Serve(lis); err != nil {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("serve gRPC: %w", err)
		}
	case <-ctx.Done():
		stopped := make(chan struct{})

		go func() {
			s.Server.GracefulStop()
			close(stopped)
		}()

		select {
		case <-stopped:
		case <-time.After(s.config.ShutdownTimeout):
			s.Server.Stop()
		}
	}

	return nil
}

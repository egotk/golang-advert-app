package corehttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	corelogger "github.com/egotk/golang-advert-app/internal/core/logger"
)

type Server struct {
	mux        *http.ServeMux
	config     Config
	log        corelogger.Logger
	middleware []Middleware
}

func NewServer(
	config Config,
	log corelogger.Logger,
	middleware ...Middleware,
) *Server {
	return &Server{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (s *Server) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		s.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router),
		)
	}
}

func (s *Server) Run(ctx context.Context) error {
	mux := ChainMiddleware(s.mux, s.middleware...)

	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Warn("start HTTP server", corelogger.String("addr", s.config.Addr))

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve HTTP: %w", err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.config.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}
	}

	return nil
}

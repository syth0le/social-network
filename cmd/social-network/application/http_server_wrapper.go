package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"social-network/cmd/social-network/configuration"
	"social-network/internal/utils"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ServerOption struct {
	Port int
	Mux  *chi.Mux
}

type httpServerOption struct {
	adminServerOption   *ServerOption
	publicServersOption []*ServerOption
}

type HTTPServerOption func(*httpServerOption)

func WithAdminServer(cfg configuration.ServerConfig) HTTPServerOption {
	return func(opts *httpServerOption) {
		opts.adminServerOption = &ServerOption{Port: cfg.Port}
	}
}

func WithPublicServer(cfg configuration.ServerConfig, mux *chi.Mux) HTTPServerOption {
	return func(opts *httpServerOption) {
		opts.publicServersOption = append(opts.publicServersOption, &ServerOption{Port: cfg.Port, Mux: mux})
	}
}

type HTTPServerWrapper struct {
	logger  *zap.Logger
	servers []*http.Server
}

func NewHTTPServerWrapper(logger *zap.Logger, opts ...HTTPServerOption) *HTTPServerWrapper {
	options := &httpServerOption{
		adminServerOption:   nil,
		publicServersOption: nil,
	}

	for _, o := range opts {
		o(options)
	}

	var servers []*http.Server

	if options.adminServerOption != nil {
		// todo make admin server
		servers = append(servers, newNetHTTPServer(logger, options.adminServerOption.Port))
	}

	for _, option := range options.publicServersOption {
		servers = append(servers, newNetHTTPServer(logger, option.Port))
	}

	return &HTTPServerWrapper{
		logger:  logger,
		servers: servers,
	}
}

func (h *HTTPServerWrapper) Run() []func() error {
	runFunc := func(server *http.Server) error {
		h.logger.Sugar().Infof("run http server on addr: %s", server.Addr)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server listen and serve: %w", err)
		}
		return nil
	}

	response := make([]func() error, 0, len(h.servers))
	for _, server := range h.servers {
		server := server
		response = append(response, func() error {
			return runFunc(server)
		})
	}
	return response
}

func (h *HTTPServerWrapper) GracefulStop() []func() error {
	gracefulFunc := func(server *http.Server) error {
		err := server.Shutdown(context.Background())
		if err != nil {
			return fmt.Errorf("server shutdown: %w", err)
		}
		return nil
	}

	response := make([]func() error, 0, len(h.servers))
	for _, server := range h.servers {
		server := server
		response = append(response, func() error {
			return gracefulFunc(server)
		})
	}
	return response
}

func newNetHTTPServer(logger *zap.Logger, port int) *http.Server {
	// TODO: admin server wrapper
	mux := chi.NewMux()
	mux.Use(utils.LoggerMiddleware(logger))
	mux.Get("/ping", pingHandler())

	return &http.Server{
		Addr:     fmt.Sprintf(":%d", port),
		Handler:  mux,
		ErrorLog: log.New(os.Stderr, "", 0), // TODO
	}
}

func pingHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK\n"))
	}
}

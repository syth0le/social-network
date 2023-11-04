package application

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type httpServerOption struct {
	adminServerOption   *struct{}
	publicServersOption []*struct{}
}

type HTTPServerOption func(*httpServerOption)

func WithAdminServer(cfg struct{}) HTTPServerOption {
	return func(opts *httpServerOption) {
		opts.adminServerOption = &struct{}{}
	}
}

func WithPublicServer(cfg struct{}, mux *chi.Mux) HTTPServerOption {
	// TODO: use mux
	return func(opts *httpServerOption) {
		opts.publicServersOption = append(opts.publicServersOption, &struct{}{})
	}
}

type HTTPServerWrapper struct {
	logger  zap.Logger
	servers []*http.Server
}

func NewHTTPServerWrapper(logger zap.Logger, opts ...HTTPServerOption) *HTTPServerWrapper {
	options := &httpServerOption{
		adminServerOption:   nil,
		publicServersOption: nil,
	}

	for _, o := range opts {
		o(options)
	}

	servers := []*http.Server{}

	if options.adminServerOption != nil {
		// todo make admin server
		// servers = append(servers, nil)
	}

	for _, option := range options.publicServersOption {
		servers = append(servers, newNetHTTPServer(logger, option.port))
	}

	return &HTTPServerWrapper{
		logger:  logger,
		servers: servers,
	}
}

func (h *HTTPServerWrapper) Run() error {
	for _, server := range h.servers {
		err := server.ListenAndServe()
		if err != nil {
			return fmt.Errorf("stop http server, addr: %s", server.Addr, err)
		}
	}
	return nil
}

func newNetHTTPServer(logger zap.Logger, port int) *http.Server {
	// TODO: admin server wrapper
	mux := chi.NewMux()
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

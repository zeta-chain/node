package zetaclient

import (
	"context"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// HTTPServer provide http endpoint for tss server
type HTTPServer struct {
	logger    zerolog.Logger
	s         *http.Server
	pendingTx uint64
}

// NewHTTPServer should only listen to the loopback
func NewHTTPServer(tssAddr string) *HTTPServer {
	hs := &HTTPServer{
		logger: log.With().Str("module", "http").Logger(),
	}
	s := &http.Server{
		Addr:    tssAddr,
		Handler: hs.Handlers(),
	}
	hs.s = s
	return hs
}

// NewHandler registers the API routes and returns a new HTTP handler
func (t *HTTPServer) Handlers() http.Handler {
	router := mux.NewRouter()
	router.Handle("/ping", http.HandlerFunc(t.pingHandler)).Methods(http.MethodGet)
	router.Handle("/metrics", promhttp.Handler())
	router.Use(logMiddleware())
	return router
}

func (t *HTTPServer) Start() error {
	if t.s == nil {
		return errors.New("invalid http server instance")
	}
	if err := t.s.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			return fmt.Errorf("fail to start http server: %w", err)
		}
	}

	return nil
}

func logMiddleware() mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debug().
				Str("route", r.URL.Path).
				Str("port", r.URL.Port()).
				Str("method", r.Method).
				Msg("HTTP request received")

			handler.ServeHTTP(w, r)
		})
	}
}

func (t *HTTPServer) Stop() error {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := t.s.Shutdown(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to shutdown the HTTP server gracefully")
	}
	return err
}

func (t *HTTPServer) pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

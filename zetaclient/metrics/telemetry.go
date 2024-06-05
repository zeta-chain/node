package metrics

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// TelemetryServer provide http endpoint for Tss server
type TelemetryServer struct {
	logger         zerolog.Logger
	s              *http.Server
	p2pid          string
	mu             sync.Mutex
	ipAddress      string
	HotKeyBurnRate *BurnRate
}

// NewTelemetryServer should only listen to the loopback
func NewTelemetryServer() *TelemetryServer {
	hs := &TelemetryServer{
		logger:         log.With().Str("module", "http").Logger(),
		HotKeyBurnRate: NewBurnRate(100),
	}
	s := &http.Server{
		Addr:              ":8123",
		Handler:           hs.Handlers(),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}
	hs.s = s
	return hs
}

// SetP2PID sets p2pid
func (t *TelemetryServer) SetP2PID(p2pid string) {
	t.mu.Lock()
	t.p2pid = p2pid
	t.mu.Unlock()
}

// GetP2PID gets p2pid
func (t *TelemetryServer) GetP2PID() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.p2pid
}

// SetIPAddress sets p2pid
func (t *TelemetryServer) SetIPAddress(ip string) {
	t.mu.Lock()
	t.ipAddress = ip
	t.mu.Unlock()
}

// GetIPAddress gets p2pid
func (t *TelemetryServer) GetIPAddress() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.ipAddress
}

// AddFeeEntry adds fee entry
func (t *TelemetryServer) AddFeeEntry(block int64, amount int64) {
	t.mu.Lock()
	err := t.HotKeyBurnRate.AddFee(amount, block)
	if err != nil {
		log.Error().Err(err).Msg("failed to update hotkey burn rate")
	}
	t.mu.Unlock()
}

// Handlers registers the API routes and returns a new HTTP handler
func (t *TelemetryServer) Handlers() http.Handler {
	router := mux.NewRouter()
	router.Handle("/ping", http.HandlerFunc(t.pingHandler)).Methods(http.MethodGet)
	router.Handle("/p2p", http.HandlerFunc(t.p2pHandler)).Methods(http.MethodGet)
	router.Handle("/ip", http.HandlerFunc(t.ipHandler)).Methods(http.MethodGet)
	router.Handle("/hotkeyburnrate", http.HandlerFunc(t.hotKeyFeeBurnRate)).Methods(http.MethodGet)

	router.Use(logMiddleware())

	return router
}

// Start starts the telemetry server
func (t *TelemetryServer) Start() error {
	if t.s == nil {
		return errors.New("invalid http server instance")
	}
	if err := t.s.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("fail to start http server: %w", err)
		}
	}

	return nil
}

// Stop stops the telemetry server
func (t *TelemetryServer) Stop() error {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := t.s.Shutdown(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to shutdown the HTTP server gracefully")
	}
	return err
}

// pingHandler returns a 200 OK response
func (t *TelemetryServer) pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// p2pHandler returns the p2p id
func (t *TelemetryServer) p2pHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	t.mu.Lock()
	defer t.mu.Unlock()
	fmt.Fprintf(w, "%s", t.p2pid)
}

// ipHandler returns the ip address
func (t *TelemetryServer) ipHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	t.mu.Lock()
	defer t.mu.Unlock()
	fmt.Fprintf(w, "%s", t.ipAddress)
}

// hotKeyFeeBurnRate returns the hot key fee burn rate
func (t *TelemetryServer) hotKeyFeeBurnRate(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	t.mu.Lock()
	defer t.mu.Unlock()
	fmt.Fprintf(w, "%v", t.HotKeyBurnRate.GetBurnRate())
}

// logMiddleware logs the incoming HTTP request
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

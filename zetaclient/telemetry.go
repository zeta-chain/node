package zetaclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/zeta-chain/zetacore/common"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/zetaclient/types"
)

// TelemetryServer provide http endpoint for Tss server
type TelemetryServer struct {
	logger                 zerolog.Logger
	s                      *http.Server
	p2pid                  string
	lastScannedBlockNumber map[int64]int64 // chainid => block number
	lastCoreBlockNumber    int64
	mu                     sync.Mutex
	lastStartTimestamp     time.Time
	status                 types.Status
}

// NewTelemetryServer should only listen to the loopback
func NewTelemetryServer() *TelemetryServer {
	hs := &TelemetryServer{
		logger:                 log.With().Str("module", "http").Logger(),
		lastScannedBlockNumber: make(map[int64]int64),
		lastStartTimestamp:     time.Now(),
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

func (t *TelemetryServer) GetLastStartTimestamp() time.Time {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.lastStartTimestamp
}

// setter/getter for p2pid
func (t *TelemetryServer) SetP2PID(p2pid string) {
	t.mu.Lock()
	t.p2pid = p2pid
	t.mu.Unlock()
}

func (t *TelemetryServer) GetP2PID() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.p2pid
}

// setter for lastScanned block number
func (t *TelemetryServer) SetLastScannedBlockNumber(chainID int64, blockNumber int64) {
	t.mu.Lock()
	t.lastScannedBlockNumber[chainID] = blockNumber
	t.mu.Unlock()
}

func (t *TelemetryServer) GetLastScannedBlockNumber(chainID int64) int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.lastScannedBlockNumber[chainID]
}

func (t *TelemetryServer) SetCoreBlockNumber(blockNumber int64) {
	t.mu.Lock()
	t.lastCoreBlockNumber = blockNumber
	t.mu.Unlock()
}

func (t *TelemetryServer) GetCoreBlockNumber() int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.lastCoreBlockNumber
}

func (t *TelemetryServer) SetNextNonce(nextNonce int) {
	t.mu.Lock()
	t.status.BTCNextNonce = nextNonce
	t.mu.Unlock()
}

func (t *TelemetryServer) GetNextNonce() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.status.BTCNextNonce
}

func (t *TelemetryServer) SetNumberOfUTXOs(numberOfUTXOs int) {
	t.mu.Lock()
	t.status.BTCNumberOfUTXOs = numberOfUTXOs
	t.mu.Unlock()
}

// NewHandler registers the API routes and returns a new HTTP handler
func (t *TelemetryServer) Handlers() http.Handler {
	router := mux.NewRouter()
	router.Handle("/ping", http.HandlerFunc(t.pingHandler)).Methods(http.MethodGet)
	router.Handle("/p2p", http.HandlerFunc(t.p2pHandler)).Methods(http.MethodGet)
	router.Handle("/version", http.HandlerFunc(t.versionHandler)).Methods(http.MethodGet)
	router.Handle("/lastscannedblock", http.HandlerFunc(t.lastScannedBlockHandler)).Methods(http.MethodGet)
	router.Handle("/laststarttimestamp", http.HandlerFunc(t.lastStartTimestampHandler)).Methods(http.MethodGet)
	router.Handle("/lastcoreblock", http.HandlerFunc(t.lastCoreBlockHandler)).Methods(http.MethodGet)
	router.Handle("/status", http.HandlerFunc(t.statusHandler)).Methods(http.MethodGet)
	// router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	// router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	// router.HandleFunc("/debug/pprof/", pprof.Index)
	// router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)

	//router.Handle("/pending", http.HandlerFunc(t.pendingHandler)).Methods(http.MethodGet)
	router.Use(logMiddleware())
	return router
}

func (t *TelemetryServer) Start() error {
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

func (t *TelemetryServer) Stop() error {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := t.s.Shutdown(c)
	if err != nil {
		log.Error().Err(err).Msg("Failed to shutdown the HTTP server gracefully")
	}
	return err
}

func (t *TelemetryServer) pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (t *TelemetryServer) p2pHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	t.mu.Lock()
	defer t.mu.Unlock()
	fmt.Fprintf(w, "%s", t.p2pid)
}

func (t *TelemetryServer) lastScannedBlockHandler(w http.ResponseWriter, _ *http.Request) {
	//w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	t.mu.Lock()
	defer t.mu.Unlock()
	// Convert map to JSON
	jsonBytes, err := json.Marshal(t.lastScannedBlockNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonBytes)
	if err != nil {
		t.logger.Error().Err(err).Msg("Failed to write response")
	}
}

func (t *TelemetryServer) lastCoreBlockHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	t.mu.Lock()
	defer t.mu.Unlock()
	fmt.Fprintf(w, "%d", t.lastCoreBlockNumber)
}

func (t *TelemetryServer) statusHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	t.mu.Lock()
	defer t.mu.Unlock()
	s, _ := json.MarshalIndent(t.status, "", "\t")
	fmt.Fprintf(w, "%s", s)
}

func (t *TelemetryServer) versionHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", common.Version)
}

func (t *TelemetryServer) lastStartTimestampHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	t.mu.Lock()
	defer t.mu.Unlock()
	fmt.Fprintf(w, "%s", t.lastStartTimestamp.Format(time.RFC3339))
}

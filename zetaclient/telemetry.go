package zetaclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeta-chain/zetacore/common"
	"net/http"
	"net/http/pprof"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// HTTPServer provide http endpoint for Tss server
type HTTPServer struct {
	logger                 zerolog.Logger
	s                      *http.Server
	p2pid                  string
	lastScannedBlockNumber map[int64]int64 // chainid => block number
	lastCoreBlockNumber    int64
	mu                     sync.Mutex
	lastStartTimestamp     time.Time
}

// NewHTTPServer should only listen to the loopback
func NewHTTPServer() *HTTPServer {
	hs := &HTTPServer{
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

func (t *HTTPServer) GetLastStartTimestamp() time.Time {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.lastStartTimestamp
}

// setter/getter for p2pid
func (t *HTTPServer) SetP2PID(p2pid string) {
	t.mu.Lock()
	t.p2pid = p2pid
	t.mu.Unlock()
}

func (t *HTTPServer) GetP2PID() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.p2pid
}

// setter for lastScanned block number
func (t *HTTPServer) SetLastScannedBlockNumber(chainId int64, blockNumber int64) {
	t.mu.Lock()
	t.lastScannedBlockNumber[chainId] = blockNumber
	t.mu.Unlock()
}

func (t *HTTPServer) GetLastScannedBlockNumber(chainId int64) int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.lastScannedBlockNumber[chainId]
}

func (t *HTTPServer) SetCoreBlockNumber(blockNumber int64) {
	t.mu.Lock()
	t.lastCoreBlockNumber = blockNumber
	t.mu.Unlock()
}

func (t *HTTPServer) GetCoreBlockNumber() int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.lastCoreBlockNumber
}

// NewHandler registers the API routes and returns a new HTTP handler
func (t *HTTPServer) Handlers() http.Handler {
	router := mux.NewRouter()
	router.Handle("/ping", http.HandlerFunc(t.pingHandler)).Methods(http.MethodGet)
	router.Handle("/p2p", http.HandlerFunc(t.p2pHandler)).Methods(http.MethodGet)
	router.Handle("/version", http.HandlerFunc(t.versionHandler)).Methods(http.MethodGet)
	router.Handle("/lastscannedblock", http.HandlerFunc(t.lastScannedBlockHandler)).Methods(http.MethodGet)
	router.Handle("/laststarttamstamp", http.HandlerFunc(t.lastStartTimestampHandler)).Methods(http.MethodGet)
	router.Handle("/lastcoreblock", http.HandlerFunc(t.lastCoreBlockHandler)).Methods(http.MethodGet)
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)

	//router.Handle("/pending", http.HandlerFunc(t.pendingHandler)).Methods(http.MethodGet)
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

func (t *HTTPServer) p2pHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	t.mu.Lock()
	defer t.mu.Unlock()
	fmt.Fprintf(w, "%s", t.p2pid)
}

func (t *HTTPServer) lastScannedBlockHandler(w http.ResponseWriter, _ *http.Request) {
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
	w.Write(jsonBytes)
}

func (t *HTTPServer) lastCoreBlockHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	t.mu.Lock()
	defer t.mu.Unlock()
	fmt.Fprintf(w, "%d", t.lastCoreBlockNumber)
}

func (t *HTTPServer) versionHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", common.Version)
}

func (t *HTTPServer) lastStartTimestampHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	t.mu.Lock()
	defer t.mu.Unlock()
	fmt.Fprintf(w, "%s", t.lastStartTimestamp.Format(time.RFC3339))
}

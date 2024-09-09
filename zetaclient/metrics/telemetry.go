package metrics

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/zetaclient/types"
)

// TelemetryServer provides http endpoint for Tss server
type TelemetryServer struct {
	logger                 zerolog.Logger
	s                      *http.Server
	p2pid                  string
	lastScannedBlockNumber map[int64]uint64 // chainID => block number
	lastCoreBlockNumber    int64
	mu                     sync.Mutex
	lastStartTimestamp     time.Time
	status                 types.Status
	ipAddress              string
	HotKeyBurnRate         *BurnRate
}

// NewTelemetryServer should only listen to the loopback
func NewTelemetryServer() *TelemetryServer {
	hs := &TelemetryServer{
		logger:                 log.With().Str("module", "http").Logger(),
		lastScannedBlockNumber: make(map[int64]uint64),
		lastStartTimestamp:     time.Now(),
		HotKeyBurnRate:         NewBurnRate(100),
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

// GetLastStartTimestamp returns last start timestamp
func (t *TelemetryServer) GetLastStartTimestamp() time.Time {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.lastStartTimestamp
}

// SetLastScannedBlockNumber last scanned block number for chain in telemetry and metrics
func (t *TelemetryServer) SetLastScannedBlockNumber(chain chains.Chain, blockNumber uint64) {
	t.mu.Lock()
	t.lastScannedBlockNumber[chain.ChainId] = blockNumber
	LastScannedBlockNumber.WithLabelValues(chain.Name).Set(float64(blockNumber))
	t.mu.Unlock()
}

// GetLastScannedBlockNumber returns last scanned block number for chain
func (t *TelemetryServer) GetLastScannedBlockNumber(chainID int64) uint64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.lastScannedBlockNumber[chainID]
}

// SetCoreBlockNumber sets core block number in telemetry and metrics
func (t *TelemetryServer) SetCoreBlockNumber(blockNumber int64) {
	t.mu.Lock()
	t.lastCoreBlockNumber = blockNumber
	LastCoreBlockNumber.Set(float64(blockNumber))
	t.mu.Unlock()
}

// GetCoreBlockNumber returns core block number
func (t *TelemetryServer) GetCoreBlockNumber() int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.lastCoreBlockNumber
}

// SetNumberOfUTXOs sets number of UTXOs in telemetry and metrics
func (t *TelemetryServer) SetNumberOfUTXOs(numberOfUTXOs int) {
	t.mu.Lock()
	t.status.BTCNumberOfUTXOs = numberOfUTXOs
	NumberOfUTXO.Set(float64(numberOfUTXOs))
	t.mu.Unlock()
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
	router.Handle("/version", http.HandlerFunc(t.versionHandler)).Methods(http.MethodGet)
	router.Handle("/lastscannedblock", http.HandlerFunc(t.lastScannedBlockHandler)).Methods(http.MethodGet)
	router.Handle("/laststarttimestamp", http.HandlerFunc(t.lastStartTimestampHandler)).Methods(http.MethodGet)
	router.Handle("/lastcoreblock", http.HandlerFunc(t.lastCoreBlockHandler)).Methods(http.MethodGet)
	router.Handle("/status", http.HandlerFunc(t.statusHandler)).Methods(http.MethodGet)
	router.Handle("/ip", http.HandlerFunc(t.ipHandler)).Methods(http.MethodGet)
	router.Handle("/hotkeyburnrate", http.HandlerFunc(t.hotKeyFeeBurnRate)).Methods(http.MethodGet)

	router.Use(logMiddleware())

	return router
}

// Start starts telemetry server
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

// Stop stops telemetry server
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

func (t *TelemetryServer) lastScannedBlockHandler(w http.ResponseWriter, _ *http.Request) {
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
	s, err := json.MarshalIndent(t.status, "", "\t")
	if err != nil {
		t.logger.Error().Err(err).Msg("Failed to marshal status")
	}
	fmt.Fprintf(w, "%s", s)
}

func (t *TelemetryServer) versionHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", constant.Version)
}

func (t *TelemetryServer) lastStartTimestampHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	t.mu.Lock()
	defer t.mu.Unlock()
	fmt.Fprintf(w, "%s", t.lastStartTimestamp.Format(time.RFC3339))
}

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

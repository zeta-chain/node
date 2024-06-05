package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

// Metrics is a struct that contains the http server for metrics
type Metrics struct {
	s *http.Server
}

// ZetaClientNamespace is the namespace for the metrics
const ZetaClientNamespace = "zetaclient"

var (
	// PendingTxsPerChain is a gauge that contains the number of pending transactions per chain
	PendingTxsPerChain = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "pending_txs_total",
		Help:      "Number of pending transactions per chain",
	}, []string{"chain"})

	// GetFilterLogsPerChain is a counter that contains the number of getLogs per chain
	GetFilterLogsPerChain = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: ZetaClientNamespace,
		Name:      "rpc_getFilterLogs_count",
		Help:      "Count of getLogs per chain",
	}, []string{"chain"})

	// GetBlockByNumberPerChain is a counter that contains the number of getBlockByNumber per chain
	GetBlockByNumberPerChain = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: ZetaClientNamespace,
		Name:      "rpc_getBlockByNumber_count",
		Help:      "Count of getLogs per chain",
	}, []string{"chain"})

	// TssNodeBlamePerPubKey is a counter that contains the number of tss node blame per pubkey
	TssNodeBlamePerPubKey = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: ZetaClientNamespace,
		Name:      "tss_node_blame_count",
		Help:      "Tss node blame counter per pubkey",
	}, []string{"pubkey"})

	// HotKeyBurnRate is a gauge that contains the fee burn rate of the hotkey
	HotKeyBurnRate = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "hotkey_burn_rate",
		Help:      "Fee burn rate of the hotkey",
	})

	// NumberOfUTXO is a gauge that contains the number of UTXOs
	NumberOfUTXO = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "utxo_number",
		Help:      "Number of UTXOs",
	})

	// LastScannedBlockNumber is a gauge that contains the last scanned block number per chain
	LastScannedBlockNumber = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "last_scanned_block_number",
		Help:      "Last scanned block number per chain",
	}, []string{"chain"})

	// LastCoreBlockNumber is a gauge that contains the last core block number
	LastCoreBlockNumber = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "last_core_block_number",
		Help:      "Last core block number",
	})

	// Info is a gauge that contains information about the zetaclient environment
	Info = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "info",
		Help:      "Information about Zetaclient environment",
	}, []string{"version"})

	// LastStartTime is a gauge that contains the start time in Unix time
	LastStartTime = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "last_start_timestamp_seconds",
		Help:      "Start time in Unix time",
	})

	// NumActiveMsgSigns is a gauge that contains the number of concurrent key signs
	NumActiveMsgSigns = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "num_active_message_signs",
		Help:      "Number of concurrent key signs",
	})

	// PercentageOfRateReached is a gauge that contains the percentage of the rate limiter rate reached
	PercentageOfRateReached = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "percentage_of_rate_reached",
		Help:      "Percentage of the rate limiter rate reached",
	})
)

// NewMetrics creates a new Metrics instance
func NewMetrics() (*Metrics, error) {
	handler := promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer,
		promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
			Timeout: 30 * time.Second,
		}),
	)

	s := &http.Server{
		Addr:              ":8886",
		Handler:           handler,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &Metrics{
		s,
	}, nil
}

// Start starts the metrics server
func (m *Metrics) Start() {
	log.Info().Msg("metrics server starting")
	go func() {
		if err := m.s.ListenAndServe(); err != nil {
			log.Error().Err(err).Msg("fail to start metric server")
		}
	}()
}

// Stop stops the metrics server
func (m *Metrics) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	return m.s.Shutdown(ctx)
}

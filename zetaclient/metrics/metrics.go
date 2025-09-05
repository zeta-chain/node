// Package metrics provides metrics functionalities for the zetaclient
package metrics

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"cosmossdk.io/errors"
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

	// GetBlockNumberPerChain is a counter that contains the number of getBlockNumber per chain
	GetBlockNumberPerChain = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: ZetaClientNamespace,
		Name:      "rpc_getBlockNumber_count",
		Help:      "Count of blockNumber per chain",
	}, []string{"chain"})

	// TSSNodeBlamePerPubKey is a counter that contains the number of tss node blame per pubkey
	TSSNodeBlamePerPubKey = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: ZetaClientNamespace,
		Name:      "tss_node_blame_count",
		Help:      "TSS node blame counter per pubkey",
	}, []string{"pubkey"})

	// RelayerKeyBalance is a gauge that contains the relayer key balance of the chain
	RelayerKeyBalance = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "relayer_key_balance",
		Help:      "Relayer key balance of the chain",
	}, []string{"chain"})

	// HotKeyBurnRate is a gauge that contains the fee burn rate of the hotkey
	HotKeyBurnRate = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "hotkey_burn_rate",
		Help:      "Fee burn rate of the hotkey",
	})

	// NumberOfUTXO is a gauge that contains the number of UTXOs
	NumberOfUTXO = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "utxo_number",
		Help:      "Number of UTXOs",
	}, []string{"chain"})

	// LastScannedBlockNumber is a gauge that contains the last scanned block number per chain
	LastScannedBlockNumber = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "last_scanned_block_number",
		Help:      "Last scanned block number per chain",
	}, []string{"chain"})

	// LatestBlockLatency is a gauge that contains the block latency for each observed chain
	LatestBlockLatency = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "latest_block_latency",
		Help:      "Latency of last block for observed chains",
	}, []string{"chain"})

	// LastCoreBlockNumber is a gauge that contains the last core block number
	LastCoreBlockNumber = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "last_core_block_number",
		Help:      "Last core block number",
	})

	// CoreBlockLatency is a gauge that measures the difference between system time and
	// block time from zetacore
	CoreBlockLatency = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "core_block_latency",
		Help:      "Difference between system time and block time from zetacore",
	})

	// CoreBlockLatencySleep is a gauge of the duration we sleep before signing
	CoreBlockLatencySleep = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "core_block_latency_sleep",
		Help:      "The duration we sleep before signing",
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

	// SignLatency is a histogram of of the TSS keysign latency
	SignLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: ZetaClientNamespace,
		Name:      "sign_latency",
		Help:      "Histogram of the TSS keysign latency",
		Buckets:   []float64{1, 7, 15, 30, 60, 120, 240},
	}, []string{"result"})

	// RPCInProgress is a gauge that contains the number of RPCs requests in progress
	RPCInProgress = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "rpc_in_progress",
		Help:      "Number of RPC requests in progress",
	}, []string{"host"})

	// RPCCount is a counter that contains the number of total RPC requests
	RPCCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: ZetaClientNamespace,
			Name:      "rpc_count",
			Help:      "A counter for number of total RPC requests",
		},
		[]string{"host", "code"},
	)

	// RPCLatency is a histogram of the RPC latency
	RPCLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: ZetaClientNamespace,
			Name:      "rpc_duration_seconds",
			Help:      "A histogram of the RPC duration in seconds",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"host"},
	)

	NumConnectedPeers = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "num_connected_peers",
		Help:      "The number of connected peers (authenticated keygen peers)",
	})

	// SchedulerTaskInvocationCounter tracks invocations categorized by status, group, and name
	SchedulerTaskInvocationCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: ZetaClientNamespace,
			Name:      "scheduler_task_invocations_total",
			Help:      "Total number of task invocations",
		},
		[]string{"status", "task_group", "task_name"},
	)

	// SchedulerTaskExecutionDuration measures the execution duration of tasks
	SchedulerTaskExecutionDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: ZetaClientNamespace,
			Name:      "scheduler_task_duration_seconds",
			Help:      "Histogram of task execution duration in seconds",
			Buckets:   []float64{0.05, 0.1, 0.2, 0.3, 0.5, 1, 1.5, 2, 3, 5, 7.5, 10, 15}, // 50ms to 15s
		},
		[]string{"status", "task_group", "task_name"},
	)

	// NumTrackerReporters is a gauge that tracks the number of active tracker reporters
	NumTrackerReporters = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "num_tracker_reporters",
		Help:      "The number of active tracker reporters",
	}, []string{"chain"})

	RPCClientCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: ZetaClientNamespace,
			Name:      "rpc_client_calls_total",
			Help:      "Total number of rpc calls",
		},
		[]string{"status", "client", "method"},
	)

	RPCClientDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: ZetaClientNamespace,
			Name:      "rpc_client_duration_seconds",
			Help:      "Histogram of rpc client calls durations in seconds",
			Buckets:   []float64{0.05, 0.1, 0.2, 0.3, 0.5, 1, 1.5, 2, 3, 5, 7.5, 10, 15}, // 50ms to 15s
		},
		[]string{"client"},
	)
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
func (m *Metrics) Start(_ context.Context) error {
	log.Info().Msg("metrics server starting")

	if err := m.s.ListenAndServe(); err != nil {
		return errors.Wrap(err, "fail to start metric server")
	}

	return nil
}

// Stop stops the metrics server
func (m *Metrics) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.s.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("failed to shutdown metrics server")
	}
}

// GetInstrumentedHTTPClient sets up a http client that emits prometheus metrics
func GetInstrumentedHTTPClient(endpoint string) (*http.Client, error) {
	host := endpoint
	// try to parse as url (so that we do not expose auth uuid in metrics)
	endpointURL, err := url.Parse(endpoint)
	if err == nil {
		host = endpointURL.Host
	}
	labels := prometheus.Labels{"host": host}
	rpcCounterMetric, err := RPCCount.CurryWith(labels)
	if err != nil {
		return nil, err
	}
	rpcLatencyMetric, err := RPCLatency.CurryWith(labels)
	if err != nil {
		return nil, err
	}

	transport := http.DefaultTransport
	transport = promhttp.InstrumentRoundTripperDuration(rpcLatencyMetric, transport)
	transport = promhttp.InstrumentRoundTripperCounter(rpcCounterMetric, transport)
	transport = promhttp.InstrumentRoundTripperInFlight(RPCInProgress.With(labels), transport)

	return &http.Client{
		Transport: transport,
	}, nil
}

// ReportBlockLatency records the latency between the current time
// an the latest block time for a chain as a metric
func ReportBlockLatency(chainName string, latestBlockTime time.Time) {
	elapsedTime := time.Since(latestBlockTime)
	LatestBlockLatency.WithLabelValues(chainName).Set(elapsedTime.Seconds())
}

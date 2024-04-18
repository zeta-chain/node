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

type Metrics struct {
	s *http.Server
}

const ZetaClientNamespace = "zetaclient"

var (
	PendingTxsPerChain = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "pending_txs_total",
		Help:      "Number of pending transactions per chain",
	}, []string{"chain"})

	GetFilterLogsPerChain = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: ZetaClientNamespace,
		Name:      "rpc_getFilterLogs_count",
		Help:      "Count of getLogs per chain",
	}, []string{"chain"})

	GetBlockByNumberPerChain = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: ZetaClientNamespace,
		Name:      "rpc_getBlockByNumber_count",
		Help:      "Count of getLogs per chain",
	}, []string{"chain"})

	TssNodeBlamePerPubKey = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: ZetaClientNamespace,
		Name:      "tss_node_blame_count",
		Help:      "Tss node blame counter per pubkey",
	}, []string{"pubkey"})

	HotKeyBurnRate = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "hotkey_burn_rate",
		Help:      "Fee burn rate of the hotkey",
	})

	NumberOfUTXO = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "utxo_number",
		Help:      "Number of UTXOs",
	})

	LastScannedBlockNumber = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "last_scanned_block_number",
		Help:      "Last scanned block number per chain",
	}, []string{"chain"})

	LastCoreBlockNumber = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "last_core_block_number",
		Help:      "Last core block number",
	})

	Info = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "info",
		Help:      "Information about Zetaclient environment",
	}, []string{"version"})

	LastStartTime = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "last_start_timestamp_seconds",
		Help:      "Start time in Unix time",
	})

	NumActiveMsgSigns = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: ZetaClientNamespace,
		Name:      "num_active_message_signs",
		Help:      "Number of concurrent key signs",
	})
)

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

func (m *Metrics) Start() {
	log.Info().Msg("metrics server starting")
	go func() {
		if err := m.s.ListenAndServe(); err != nil {
			log.Error().Err(err).Msg("fail to start metric server")
		}
	}()
}

func (m *Metrics) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	return m.s.Shutdown(ctx)
}

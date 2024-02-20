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

var (
	PendingTxsPerChain = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "zetaclient",
		Name:      "pending_txs_total",
		Help:      "Number of pending transactions per chain",
	}, []string{"chain"})

	GetFilterLogsPerChain = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "zetaclient",
		Name:      "rpc_getFilterLogs_count",
		Help:      "Count of getLogs per chain",
	}, []string{"chain"})

	GetBlockByNumberPerChain = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "zetaclient",
		Name:      "rpc_getBlockByNumber_count",
		Help:      "Count of getLogs per chain",
	}, []string{"chain"})

	TssNodeBlamePerPubKey = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "zetaclient",
		Name:      "tss_node_blame_count",
		Help:      "Tss node blame counter per pubkey",
	}, []string{"pubkey"})

	HotKeyBurnRate = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "zetaclient",
		Name:      "hotkey_burn_rate",
		Help:      "Fee burn rate of the hotkey",
	})
)

func NewMetrics() (*Metrics, error) {
	server := http.NewServeMux()

	server.Handle("/metrics",
		promhttp.InstrumentMetricHandler(
			prometheus.DefaultRegisterer,
			promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
				Timeout: 30 * time.Second,
			}),
		),
	)

	s := &http.Server{
		Addr:              ":8886",
		Handler:           server,
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

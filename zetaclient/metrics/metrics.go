package metrics

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type Metrics struct {
	s *http.Server
}

var (
	counters = map[string]prometheus.Counter{}

	gauges = map[string]prometheus.Gauge{}
)

func NewMetrics() (*Metrics, error) {
	var pendingSendGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "num_pending_send",
		Help: "current number of Sends from the ZetaCore pendingSend API",
	})

	prometheus.MustRegister(pendingSendGauge)

	pendingSendGauge.Set(234)

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
		Addr:        fmt.Sprintf(":8080"),
		Handler:     server,
		ReadTimeout: 5 * time.Second,
	}

	return &Metrics{
		s,
	}, nil
}

func (m *Metrics) Start() {
	log.Info().Msg("Metrics server starting...")
	go func() {
		if err := m.s.ListenAndServe(); err != nil {
			log.Error().Err(err).Msg("fail to stop metric server")
		}
	}()
}

func (m *Metrics) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	return m.s.Shutdown(ctx)
}

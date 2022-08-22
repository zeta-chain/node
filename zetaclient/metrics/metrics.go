package metrics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

type Metrics struct {
	s *http.Server
}

type MetricName int

const (
	//GAUGE_PENDING_TX MetricName = iota
	//
	//COUNTER_NUM_RPCS
	PENDING_TXS = "pending_txs"
)

var (
	Counters = map[string]prometheus.Counter{}

	Gauges = map[string]prometheus.Gauge{}
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
		Addr:              fmt.Sprintf(":8886"),
		Handler:           server,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &Metrics{
		s,
	}, nil
}

func (m *Metrics) RegisterCounter(name string, help string) error {
	if _, found := Counters[name]; found {
		return fmt.Errorf("counter %s already registered", name)
	}
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: name,
		Help: help,
	})
	prometheus.MustRegister(counter)
	Counters[name] = counter
	return nil
}

func (m *Metrics) RegisterGauge(name string, help string) error {
	if _, found := Gauges[name]; found {
		return fmt.Errorf("gauge %s already registered", name)
	}

	var gauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: name,
		Help: help,
	})
	prometheus.MustRegister(gauge)
	Gauges[name] = gauge
	return nil
}

func (m *Metrics) Start() {
	log.Info().Msg("Metrics server starting...")
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

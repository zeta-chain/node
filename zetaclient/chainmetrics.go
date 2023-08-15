package zetaclient

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
)

const MetricGroup = "zetaclient"

type ChainMetrics struct {
	chain   string
	metrics *metrics.Metrics
}

func NewChainMetrics(chain string, metrics *metrics.Metrics) *ChainMetrics {
	return &ChainMetrics{
		chain,
		metrics,
	}
}

func (m *ChainMetrics) GetPromGauge(name string) (prometheus.Gauge, error) {
	gauge, found := metrics.Gauges[m.buildGroupName(name)]
	if !found {
		return nil, errors.New("gauge not found")
	}
	return gauge, nil
}

func (m *ChainMetrics) RegisterPromGauge(name string, help string) error {
	gaugeName := m.buildGroupName(name)
	return m.metrics.RegisterGauge(gaugeName, help)
}

func (m *ChainMetrics) GetPromCounter(name string) (prometheus.Counter, error) {
	if cnt, found := metrics.Counters[m.buildGroupName(name)]; found {
		return cnt, nil
	}
	return nil, errors.New("counter not found")

}

func (m *ChainMetrics) RegisterPromCounter(name string, help string) error {
	cntName := m.buildGroupName(name)
	return m.metrics.RegisterCounter(cntName, help)
}

func (m *ChainMetrics) buildGroupName(name string) string {
	return MetricGroup + "_" + name + "_" + m.chain
}

package observer

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
)

func (ob *Observer) GetPromGauge(name string) (prometheus.Gauge, error) {
	gauge, found := metrics.Gauges[ob.chain.String()+"_"+name]
	if !found {
		return nil, errors.New("gauge not found")
	}
	return gauge, nil
}

func (ob *Observer) RegisterPromGauge(name string, help string) error {
	gaugeName := ob.chain.String() + "_" + name
	return ob.metrics.RegisterGauge(gaugeName, help)
}

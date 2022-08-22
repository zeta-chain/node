package zetaclient

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
)

func (ob *ChainObserver) GetPromGauge(name string) (prometheus.Gauge, error) {
	if gauge, found := metrics.Gauges[ob.chain.String()+"_"+name]; found {
		return gauge, nil
	} else {
		return nil, errors.New("gauge not found")
	}
}

func (ob *ChainObserver) RegisterPromGauge(name string, help string) error {
	gaugeName := ob.chain.String() + "_" + name
	return ob.metrics.RegisterGauge(gaugeName, help)
}

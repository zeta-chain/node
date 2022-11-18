package zetaclient

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
)

func (ob *EVMChainClient) GetPromGauge(name string) (prometheus.Gauge, error) {
	gauge, found := metrics.Gauges[ob.chain.String()+"_"+name]
	if !found {
		return nil, errors.New("gauge not found")
	}
	return gauge, nil
}

func (ob *EVMChainClient) RegisterPromGauge(name string, help string) error {
	gaugeName := ob.chain.String() + "_" + name
	return ob.metrics.RegisterGauge(gaugeName, help)
}

func (ob *BitcoinChainClient) GetPromGauge(name string) (prometheus.Gauge, error) {
	gauge, found := metrics.Gauges[ob.chain.String()+"_"+name]
	if !found {
		return nil, errors.New("gauge not found")
	}
	return gauge, nil
}

func (ob *BitcoinChainClient) RegisterPromGauge(name string, help string) error {
	gaugeName := ob.chain.String() + "_" + name
	return ob.metrics.RegisterGauge(gaugeName, help)
}

package observer

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/zeta-chain/node/zetaclient/chains/base"
)

// Observer SUI observer
type Observer struct {
	*base.Observer
	client RPC
}

// RPC represents subset of SUI RPC methods.
type RPC interface {
	HealthCheck(ctx context.Context) (time.Time, error)
}

// New Observer constructor.
func New(baseObserver *base.Observer, client RPC) *Observer {
	return &Observer{
		Observer: baseObserver,
		client:   client,
	}
}

// CheckRPCStatus checks the RPC status of the chain.
func (ob *Observer) CheckRPCStatus(ctx context.Context) error {
	blockTime, err := ob.client.HealthCheck(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to check rpc health")
	}

	// It's not a "real" block latency as SUI uses concept of "checkpoints"
	ob.ReportBlockLatency(blockTime)

	return nil
}

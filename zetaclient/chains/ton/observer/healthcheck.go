package observer

import (
	"context"

	"github.com/pkg/errors"

	"github.com/zeta-chain/node/zetaclient/metrics"
)

// CheckRPCStatus checks the status of the TON client and reports the block latency.
func (ob *Observer) CheckRPCStatus(ctx context.Context) error {
	blockTime, err := ob.tonRepo.CheckHealth(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to check TON client health")
	}

	if blockTime == nil {
		return errors.Wrap(err, "invalid block time (internal error)")
	}

	metrics.ReportBlockLatency(ob.Chain().Name, *blockTime)
	return nil
}

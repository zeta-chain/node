package observer

import (
	"context"

	"github.com/pkg/errors"
)

// ObserveGasPrice fetches on-chain gas config and reports it to Zetacore.
func (ob *Observer) ObserveGasPrice(ctx context.Context) error {
	// Gets the latest gas price and block number.
	gasPrice, blockNumber, err := ob.tonRepo.GetGasPrice(ctx)
	if err != nil {
		return err
	}

	// There's no concept of priority fee in TON
	const priorityFee = 0

	_, err = ob.ZetacoreClient().PostVoteGasPrice(ctx, ob.Chain(), gasPrice, priorityFee, blockNumber)
	if err != nil {
		return errors.Wrap(err, "unable to post gas price")
	}

	ob.setLatestGasPrice(gasPrice)

	return nil
}

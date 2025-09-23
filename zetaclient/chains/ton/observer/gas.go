package observer

import (
	"context"

	"github.com/pkg/errors"
)

// // ObserveGasPrice fetches on-chain gas config and reports it to Zetacore.
// func (ob *Observer) ObserveGasPrice(ctx context.Context) error {
// 	// Gets the latest gas price and block number.
// 	gasPrice, blockNumber, err := ob.repo.GetGasPrice(ctx)
// 	if err != nil {
// 		return err
// 	}
//
// 	chain := ob.Chain()
//
// 	err = ob.repo.VoteGasPrice(ctx, chain, gasPrice, blockNumber)
// 	if err != nil {
// 		return err
// 	}
//
// 	ob.setLatestGasPrice(gasPrice) // TODO: should this only happen after?
//
// 	return nil
// }

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
		return errors.Wrap(err, "failed to post gas price")
	}

	ob.setLatestGasPrice(gasPrice)

	return nil
}

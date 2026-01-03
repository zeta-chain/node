package observer

import (
	"context"
)

// ObserveGasPrice fetches on-chain gas information and reports it to zetacore.
func (ob *Observer) ObserveGasPrice(ctx context.Context) error {
	// Gets the latest gas price and block number.
	gasPrice, block, err := ob.tonRepo.GetGasPrice(ctx)
	if err != nil {
		return err
	}

	// There's no concept of priority fee in TON.
	const priorityFee = 0

	var (
		logger     = ob.Logger().Chain
		multiplier = ob.ChainParams().GasPriceMultiplier
	)

	_, err = ob.ZetaRepo().VoteGasPrice(ctx, logger, gasPrice, multiplier, priorityFee, block)
	if err != nil {
		return err
	}

	ob.setLatestGasPrice(gasPrice)

	return nil
}

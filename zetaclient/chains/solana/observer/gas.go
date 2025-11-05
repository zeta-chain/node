package observer

import (
	"context"

	"github.com/gagliardetto/solana-go/rpc"
)

// PostGasPrice posts gas prices to zetacore.
func (ob *Observer) PostGasPrice(ctx context.Context) error {
	priorityFee, err := ob.solanaRepo.GetPriorityFee(ctx)
	if err != nil {
		return err
	}

	slot, err := ob.solanaRepo.GetSlot(ctx, rpc.CommitmentConfirmed)
	if err != nil {
		return err
	}

	// There is no Ethereum-like gas price in Solana, so we only post the priority fee.
	_, err = ob.ZetaRepo().VoteGasPrice(ctx, ob.Logger().Chain, 1, priorityFee, slot)
	return err
}

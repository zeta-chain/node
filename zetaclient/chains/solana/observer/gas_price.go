package observer

import (
	"context"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
)

// WatchGasPrice watches the gas price of the chain and posts it to the zetacore
func (ob *Observer) WatchGasPrice(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			slot, err := ob.solClient.GetSlot(context.Background(), rpc.CommitmentConfirmed)
			if err != nil {
				ob.Logger().GasPrice.Err(err).Msg("GetSlot error")
				continue
			}
			// FIXME: what's the fee rate of compute unit? How to query?
			txhash, err := ob.ZetacoreClient().PostVoteGasPrice(ctx, ob.Chain(), 1, "", slot)
			if err != nil {
				ob.Logger().GasPrice.Err(err).Msg("PostGasPrice error")
				continue
			}
			ob.Logger().GasPrice.Info().Msgf("gas price posted: %s", txhash)
		case <-ob.StopChannel():
			ob.Logger().GasPrice.Info().Msgf("WatchGasPrice stopped for chain %d", ob.Chain().ChainId)
			return
		}
	}
}

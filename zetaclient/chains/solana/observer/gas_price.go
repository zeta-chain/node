package observer

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"

	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
)

// WatchGasPrice watches the gas price of the chain and posts it to the zetacore
func (ob *Observer) WatchGasPrice(ctx context.Context) error {
	// report gas price right away as the ticker takes time to kick in
	err := ob.PostGasPrice(ctx)
	if err != nil {
		ob.Logger().GasPrice.Error().Err(err).Msgf("PostGasPrice error for chain %d", ob.Chain().ChainId)
	}

	// start gas price ticker
	ticker, err := clienttypes.NewDynamicTicker(
		fmt.Sprintf("Solana_WatchGasPrice_%d", ob.Chain().ChainId),
		ob.GetChainParams().GasPriceTicker,
	)
	if err != nil {
		return errors.Wrapf(err, "NewDynamicTicker error")
	}
	ob.Logger().GasPrice.Info().Msgf("WatchGasPrice started for chain %d with interval %d",
		ob.Chain().ChainId, ob.GetChainParams().GasPriceTicker)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C():
			if !ob.GetChainParams().IsSupported {
				continue
			}
			err = ob.PostGasPrice(ctx)
			if err != nil {
				ob.Logger().GasPrice.Error().Err(err).Msgf("PostGasPrice error for chain %d", ob.Chain().ChainId)
			}
			ticker.UpdateInterval(ob.GetChainParams().GasPriceTicker, ob.Logger().GasPrice)
		case <-ob.StopChannel():
			ob.Logger().GasPrice.Info().Msgf("WatchGasPrice stopped for chain %d", ob.Chain().ChainId)
			return nil
		}
	}
}

// PostGasPrice posts gas price to zetacore
func (ob *Observer) PostGasPrice(ctx context.Context) error {
	// get current slot
	slot, err := ob.solClient.GetSlot(context.Background(), rpc.CommitmentConfirmed)
	if err != nil {
		return errors.Wrap(err, "GetSlot error")
	}

	// post gas price to zetacore
	// FIXME: what's the fee rate of compute unit? How to query?
	_, err = ob.ZetacoreClient().PostVoteGasPrice(ctx, ob.Chain(), 1, 0, slot)
	if err != nil {
		return errors.Wrap(err, "PostVoteGasPrice error")
	}

	return nil
}

package observer

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"

	zetamath "github.com/zeta-chain/node/pkg/math"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
)

const (
	// SolanaTransactionFee is the static fee per transaction, 5k lamports.
	SolanaTransactionFee = 5000

	// MicroLamportsPerLamport is the number of micro lamports in a lamport.
	MicroLamportsPerLamport = 1_000_000

	// SolanaDefaultComputeBudget is the default compute budget for a transaction.
	SolanaDefaultComputeBudget = 200_000

	// Solana uses micro lamports (0.000001 lamports) as the smallest unit of gas price.
	// The gas fee formula 'gasFee = gasPrice * gasLimit' won't fit Solana in the ZRC20 SOL contract.
	// We could use lamports as the unit of gas price and 10K CU as the smallest unit of compute units.
	// SolanaDefaultGasPrice10KCUs is the default gas price (in lamports) per 10K compute units.
	SolanaDefaultGasPrice10KCUs = 100

	// SolanaDefaultGasLimit is the default compute units (in 10K CU) for a transaction.
	SolanaDefaultGasLimit10KCU = 50
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
	slot, err := ob.solClient.GetSlot(ctx, rpc.CommitmentConfirmed)
	if err != nil {
		return errors.Wrap(err, "GetSlot error")
	}

	// query recent priority fees
	recentFees, err := ob.solClient.GetRecentPrioritizationFees(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "GetRecentPrioritizationFees error")
	}

	// locate median priority fee
	priorityFees := make([]uint64, len(recentFees))
	for i, fee := range recentFees {
		if fee.PrioritizationFee > 0 {
			priorityFees[i] = fee.PrioritizationFee
		}
	}
	// the priority fee is in increments of 0.000001 lamports (micro lamports)
	medianFee := zetamath.SliceMedianValue(priorityFees, true)

	// there is no Ethereum-like gas price in Solana, we only post priority fee for now
	_, err = ob.ZetacoreClient().PostVoteGasPrice(ctx, ob.Chain(), 1, medianFee, slot)
	if err != nil {
		return errors.Wrapf(err, "PostVoteGasPrice error for chain %d", ob.Chain().ChainId)
	}

	return nil
}

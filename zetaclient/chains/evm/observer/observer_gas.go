package observer

import (
	"context"
	"fmt"
	"math/big"

	"github.com/pkg/errors"

	clienttypes "github.com/zeta-chain/node/zetaclient/types"
)

// WatchGasPrice watches evm chain for gas prices and post to zetacore
// TODO(revamp): move inner logic to a separate function
func (ob *Observer) WatchGasPrice(ctx context.Context) error {
	// report gas price right away as the ticker takes time to kick in
	err := ob.PostGasPrice(ctx)
	if err != nil {
		ob.Logger().GasPrice.Error().Err(err).Msgf("PostGasPrice error for chain %d", ob.Chain().ChainId)
	}

	// start gas price ticker
	ticker, err := clienttypes.NewDynamicTicker(
		fmt.Sprintf("EVM_WatchGasPrice_%d", ob.Chain().ChainId),
		ob.GetChainParams().GasPriceTicker,
	)
	if err != nil {
		ob.Logger().GasPrice.Error().Err(err).Msg("NewDynamicTicker error")
		return err
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
			ob.Logger().GasPrice.Info().Msg("WatchGasPrice stopped")
			return nil
		}
	}
}

// PostGasPrice posts gas price to zetacore.
func (ob *Observer) PostGasPrice(ctx context.Context) error {
	// GAS PRICE
	gasPrice, err := ob.evmClient.SuggestGasPrice(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to suggest gas price")
	}

	// PRIORITY FEE (EIP-1559)
	priorityFee, err := ob.determinePriorityFee(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to determine priority fee")
	}

	blockNum, err := ob.evmClient.BlockNumber(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get block number")
	}

	_, err = ob.
		ZetacoreClient().
		PostVoteGasPrice(ctx, ob.Chain(), gasPrice.Uint64(), priorityFee.Uint64(), blockNum)

	if err != nil {
		return errors.Wrap(err, "unable to post vote for gas price")
	}

	return nil
}

// determinePriorityFee determines the chain priority fee.
// Returns zero for non EIP-1559 (London fork) chains.
func (ob *Observer) determinePriorityFee(ctx context.Context) (*big.Int, error) {
	supported, err := ob.supportsPriorityFee(ctx)
	switch {
	case err != nil:
		return nil, err
	case !supported:
		// noop
		return big.NewInt(0), nil
	}

	fee, err := ob.evmClient.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to suggest gas tip cap")
	}

	return fee, nil
}

// supportsPriorityFee checks if the chain supports EIP-1559 (London fork).
// uses cache so actual RPC call is made only once.
func (ob *Observer) supportsPriorityFee(ctx context.Context) (bool, error) {
	// noop
	if ob.priorityFeeConfig.checked {
		return ob.priorityFeeConfig.supported, nil
	}

	baseFee, err := ob.getChainBaseFee(ctx)
	if err != nil {
		return false, errors.Wrap(err, "unable to get base fee")
	}

	ob.Logger().GasPrice.Info().
		Int64("observer.chain_id", ob.Chain().ChainId).
		Str("observer.base_fee", baseFee.String()).
		Msg("Fetched base fee for chain")

	// EIP-1559 is supported if base fee is not zero.
	// Not that, for example, BSC supports EIP-1559 but base fee is zero.
	isSupported := baseFee != nil

	ob.Mu().Lock()
	defer ob.Mu().Unlock()

	ob.priorityFeeConfig.checked = true
	ob.priorityFeeConfig.supported = isSupported

	return isSupported, nil
}

// getChainBaseFee fetches baseFee from latest block's header.
func (ob *Observer) getChainBaseFee(ctx context.Context) (*big.Int, error) {
	// get latest block
	header, err := ob.evmClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get latest block header")
	}

	if header.BaseFee == nil {
		return big.NewInt(0), nil
	}

	return header.BaseFee, nil
}

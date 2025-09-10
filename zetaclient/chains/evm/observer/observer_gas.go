package observer

import (
	"context"
	"math/big"

	"github.com/pkg/errors"
)

// PostGasPrice posts gas price to zetacore.
func (ob *Observer) PostGasPrice(ctx context.Context) error {
	// GAS PRICE
	gasPrice, err := ob.evmClient.SuggestGasPrice(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to suggest gas price")
	}

	// The zetaclients now only build legacy tx rather than EIP-1559 tx.
	// Hardcode priority fee to zero to avoid gas price bump failure in the zetacore:
	// https://github.com/zeta-chain/node/blob/release%2Fv30/x/crosschain/keeper/abci.go#L182
	// https://github.com/zeta-chain/node/issues/3221
	priorityFee := uint64(0)

	// PRIORITY FEE (EIP-1559)
	// priorityFee, err := ob.determinePriorityFee(ctx)
	// if err != nil {
	// 	return errors.Wrap(err, "unable to determine priority fee")
	// }

	blockNum, err := ob.evmClient.BlockNumber(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get block number")
	}

	_, err = ob.
		ZetacoreClient().
		PostVoteGasPrice(ctx, ob.Chain(), gasPrice.Uint64(), priorityFee, blockNum)

	if err != nil {
		return errors.Wrap(err, "unable to post vote for gas price")
	}

	return nil
}

// DeterminePriorityFee determines the chain priority fee.
// Returns zero for non EIP-1559 (London fork) chains.
func (ob *Observer) DeterminePriorityFee(ctx context.Context) (*big.Int, error) {
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
		Str("base_fee", baseFee.String()).
		Msg("fetched base fee")

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

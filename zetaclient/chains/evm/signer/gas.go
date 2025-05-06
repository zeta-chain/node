package signer

import (
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// maxGasLimit is the maximum gas limit cap for EVM chain outbound to prevent excessive gas
const maxGasLimit = 2_500_000

// contractCallMinGasLimit is the minimum gas limit for contract calls to prevent intrinsic low gas limit error
const contractCallMinGasLimit = 100_000

// Gas represents gas parameters for EVM transactions.
//
// This is pretty interesting because all EVM chains now support EIP-1559, but some chains do it in a specific way
// https://eips.ethereum.org/EIPS/eip-1559
// https://www.blocknative.com/blog/eip-1559-fees
// https://github.com/bnb-chain/BEPs/blob/master/BEPs/BEP226.md (tl;dr: baseFee is always zero)
//
// However, this doesn't affect tx creation nor broadcasting
type Gas struct {
	Limit uint64

	// This is a "total" gasPrice per 1 unit of gas.
	// GasPrice for pre EIP-1559 transactions or maxFeePerGas for EIP-1559.
	Price *big.Int

	// PriorityFee a fee paid directly to validators for EIP-1559.
	PriorityFee *big.Int
}

func (g Gas) validate() error {
	switch {
	case g.Limit == 0:
		return errors.New("gas limit is zero")
	case g.Price == nil:
		return errors.New("max fee per unit is nil")
	case g.PriorityFee == nil:
		return errors.New("priority fee per unit is nil")
	case g.Price.Cmp(g.PriorityFee) == -1:
		return fmt.Errorf(
			"max fee per unit (%d) is less than priority fee per unit (%d)",
			g.Price.Int64(),
			g.PriorityFee.Int64(),
		)
	default:
		return nil
	}
}

// isLegacy determines whether the gas is meant for LegacyTx{} (pre EIP-1559)
// or DynamicFeeTx{} (post EIP-1559).
//
// Returns true if priority fee is <= 0.
//
//nolint:unused // https://github.com/zeta-chain/node/issues/3221
func (g Gas) isLegacy() bool {
	return g.PriorityFee.Sign() < 1
}

func gasFromCCTX(cctx *types.CrossChainTx, logger zerolog.Logger) (Gas, error) {
	var (
		params = cctx.GetCurrentOutboundParam()
		limit  = params.CallOptions.GasLimit
	)

	isWithdrawWithNoCall := !cctx.IsCurrentOutboundRevert() && cctx.InboundParams.IsCrossChainCall
	isRevertWithNoCall := cctx.IsCurrentOutboundRevert() && cctx.RevertOptions.CallOnRevert

	if limit > maxGasLimit {
		limit = maxGasLimit
		logger.Warn().
			Uint64("cctx.initial_gas_limit", params.CallOptions.GasLimit).
			Uint64("cctx.gas_limit", limit).
			Msgf("Gas limit is too high; Setting to the maximum (%d)", maxGasLimit)
	} else if limit < contractCallMinGasLimit && (params.CoinType != coin.CoinType_Gas || isWithdrawWithNoCall || isRevertWithNoCall) {
		// currently a gas limit that is too low will not make the transaction fail but just completely block the network because the outbound can't be broadcasted
		// this check ensure that the gas limit is high enough for the outbound to be broadcasted
		// a minimal gas limit is set to the tx, this check is skipped above in the case where a simple gas withdraw is performed or a gas deposit revert without revert call
		// TODO: gas limit that is too low should now block outbound at all
		// https://github.com/zeta-chain/node/issues/3725
		limit = contractCallMinGasLimit
		logger.Warn().
			Uint64("cctx.initial_gas_limit", params.CallOptions.GasLimit).
			Uint64("cctx.gas_limit", limit).
			Msgf("Gas limit is too low for contract call; Setting to the minimum (%d)", contractCallMinGasLimit)
	}

	gasPrice, err := bigIntFromString(params.GasPrice)
	if err != nil {
		return Gas{}, errors.Wrap(err, "unable to parse gasPrice")
	}

	priorityFee, err := bigIntFromString(params.GasPriorityFee)
	switch {
	case err != nil:
		return Gas{}, errors.Wrap(err, "unable to parse priorityFee")
	case gasPrice.Cmp(priorityFee) == -1:
		logger.Warn().
			Str("cctx.initial_priority_fee", priorityFee.String()).
			Str("cctx.forced_priority_fee", gasPrice.String()).
			Msg("gasPrice is less than priorityFee, setting priorityFee = gasPrice")

		// this should in theory never happen, but this reported bug might be a cause: https://github.com/zeta-chain/node/issues/2954
		// in this case we lower the priorityFee to the gasPrice to ensure the transaction is valid
		// the only potential issue is the transaction might not cover the baseFee
		// the gas stability pool mechanism help to mitigate this issue
		priorityFee = big.NewInt(0).Set(gasPrice)
	}

	return Gas{
		Limit:       limit,
		Price:       gasPrice,
		PriorityFee: priorityFee,
	}, nil
}

func bigIntFromString(s string) (*big.Int, error) {
	if s == "" || s == "0" {
		return big.NewInt(0), nil
	}

	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("unable to parse %q as big.Int", s)
	}

	if v.Sign() == -1 {
		return nil, fmt.Errorf("big.Int is negative: %d", v.Int64())
	}

	return v, nil
}

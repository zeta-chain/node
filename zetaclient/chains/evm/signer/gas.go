package signer

import (
	"math/big"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

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

	// MaxFeePerUnit absolute maximum we're willing to pay per unit of gas to get tx included in a block.
	MaxFeePerUnit *big.Int

	// PriorityFeePerUnit optional fee paid directly to validators.
	PriorityFeePerUnit *big.Int
}

func (g Gas) validate() error {
	if g.Limit == 0 {
		return errors.New("gas limit is zero")
	}

	if g.MaxFeePerUnit == nil {
		return errors.New("max fee per unit is nil")
	}

	if g.PriorityFeePerUnit == nil {
		return errors.New("priority fee per unit is nil")
	}

	return nil
}

// isLegacy determines whether the gas is meant for LegacyTx{} (pre EIP-1559)
// or DynamicFeeTx{} (post EIP-1559).
//
// Returns true if priority fee is <= 0.
func (g Gas) isLegacy() bool {
	return g.PriorityFeePerUnit.Sign() < 1
}

// makeGasFromCCTX creates Gas struct based from CCTX and priorityFee.
func makeGasFromCCTX(cctx *types.CrossChainTx, priorityFee *big.Int, logger zerolog.Logger) (Gas, error) {
	if priorityFee == nil {
		return Gas{}, errors.New("priorityFee is nil")
	}

	var (
		outboundParams = cctx.GetCurrentOutboundParam()
		limit          = outboundParams.GasLimit
	)

	switch {
	case limit < MinGasLimit:
		limit = MinGasLimit
		logger.Warn().
			Uint64("cctx.initial_gas_limit", outboundParams.GasLimit).
			Uint64("cctx.gas_limit", limit).
			Msgf("Gas limit is too low. Setting to the minimum (%d)", MinGasLimit)
	case limit > MaxGasLimit:
		limit = MaxGasLimit
		logger.Warn().
			Uint64("cctx.initial_gas_limit", outboundParams.GasLimit).
			Uint64("cctx.gas_limit", limit).
			Msgf("Gas limit is too high; Setting to the maximum (%d)", MaxGasLimit)
	}

	maxFee, ok := new(big.Int).SetString(outboundParams.GasPrice, 10)
	if !ok {
		return Gas{}, errors.New("unable to parse gasPrice from " + outboundParams.GasPrice)
	}

	// is maxFee < priorityFee
	if maxFee.Cmp(priorityFee) == -1 {
		return Gas{}, errors.Errorf("maxFee (%d) is less than priorityFee (%d)", maxFee.Int64(), priorityFee.Int64())
	}

	return Gas{
		Limit:              limit,
		MaxFeePerUnit:      maxFee,
		PriorityFeePerUnit: priorityFee,
	}, nil
}

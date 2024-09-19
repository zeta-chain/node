package signer

import (
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/x/crosschain/types"
)

const (
	minGasLimit = 100_000
	maxGasLimit = 1_000_000
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
func (g Gas) isLegacy() bool {
	return g.PriorityFee.Sign() < 1
}

func gasFromCCTX(cctx *types.CrossChainTx, logger zerolog.Logger) (Gas, error) {
	var (
		params = cctx.GetCurrentOutboundParam()
		limit  = params.GasLimit
	)

	switch {
	case limit < minGasLimit:
		limit = minGasLimit
		logger.Warn().
			Uint64("cctx.initial_gas_limit", params.GasLimit).
			Uint64("cctx.gas_limit", limit).
			Msgf("Gas limit is too low. Setting to the minimum (%d)", minGasLimit)
	case limit > maxGasLimit:
		limit = maxGasLimit
		logger.Warn().
			Uint64("cctx.initial_gas_limit", params.GasLimit).
			Uint64("cctx.gas_limit", limit).
			Msgf("Gas limit is too high; Setting to the maximum (%d)", maxGasLimit)
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
		return Gas{}, fmt.Errorf("gasPrice (%d) is less than priorityFee (%d)", gasPrice.Int64(), priorityFee.Int64())
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

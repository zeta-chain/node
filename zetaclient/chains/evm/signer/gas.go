package signer

import (
	"math/big"

	"github.com/pkg/errors"
)

// Gas represents gas parameters for EVM transactions.
//
// This is pretty interesting because all EVM chains now support EIP-1559, but some chains do it in a specific way
// https://eips.ethereum.org/EIPS/eip-1559
// https://www.blocknative.com/blog/eip-1559-fees
// https://github.com/bnb-chain/BEPs/blob/master/BEPs/BEP226.md (tl;dr: baseFee is always zero)
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
	return g.PriorityFeePerUnit.Sign() <= 1
}

package base

import (
	cosmoserrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	cctypes "github.com/zeta-chain/node/x/crosschain/types"
)

// OutboundAmountUint64 validates an amount can fit into uint64 and returns it.
func OutboundAmountUint64(amount sdkmath.Uint) (uint64, error) {
	if amount.BigInt().BitLen() > 64 {
		return 0, cosmoserrors.Wrap(cctypes.ErrInvalidWithdrawalAmount, "amount exceeds uint64 range")
	}

	return amount.Uint64(), nil
}

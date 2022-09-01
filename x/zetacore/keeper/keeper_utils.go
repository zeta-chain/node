package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) isAuthorized(ctx sdk.Context, address string) bool {
	validators := k.StakingKeeper.GetAllValidators(ctx)
	return IsBondedValidator(address, validators)
}

func (k Keeper) hasSuperMajorityValidators(ctx sdk.Context, signers []string) bool {
	numSigners := len(signers)
	validators := k.StakingKeeper.GetAllValidators(ctx)
	numValidValidators := 0
	for _, v := range validators {
		if v.IsBonded() {
			numValidValidators++
		}
	}
	threshold := numValidValidators * 2 / 3
	if threshold < 2 {
		threshold = 2
	}
	if numValidValidators == 1 {
		threshold = 1
	}
	return numSigners == threshold
}

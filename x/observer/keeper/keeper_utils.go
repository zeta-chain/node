package keeper

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) IsValidator(ctx sdk.Context, creator string) bool {
	validators := k.stakingKeeper.GetAllValidators(ctx)
	isValidator := false
	creatorAddr, _ := sdk.AccAddressFromBech32(creator)
	for _, v := range validators {
		valAddr, err := sdk.ValAddressFromBech32(v.OperatorAddress)
		if err != nil {
			continue
		}
		if v.IsBonded() && bytes.Compare(creatorAddr.Bytes(), valAddr.Bytes()) == 0 && v.Jailed == false {
			isValidator = true
			break
		}
	}
	return isValidator
}

package keeper

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func IsBondedValidator(creator string, validators []stakingtypes.Validator) bool {
	creatorAddr, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return false
	}
	for _, v := range validators {
		valAddr, err := sdk.ValAddressFromBech32(v.OperatorAddress)
		if err != nil {
			continue
		}
		//TODO: How about Jailed?
		if v.IsBonded() && bytes.Compare(creatorAddr.Bytes(), valAddr.Bytes()) == 0 {
			return true
		}
	}
	return false
}

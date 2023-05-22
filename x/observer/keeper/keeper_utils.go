package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func (k Keeper) IsValidator(ctx sdk.Context, creator string) error {
	valAddress, err := types.GetOperatorAddressFromAccAddress(creator)
	if err != nil {
		return err
	}
	validator, found := k.stakingKeeper.GetValidator(ctx, valAddress)
	if !found {
		return types.ErrNotValidator
	}

	if validator.Jailed == true || validator.IsBonded() == false {
		return types.ErrValidatorStatus
	}
	return nil

}

func (k Keeper) CheckObserverDelegation(ctx sdk.Context, accAddress string, chain *common.Chain) error {
	selfdelAddr, _ := sdk.AccAddressFromBech32(accAddress)
	valAddress, err := types.GetOperatorAddressFromAccAddress(accAddress)
	if err != nil {
		return err
	}
	validator, found := k.stakingKeeper.GetValidator(ctx, valAddress)
	if !found {
		return types.ErrNotValidator
	}

	delegation, found := k.stakingKeeper.GetDelegation(ctx, selfdelAddr, valAddress)
	if !found {
		return types.ErrSelfDelegation
	}
	obsParams := k.GetParams(ctx).GetParamsForChain(chain)
	if !obsParams.IsSupported {
		return errors.Wrap(types.ErrSupportedChains, fmt.Sprintf("Chain not suported %s ", chain.String()))
	}

	tokens := validator.TokensFromShares(delegation.Shares)
	if tokens.LT(obsParams.MinObserverDelegation) {
		return types.ErrCheckObserverDelegation
	}
	return nil
}

package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
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

func (k Keeper) CheckObserverDelegation(ctx sdk.Context, accAddress string, chain *types.Chain, observationType types.ObservationType) error {
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
		return types.ErrSelfDelgation
	}
	obsParams, found := k.GetParams(ctx).GetParamsForChainAndType(chain, observationType)
	if !found {
		return errors.Wrap(types.ErrSupportedChains, fmt.Sprintf("Params for chain and type do not exists %s , %s", chain.String(), observationType.String()))
	}

	tokens := validator.TokensFromShares(delegation.Shares)
	if tokens.LT(obsParams.MinObserverDelegation) {
		return types.ErrCheckObserverDelegation
	}
	return nil
}

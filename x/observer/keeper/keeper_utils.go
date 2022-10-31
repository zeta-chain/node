package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func (k Keeper) IsValidator(ctx sdk.Context, operatorAddress string) error {
	valAddress, err := sdk.ValAddressFromBech32(operatorAddress)
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

func (k Keeper) CheckObserverDelegation(ctx sdk.Context, msg *types.MsgAddObserver) error {
	delAddr, _ := sdk.AccAddressFromBech32(msg.Creator)
	valAddress, err := sdk.ValAddressFromBech32(msg.ObserverOperator)
	if err != nil {
		return err
	}
	validator, found := k.stakingKeeper.GetValidator(ctx, valAddress)
	if !found {
		return types.ErrNotValidator
	}
	delegation, found := k.stakingKeeper.GetDelegation(ctx, delAddr, valAddress)
	if !found {
		return types.ErrSelfDelgation
	}
	obsParams, found := k.GetParams(ctx).GetParamsForChainAndType(msg.ObserverChain, msg.ObservationType)
	if !found {
		return errors.Wrap(types.ErrSupportedChains, fmt.Sprintf("Params for chain and type do not exists %s , %s", msg.ObserverChain.String(), msg.ObservationType.String()))
	}

	tokens := validator.TokensFromShares(delegation.Shares)
	if tokens.LTE(obsParams.MinObserverDelegation) {
		return types.ErrCheckObserverDelegation
	}
	return nil
}

package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func (k msgServer) UpdateObserver(goCtx context.Context, msg *types.MsgUpdateObserver) (*types.MsgUpdateObserverResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	chains := k.GetParams(ctx).GetSupportedChains()
	for _, chain := range chains {
		if !k.IsObserverPresentInMappers(ctx, msg.OldObserverAddress, chain) {
			return nil, errorsmod.Wrap(types.ErrNotAuthorized, fmt.Sprintf("Observer address is not authorized for chain : %s", chain.String()))
		}
	}

	err := k.IsValidator(ctx, msg.NewObserverAddress)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrUpdateObserver, err.Error())
	}

	ok, err := k.CheckUpdateReason(ctx, msg)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrUpdateObserver, err.Error())
	}
	if !ok {
		return nil, errorsmod.Wrap(types.ErrUpdateObserver, fmt.Sprintf("Unable to update observer with update reason : %s", msg.UpdateReason))
	}

	// Update all mappers so that ballots can be created for the new observer address
	k.UpdateObserverAddress(ctx, msg.OldObserverAddress, msg.NewObserverAddress)

	// Update the node account with the new operator address
	nodeAccount, found := k.GetNodeAccount(ctx, msg.OldObserverAddress)
	if !found {
		return nil, errorsmod.Wrap(types.ErrNodeAccountNotFound, fmt.Sprintf("Observer node account not found : %s", msg.Creator))
	}
	newNodeAccount := nodeAccount
	newNodeAccount.Operator = msg.NewObserverAddress

	// Remove an old node account, so that number of node accounts remains the same as the number of observers in the system
	k.RemoveNodeAccount(ctx, msg.OldObserverAddress)
	k.SetNodeAccount(ctx, newNodeAccount)

	// Check LastBlockObserver count just to be safe
	observerMappers := k.GetAllObserverMappers(ctx)
	totalObserverCountCurrentBlock := uint64(0)
	for _, mapper := range observerMappers {
		totalObserverCountCurrentBlock += uint64(len(mapper.ObserverList))
	}
	lastBlockCount, found := k.GetLastObserverCount(ctx)
	if !found {
		return nil, errorsmod.Wrap(types.ErrLastObserverCountNotFound, fmt.Sprintf("Observer count not found"))
	}
	if lastBlockCount.Count != totalObserverCountCurrentBlock {
		return nil, errorsmod.Wrap(types.ErrUpdateObserver, fmt.Sprintf("Observer count mismatch"))
	}
	return &types.MsgUpdateObserverResponse{}, nil
}

func (k Keeper) CheckUpdateReason(ctx sdk.Context, msg *types.MsgUpdateObserver) (bool, error) {
	switch msg.UpdateReason {
	case types.ObserverUpdateReason_Tombstoned:
		{
			if msg.Creator != msg.OldObserverAddress {
				return false, errorsmod.Wrap(types.ErrUpdateObserver, fmt.Sprintf("Creator address and old observer address need to be same for updating tombstoned observer"))
			}
			return k.IsOperatorTombstoned(ctx, msg.Creator)
		}
	case types.ObserverUpdateReason_AdminUpdate:
		{
			if msg.Creator != k.GetParams(ctx).GetAdminPolicyAccount(types.Policy_Type_group2) {
				return false, types.ErrNotAuthorizedPolicy
			}
			return true, nil
		}
	}
	return false, nil
}

func (k Keeper) UpdateObserverAddress(ctx sdk.Context, oldObserverAddress, newObserverAddress string) {
	observerMappers := k.GetAllObserverMappers(ctx)
	for _, om := range observerMappers {
		UpdateObserverList(om.ObserverList, oldObserverAddress, newObserverAddress)
		k.SetObserverMapper(ctx, om)
	}
}

func UpdateObserverList(list []string, oldObserverAddresss, newObserverAddress string) {
	for i, observer := range list {
		if observer == oldObserverAddresss {
			list[i] = newObserverAddress
		}
	}
}

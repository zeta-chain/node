package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/observer/types"
)

// UpdateObserver handles updating an observer address
// Authorized: admin policy (admin update), old observer address (if the
// reason is that the observer was tombstoned).
func (k msgServer) UpdateObserver(
	goCtx context.Context,
	msg *types.MsgUpdateObserver,
) (*types.MsgUpdateObserverResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ok, err := k.CheckUpdateReason(ctx, msg)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrUpdateObserver, err.Error())
	}
	if !ok {
		return nil, errorsmod.Wrapf(
			types.ErrUpdateObserver,
			"Unable to update observer with update reason : %s", msg.UpdateReason)
	}

	// We do not use CheckObserverCanVote here because we want to allow tombstoned observers to be updated
	if !k.IsAddressPartOfObserverSet(ctx, msg.OldObserverAddress) {
		return nil, errorsmod.Wrapf(
			types.ErrNotObserver,
			"Observer address is not authorized : %s", msg.OldObserverAddress)
	}

	// The New address should be a validator, not jailed and bonded
	err = k.IsValidator(ctx, msg.NewObserverAddress)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrUpdateObserver, err.Error())
	}

	// Update all mappers so that ballots can be created for the new observer address
	err = k.UpdateObserverAddress(ctx, msg.OldObserverAddress, msg.NewObserverAddress)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrUpdateObserver, err.Error())
	}

	// Update the node account with the new operator address
	nodeAccount, found := k.GetNodeAccount(ctx, msg.OldObserverAddress)
	if !found {
		return nil, errorsmod.Wrapf(
			types.ErrNodeAccountNotFound,
			"Observer node account not found : %s", msg.OldObserverAddress)
	}
	newNodeAccount := nodeAccount
	newNodeAccount.Operator = msg.NewObserverAddress

	// Remove an old node account, so that number of node accounts remains the same as the number of observers in the system
	k.RemoveNodeAccount(ctx, msg.OldObserverAddress)
	k.SetNodeAccount(ctx, newNodeAccount)

	// Check LastBlockObserver count just to be safe
	observerSet, found := k.GetObserverSet(ctx)
	if !found {
		return nil, errorsmod.Wrap(types.ErrObserverSetNotFound, "Observer set not found")
	}
	totalObserverCountCurrentBlock := observerSet.LenUint()
	lastBlockCount, found := k.GetLastObserverCount(ctx)
	if !found {
		return nil, errorsmod.Wrap(types.ErrLastObserverCountNotFound, "Observer count not found")
	}
	if lastBlockCount.Count != totalObserverCountCurrentBlock {
		return nil, errorsmod.Wrapf(
			types.ErrUpdateObserver,
			"Observer count mismatch current block: %d , last block: %d",
			totalObserverCountCurrentBlock,
			lastBlockCount.Count,
		)
	}
	return &types.MsgUpdateObserverResponse{}, nil
}

func (k Keeper) CheckUpdateReason(ctx sdk.Context, msg *types.MsgUpdateObserver) (bool, error) {
	switch msg.UpdateReason {
	case types.ObserverUpdateReason_Tombstoned:
		{
			if msg.Creator != msg.OldObserverAddress {
				return false, errorsmod.Wrap(
					types.ErrUpdateObserver,
					"Creator address and old observer address need to be same for updating tombstoned observer",
				)
			}
			return k.IsOperatorTombstoned(ctx, msg.Creator)
		}
	case types.ObserverUpdateReason_AdminUpdate:
		{
			// Operational policy is required to update an observer for admin update
			err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
			if err != nil {
				return false, errorsmod.Wrap(authoritytypes.ErrUnauthorized, err.Error())
			}
			return true, nil
		}
	}
	return false, nil
}

func UpdateObserverList(list []string, oldObserverAddresss, newObserverAddress string) {
	for i, observer := range list {
		if observer == oldObserverAddresss {
			list[i] = newObserverAddress
		}
	}
}

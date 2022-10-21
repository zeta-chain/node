package keeper

import (
	"context"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

func (k msgServer) AddToOutTxTracker(goCtx context.Context, msg *types.MsgAddToOutTxTracker) (*types.MsgAddToOutTxTrackerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validators := k.StakingKeeper.GetAllValidators(ctx)
	if !IsBondedValidator(msg.Creator, validators) && msg.Creator != types.AdminKey {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("signer %s is not a bonded validator", msg.Creator))
	}
	nonceString := strconv.Itoa(int(msg.Nonce))
	index := fmt.Sprintf("%s-%s", msg.Chain, nonceString)
	tracker, found := k.GetOutTxTracker(ctx, index)
	hash := types.TxHashList{
		TxHash: msg.TxHash,
		Signer: msg.Creator,
	}
	if !found {
		k.SetOutTxTracker(ctx, types.OutTxTracker{
			Index:    index,
			Chain:    msg.Chain,
			Nonce:    nonceString,
			HashList: []*types.TxHashList{&hash},
		})
		return &types.MsgAddToOutTxTrackerResponse{}, nil
	}
	var isDup = false
	for _, hash := range tracker.HashList {
		if strings.EqualFold(hash.TxHash, msg.TxHash) {
			isDup = true
			break
		}
	}
	if !isDup {
		tracker.HashList = append(tracker.HashList, &hash)
		k.SetOutTxTracker(ctx, tracker)
	}
	return &types.MsgAddToOutTxTrackerResponse{}, nil
}

func (k msgServer) RemoveFromOutTxTracker(goCtx context.Context, msg *types.MsgRemoveFromOutTxTracker) (*types.MsgRemoveFromOutTxTrackerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validators := k.StakingKeeper.GetAllValidators(ctx)
	if !IsBondedValidator(msg.Creator, validators) && msg.Creator != types.AdminKey {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("signer %s is not a bonded validator", msg.Creator))
	}
	nonceString := strconv.Itoa(int(msg.Nonce))
	index := fmt.Sprintf("%s-%s", msg.Chain, nonceString)

	k.RemoveOutTxTracker(ctx, index)
	return &types.MsgRemoveFromOutTxTrackerResponse{}, nil
}

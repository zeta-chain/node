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

func (k msgServer) AddToWatchList(goCtx context.Context, msg *types.MsgAddToWatchList) (*types.MsgAddToWatchListResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validators := k.StakingKeeper.GetAllValidators(ctx)
	if !IsBondedValidator(msg.Creator, validators) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("signer %s is not a bonded validator", msg.Creator))
	}
	nonceString := strconv.Itoa(int(msg.Nonce))
	index := fmt.Sprintf("%s/%s", msg.Chain, nonceString)
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
		return &types.MsgAddToWatchListResponse{}, nil
	}
	var isDup = false
	for _, hash := range tracker.HashList {
		if strings.EqualFold(hash.TxHash, msg.TxHash) {
			isDup = true
		}
	}
	if !isDup {
		tracker.HashList = append(tracker.HashList, &hash)
		k.SetOutTxTracker(ctx, tracker)
	}
	return &types.MsgAddToWatchListResponse{}, nil
}

package keeper

import (
	"context"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

func (k msgServer) RemoveFromWatchList(goCtx context.Context, msg *types.MsgRemoveFromWatchList) (*types.MsgRemoveFromWatchListResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validators := k.StakingKeeper.GetAllValidators(ctx)
	if !IsBondedValidator(msg.Creator, validators) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("signer %s is not a bonded validator", msg.Creator))
	}
	nonceString := strconv.Itoa(int(msg.Nonce))
	index := fmt.Sprintf("%s/%s", msg.Chain, nonceString)

	k.RemoveOutTxTracker(ctx, index)
	return &types.MsgRemoveFromWatchListResponse{}, nil
}

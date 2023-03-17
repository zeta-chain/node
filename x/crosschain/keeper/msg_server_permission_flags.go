package keeper

import (
	"context"
	types2 "github.com/zeta-chain/zetacore/x/observer/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func (k msgServer) UpdatePermissionFlags(goCtx context.Context, msg *types.MsgUpdatePermissionFlags) (*types.MsgUpdatePermissionFlagsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.zetaObserverKeeper.GetParams(ctx).GetAdminPolicyAccount(types2.Policy_Type_stop_inbound_cctx) {
		return &types.MsgUpdatePermissionFlagsResponse{}, types.ErrNotAuthorized.Wrap("creator does not have enough permissions to set this flag")
	}
	// Check if the value exists
	flags, isFound := k.GetPermissionFlags(ctx)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "not set")
	}
	flags.IsInboundEnabled = msg.IsInboundEnabled
	k.SetPermissionFlags(ctx, flags)

	return &types.MsgUpdatePermissionFlagsResponse{}, nil
}

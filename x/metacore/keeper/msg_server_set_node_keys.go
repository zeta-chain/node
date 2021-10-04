package keeper

import (
	"context"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SetNodeKeys(goCtx context.Context, msg *types.MsgSetNodeKeys) (*types.MsgSetNodeKeysResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("msg creator %s not valid", msg.Creator))
	}
	na, found := k.GetNodeAccount(ctx, msg.Creator)
	if !found {
		na = types.NodeAccount{
			Creator:     msg.Creator,
			Index:       msg.Creator,
			NodeAddress: addr,
			PubkeySet:   msg.PubkeySet,
			NodeStatus:  types.NodeStatus_Unknown,
		}
	} else {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("msg creator %s already has a node account", msg.Creator))
	}

	k.SetNodeAccount(ctx, na)

	return &types.MsgSetNodeKeysResponse{}, nil
}

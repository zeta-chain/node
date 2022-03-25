package keeper

import (
	"context"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

func (k msgServer) SetNodeKeys(goCtx context.Context, msg *types.MsgSetNodeKeys) (*types.MsgSetNodeKeysResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("msg creator %s not valid", msg.Creator))
	}
	_, found := k.GetNodeAccount(ctx, msg.Creator)
	if !found {
		na := types.NodeAccount{
			Creator:     msg.Creator,
			Index:       msg.Creator,
			NodeAddress: addr,
			PubkeySet:   msg.PubkeySet,
			NodeStatus:  types.NodeStatus_Unknown,
		}
		k.SetNodeAccount(ctx, na)
	} else {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("msg creator %s already has a node account", msg.Creator))
	}

	return &types.MsgSetNodeKeysResponse{}, nil
}

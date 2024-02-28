package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/authority/types"
)

// UpdatePolicies updates policies
func (k msgServer) UpdatePolicies(goCtx context.Context, msg *types.MsgUpdatePolicies) (*types.MsgUpdatePoliciesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: check if authorized

	k.SetPolicies(ctx, msg.Policies)

	return &types.MsgUpdatePoliciesResponse{}, nil
}

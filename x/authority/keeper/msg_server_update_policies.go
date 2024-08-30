package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/zeta-chain/node/x/authority/types"
)

// UpdatePolicies updates policies
func (k msgServer) UpdatePolicies(
	goCtx context.Context,
	msg *types.MsgUpdatePolicies,
) (*types.MsgUpdatePoliciesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check called by governance
	if k.govAddr.String() != msg.Creator {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority, expected %s, got %s",
			k.govAddr.String(),
			msg.Creator,
		)
	}

	// set policies
	k.SetPolicies(ctx, msg.Policies)

	return &types.MsgUpdatePoliciesResponse{}, nil
}

package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/observer/types"
)

// UpdateV2ZetaFlows updates the V2 ZETA gateway flows flag.
func (k msgServer) UpdateV2ZetaFlows(
	goCtx context.Context,
	msg *types.MsgUpdateV2ZetaFlows,
) (*types.MsgUpdateV2ZetaFlowsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, errors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// check if the value exists,
	// if not, set the default value for the flags
	flags, isFound := k.GetCrosschainFlags(ctx)
	if !isFound {
		flags = *types.DefaultCrosschainFlags()
	}

	// update V2 ZETA flows flag
	flags.IsV2ZetaEnabled = msg.IsV2ZetaEnabled

	k.SetCrosschainFlags(ctx, flags)

	return &types.MsgUpdateV2ZetaFlowsResponse{}, nil
}

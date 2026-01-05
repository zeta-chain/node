package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/observer/types"
)

// UpdateV2ZetaFlows updates the IsV2ZetaEnabled flag.
// The flag is updated by the policy account with the groupOperational policy type.
func (k msgServer) UpdateV2ZetaFlows(
	goCtx context.Context,
	msg *types.MsgUpdateV2ZetaFlows,
) (*types.MsgUpdateV2ZetaFlowsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, errors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// check if the value exists,
	// if not, set the default value for the Inbound and Outbound flags only
	flags, isFound := k.GetCrosschainFlags(ctx)
	if !isFound {
		flags = *types.DefaultCrosschainFlags()
		flags.GasPriceIncreaseFlags = nil
	}

	flags.IsV2ZetaEnabled = msg.IsV2ZetaEnabled

	k.SetCrosschainFlags(ctx, flags)

	return &types.MsgUpdateV2ZetaFlowsResponse{}, nil
}

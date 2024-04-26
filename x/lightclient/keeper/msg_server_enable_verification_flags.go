package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

// UpdateVerificationFlags updates the light client verification flags.
// This disables/enables blocks verification of the light client for the specified chain.
// Emergency group can disable flags, it requires operational group if at least one flag is being enabled
func (k msgServer) EnableVerificationFlags(goCtx context.Context, msg *types.MsgEnableVerificationFlags) (
	*types.MsgEnableVerificationFlagsResponse,
	error,
) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	if !k.GetAuthorityKeeper().IsAuthorized(ctx, msg.Creator, authoritytypes.PolicyType_groupOperational) {
		return nil, authoritytypes.ErrUnauthorized
	}

	for _, chainID := range msg.ChainIdList {
		// set the verification flags
		k.SetVerificationFlags(ctx, types.VerificationFlags{
			ChainId: chainID,
			Enabled: true,
		})
	}

	return &types.MsgEnableVerificationFlagsResponse{}, nil
}

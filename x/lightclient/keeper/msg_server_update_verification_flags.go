package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

// UpdateVerificationFlags updates the crosschain related flags.
// Emergency group can disable flags while operation group can enable/disable
func (k msgServer) UpdateVerificationFlags(goCtx context.Context, msg *types.MsgUpdateVerificationFlags) (
	*types.MsgUpdateVerificationFlagsResponse,
	error,
) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	requiredGroup := authoritytypes.PolicyType_groupEmergency
	if msg.VerificationFlags.EthTypeChainEnabled || msg.VerificationFlags.BtcTypeChainEnabled {
		requiredGroup = authoritytypes.PolicyType_groupOperational
	}

	// check permission
	if !k.GetAuthorityKeeper().IsAuthorized(ctx, msg.Creator, requiredGroup) {
		return &types.MsgUpdateVerificationFlagsResponse{}, authoritytypes.ErrUnauthorized
	}

	k.SetVerificationFlags(ctx, msg.VerificationFlags)

	return &types.MsgUpdateVerificationFlagsResponse{}, nil
}

package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

// DisableVerificationFlags disables the verification flags for the given chain IDs
// Disabled chains do not allow the submissions of block headers or using it to verify the correctness of proofs
func (k msgServer) DisableVerificationFlags(goCtx context.Context, msg *types.MsgDisableVerificationFlags) (*types.MsgDisableVerificationFlagsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	if !k.GetAuthorityKeeper().IsAuthorized(ctx, msg.Creator, authoritytypes.PolicyType_groupEmergency) {
		return nil, authoritytypes.ErrUnauthorized
	}

	for _, chainID := range msg.ChainIdList {
		// set the verification flags to false to disable verification
		k.SetVerificationFlags(ctx, types.VerificationFlags{
			ChainId: chainID,
			Enabled: false,
		})
	}

	return &types.MsgDisableVerificationFlagsResponse{}, nil
}

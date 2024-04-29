package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

// EnableVerificationFlags enables the verification flags for the given chain IDs
// Enabled chains allow the submissions of block headers and using it to verify the correctness of proofs
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

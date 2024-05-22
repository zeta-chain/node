package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

// EnableHeaderVerification enables the verification flags for the given chain IDs
// Enabled chains allow the submissions of block headers and using it to verify the correctness of proofs
func (k msgServer) EnableHeaderVerification(goCtx context.Context, msg *types.MsgEnableHeaderVerification) (
	*types.MsgEnableHeaderVerificationResponse,
	error,
) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	if !k.GetAuthorityKeeper().IsAuthorized(ctx, msg.Creator, authoritytypes.PolicyType_groupOperational) {
		return nil, authoritytypes.ErrUnauthorized
	}

	bhv, found := k.GetBlockHeaderVerification(ctx)
	if !found {
		bhv = types.BlockHeaderVerification{}
	}

	for _, chainID := range msg.ChainIdList {
		bhv.EnableChain(chainID)
	}

	k.SetBlockHeaderVerification(ctx, bhv)
	return &types.MsgEnableHeaderVerificationResponse{}, nil
}

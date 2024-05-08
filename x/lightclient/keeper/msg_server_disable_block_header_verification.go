package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

// DisableHeaderVerification disables the verification flags for the given chain IDs
// Disabled chains do not allow the submissions of block headers or using it to verify the correctness of proofs
func (k msgServer) DisableHeaderVerification(goCtx context.Context, msg *types.MsgDisableHeaderVerification) (*types.MsgDisableHeaderVerificationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	if !k.GetAuthorityKeeper().IsAuthorized(ctx, msg) {
		return nil, authoritytypes.ErrUnauthorized
	}

	bhv, found := k.GetBlockHeaderVerification(ctx)
	if !found {
		bhv = types.BlockHeaderVerification{}
	}

	for _, chainID := range msg.ChainIdList {
		bhv.DisableChain(chainID)
	}

	k.SetBlockHeaderVerification(ctx, bhv)

	return &types.MsgDisableHeaderVerificationResponse{}, nil
}

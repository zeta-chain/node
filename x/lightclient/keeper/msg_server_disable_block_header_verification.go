package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

// DisableHeaderVerification disables the verification flags for the given chain IDs
// Disabled chains do not allow the submissions of block headers or using it to verify the correctness of proofs
func (k msgServer) DisableHeaderVerification(goCtx context.Context, msg *types.MsgDisableHeaderVerification) (*types.MsgDisableHeaderVerificationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, errorsmod.Wrap(authoritytypes.ErrUnauthorized, err.Error())
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

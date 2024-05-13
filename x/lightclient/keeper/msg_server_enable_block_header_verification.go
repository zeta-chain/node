package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
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
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, errorsmod.Wrap(authoritytypes.ErrUnauthorized, err.Error())
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

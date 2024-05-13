package keeper

import (
	"context"

	cosmoserrors "cosmossdk.io/errors"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// UpdateKeygen updates the block height of the keygen and sets the status to
// "pending keygen".
//
// Authorized: admin policy group 1.
func (k msgServer) UpdateKeygen(goCtx context.Context, msg *types.MsgUpdateKeygen) (*types.MsgUpdateKeygenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	keygen, found := k.GetKeygen(ctx)
	if !found {
		return nil, types.ErrKeygenNotFound
	}
	if msg.Block <= (ctx.BlockHeight() + 10) {
		return nil, types.ErrKeygenBlockTooLow
	}

	nodeAccountList := k.GetAllNodeAccount(ctx)
	granteePubKeys := make([]string, len(nodeAccountList))
	for i, nodeAccount := range nodeAccountList {
		granteePubKeys[i] = nodeAccount.GranteePubkey.Secp256k1.String()
	}

	// update keygen
	keygen.GranteePubkeys = granteePubKeys
	keygen.BlockNumber = msg.Block
	keygen.Status = types.KeygenStatus_PendingKeygen
	k.SetKeygen(ctx, keygen)

	EmitEventKeyGenBlockUpdated(ctx, &keygen)

	return &types.MsgUpdateKeygenResponse{}, nil
}

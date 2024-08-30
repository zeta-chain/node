package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/observer/types"
)

// ResetChainNonces handles resetting chain nonces
func (k msgServer) ResetChainNonces(
	goCtx context.Context,
	msg *types.MsgResetChainNonces,
) (*types.MsgResetChainNoncesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, errors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	tss, found := k.GetTSS(ctx)
	if !found {
		return nil, types.ErrTssNotFound
	}

	// set chain nonces
	chainNonce := types.ChainNonces{
		ChainId: msg.ChainId,
		// #nosec G115 always positive
		Nonce: uint64(msg.ChainNonceHigh),
		// #nosec G115 always positive
		FinalizedHeight: uint64(ctx.BlockHeight()),
	}
	k.SetChainNonces(ctx, chainNonce)

	// set pending nonces
	p := types.PendingNonces{
		NonceLow:  msg.ChainNonceLow,
		NonceHigh: msg.ChainNonceHigh,
		Tss:       tss.TssPubkey,
		ChainId:   msg.ChainId,
	}
	k.SetPendingNonces(ctx, p)

	return &types.MsgResetChainNoncesResponse{}, nil
}

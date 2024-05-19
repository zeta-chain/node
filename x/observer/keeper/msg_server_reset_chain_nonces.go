package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/pkg/chains"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// ResetChainNonces handles resetting chain nonces
func (k msgServer) ResetChainNonces(goCtx context.Context, msg *types.MsgResetChainNonces) (*types.MsgResetChainNoncesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.GetAuthorityKeeper().IsAuthorized(ctx, msg.Creator, authoritytypes.PolicyType_groupOperational) {
		return &types.MsgResetChainNoncesResponse{}, authoritytypes.ErrUnauthorized
	}

	tss, found := k.GetTSS(ctx)
	if !found {
		return nil, types.ErrTssNotFound
	}

	chain := chains.GetChainFromChainID(msg.ChainId)
	if chain == nil {
		return nil, types.ErrSupportedChains
	}

	// set chain nonces
	chainNonce := types.ChainNonces{
		Index:   chain.ChainName.String(),
		ChainId: chain.ChainId,
		// #nosec G701 always positive
		Nonce: uint64(msg.ChainNonceHigh),
		// #nosec G701 always positive
		FinalizedHeight: uint64(ctx.BlockHeight()),
	}
	k.SetChainNonces(ctx, chainNonce)

	// set pending nonces
	p := types.PendingNonces{
		NonceLow:  msg.ChainNonceLow,
		NonceHigh: msg.ChainNonceHigh,
		ChainId:   chain.ChainId,
		Tss:       tss.TssPubkey,
	}
	k.SetPendingNonces(ctx, p)

	return &types.MsgResetChainNoncesResponse{}, nil
}

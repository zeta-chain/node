package keeper

import (
	"context"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// ResetChainNonces handles resetting chain nonces
// Authorized: policy group admin
func (k msgServer) ResetChainNonces(goCtx context.Context, msg *types.MsgResetChainNonces) (*types.MsgResetChainNoncesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	if !k.GetAuthorityKeeper().IsAuthorized(ctx, msg.Creator, authoritytypes.PolicyType_groupAdmin) {
		return &types.MsgResetChainNoncesResponse{}, types.ErrNotAuthorizedPolicy
	}

	tss, found := k.GetTSS(ctx)
	if !found {
		return nil, types.ErrTssNotFound
	}

	chain := common.GetChainFromChainID(msg.ChainId)
	if chain == nil {
		return nil, types.ErrSupportedChains
	}

	// reset chain nonces
	chainNonce := types.ChainNonces{
		Index:   chain.ChainName.String(),
		ChainId: chain.ChainId,
		Nonce:   msg.ChainNonceHigh,
		// #nosec G701 always positive
		FinalizedHeight: uint64(ctx.BlockHeight()),
	}
	k.SetChainNonces(ctx, chainNonce)

	// reset pending nonces
	p := types.PendingNonces{
		NonceLow:  int64(msg.ChainNonceLow),
		NonceHigh: int64(msg.ChainNonceHigh),
		ChainId:   chain.ChainId,
		Tss:       tss.TssPubkey,
	}
	k.SetPendingNonces(ctx, p)

	return &types.MsgResetChainNoncesResponse{}, nil
}

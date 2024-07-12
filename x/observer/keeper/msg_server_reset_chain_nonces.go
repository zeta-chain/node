package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/pkg/chains"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
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

	chain, found := chains.GetChainFromChainID(msg.ChainId, k.GetAuthorityKeeper().GetAdditionalChainList(ctx))
	if !found {
		return nil, types.ErrSupportedChains
	}

	// set chain nonces
	chainNonce := types.ChainNonces{
		Index:   chain.ChainName.String(),
		ChainId: chain.ChainId,
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
		ChainId:   chain.ChainId,
		Tss:       tss.TssPubkey,
	}
	k.SetPendingNonces(ctx, p)

	return &types.MsgResetChainNoncesResponse{}, nil
}

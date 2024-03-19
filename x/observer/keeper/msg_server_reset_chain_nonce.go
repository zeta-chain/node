package keeper

import (
	"context"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

type MsgResetChainNonces struct {
	Creator string
	ChainID int64
	//ChainNonces   types.ChainNonces
	//PendingNonces types.PendingNonces
}

func (k msgServer) ResetChainNonces(goCtx context.Context, msg *MsgResetChainNonces) (*types.MsgAddBlameVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	tss, found := k.GetTSS(ctx)
	if !found {
		return nil, errors.New("no tss found")
	}

	chain := common.GetChainFromChainID(msg.ChainID)
	if chain == nil {
		return nil, errors.New("chain not found")
	}

	// set chain nonce and pending nonce
	chainNonce := types.ChainNonces{
		Index:   chain.ChainName.String(),
		ChainId: chain.ChainId,
		Nonce:   0,
		// #nosec G701 always positive
		FinalizedHeight: uint64(ctx.BlockHeight()),
	}
	k.SetChainNonces(ctx, chainNonce)
	p := types.PendingNonces{
		NonceLow:  0,
		NonceHigh: 0,
		ChainId:   chain.ChainId,
		Tss:       tss.TssPubkey,
	}
	k.SetPendingNonces(ctx, p)

	return &types.MsgAddBlameVoteResponse{}, nil
}

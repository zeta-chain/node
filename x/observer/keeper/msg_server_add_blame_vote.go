package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
	crosschainTypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func (k msgServer) AddBlameVote(goCtx context.Context, vote *types.MsgAddBlameVote) (*types.MsgAddBlameVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// GetChainFromChainID makes sure we are getting only supported chains , if a chain support has been turned on using gov proposal, this function returns nil
	observationChain := k.GetParams(ctx).GetChainFromChainID(vote.ChainId)
	if observationChain == nil {
		return nil, sdkerrors.Wrap(crosschainTypes.ErrUnsupportedChain, fmt.Sprintf("ChainID %d, Blame vote", vote.ChainId))
	}
	// IsAuthorized does various checks against the list of observer mappers
	ok, err := k.IsAuthorized(ctx, vote.Creator, observationChain)
	if !ok {
		return nil, err
	}
	return &types.MsgAddBlameVoteResponse{}, nil
}

func (k msgServer) IsAuthorized(ctx sdk.Context, address string, chain *common.Chain) (bool, error) {
	observerMapper, found := k.GetObserverMapper(ctx, chain)
	if !found {
		return false, sdkerrors.Wrap(crosschainTypes.ErrNotAuthorized, fmt.Sprintf("observer list not present for chain %s", chain.String()))
	}
	for _, obs := range observerMapper.ObserverList {
		if obs == address {
			return true, nil
		}
	}
	return true, nil
}

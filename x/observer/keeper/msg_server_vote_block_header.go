package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	lightclienttypes "github.com/zeta-chain/node/x/lightclient/types"
	"github.com/zeta-chain/node/x/observer/types"
)

const voteBlockHeaderID = "Vote BlockHeader"

// VoteBlockHeader vote for a new block header to the storers
func (k msgServer) VoteBlockHeader(
	goCtx context.Context,
	msg *types.MsgVoteBlockHeader,
) (*types.MsgVoteBlockHeaderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if the chain is enabled
	chain, found := k.GetSupportedChainFromChainID(ctx, msg.ChainId)
	if !found {
		return nil, sdkerrors.Wrapf(
			types.ErrSupportedChains,
			"%s, ChainID %d", voteBlockHeaderID, msg.ChainId)
	}

	err := k.CheckObserverCanVote(ctx, msg.Creator)
	if err != nil {
		return nil, err
	}

	// check the new block header is valid
	parentHash, err := k.lightclientKeeper.CheckNewBlockHeader(ctx, msg.ChainId, msg.BlockHash, msg.Height, msg.Header)
	if err != nil {
		return nil, sdkerrors.Wrapf(
			lightclienttypes.ErrInvalidBlockHeader,
			"%s, parent hash %s", voteBlockHeaderID, parentHash)
	}

	_, isFinalized, isNew, err := k.VoteOnBallot(
		ctx,
		chain,
		msg.Digest(),
		types.ObservationType_InboundTx,
		msg.Creator,
		types.VoteType_SuccessObservation,
	)
	if err != nil {
		return nil, sdkerrors.Wrap(err, voteBlockHeaderID)
	}

	if !isFinalized {
		return &types.MsgVoteBlockHeaderResponse{
			BallotCreated: isNew,
			VoteFinalized: false,
		}, nil
	}

	// Add the new block header to the store.
	k.lightclientKeeper.AddBlockHeader(ctx, msg.ChainId, msg.Height, msg.BlockHash, msg.Header, parentHash)
	return &types.MsgVoteBlockHeaderResponse{
		BallotCreated: isNew,
		VoteFinalized: true,
	}, nil
}

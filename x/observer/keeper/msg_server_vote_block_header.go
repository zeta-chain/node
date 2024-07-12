package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// VoteBlockHeader vote for a new block header to the storers
func (k msgServer) VoteBlockHeader(
	goCtx context.Context,
	msg *types.MsgVoteBlockHeader,
) (*types.MsgVoteBlockHeaderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if the chain is enabled
	chain, found := k.GetSupportedChainFromChainID(ctx, msg.ChainId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrSupportedChains, "chain id: %d", msg.ChainId)
	}

	// check if observer
	if ok := k.IsNonTombstonedObserver(ctx, msg.Creator); !ok {
		return nil, types.ErrNotObserver
	}

	// check the new block header is valid
	parentHash, err := k.lightclientKeeper.CheckNewBlockHeader(ctx, msg.ChainId, msg.BlockHash, msg.Height, msg.Header)
	if err != nil {
		return nil, sdkerrors.Wrap(lightclienttypes.ErrInvalidBlockHeader, err.Error())
	}

	_, isFinalized, isNew, err := k.VoteOnBallot(ctx, chain, msg.Digest(), types.ObservationType_InboundTx, msg.Creator, types.VoteType_SuccessObservation)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to vote on ballot")
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

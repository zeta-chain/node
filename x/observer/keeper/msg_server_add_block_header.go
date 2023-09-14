package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// Updates permissions. Currently, this is only used to enable/disable the
// inbound transactions.
//
// Only the admin policy account is authorized to broadcast this message.
func (k msgServer) AddBlockHeader(goCtx context.Context, msg *types.MsgAddBlockHeader) (*types.MsgAddBlockHeaderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	observationType := types.ObservationType_InBoundTx
	chain := common.GetChainFromChainID(msg.ChainId)
	ok, err := k.IsAuthorized(ctx, msg.Creator, chain)
	if !ok {
		return nil, err
	}

	ballot, _, err := k.FindBallot(ctx, msg.Digest(), chain, observationType)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to find ballot")
	}
	ballot, err = k.AddVoteToBallot(ctx, ballot, msg.Creator, types.VoteType_SuccessObservation)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to add vote to ballot")
	}
	ballot, isFinalized := k.CheckIfFinalizingVote(ctx, ballot)
	if !isFinalized {
		return &types.MsgAddBlockHeaderResponse{}, nil
	}

	_, found := k.GetBlockHeader(ctx, msg.BlockHash)
	if found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("block header with hash %s already exists", msg.BlockHeader))
	}

	pHash, err := msg.ParentHash() // error is checked in BasicValidation in msg; check again for extra caution
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("failed to get parent hash: %s", err.Error()))
	}

	// TODO: add check for parent block header's existence here

	bh := types.BlockHeader{
		Header:     msg.BlockHeader,
		Height:     msg.Height,
		Hash:       msg.BlockHash,
		ParentHash: pHash,
		ChainId:    msg.ChainId,
	}

	k.SetBlockHeader(ctx, bh)
	return &types.MsgAddBlockHeaderResponse{}, nil
}

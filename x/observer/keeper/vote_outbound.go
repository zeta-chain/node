package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/chains"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// VoteOnOutboundBallot casts a vote on an outbound transaction observed on a connected chain (after
// it has been broadcasted to and finalized on a connected chain). If this is
// the first vote, a new ballot is created. When a threshold of votes is
// reached, the ballot is finalized.
// returns if the vote is finalized, if the ballot is new, the ballot status and the name of the observation chain
func (k Keeper) VoteOnOutboundBallot(
	ctx sdk.Context,
	ballotIndex string,
	outTxChainID int64,
	receiveStatus chains.ReceiveStatus,
	voter string,
) (isFinalized bool, isNew bool, ballot observertypes.Ballot, observationChainName string, err error) {
	// Observer Chain already checked then inbound is created
	/* EDGE CASE : Params updated in during the finalization process
	   i.e Inbound has been finalized but outbound is still pending
	*/
	observationChain, found := k.GetSupportedChainFromChainID(ctx, outTxChainID)
	if !found {
		return false, false, ballot, "", observertypes.ErrSupportedChains
	}
	if observertypes.CheckReceiveStatus(receiveStatus) != nil {
		return false, false, ballot, "", observertypes.ErrInvalidStatus
	}

	// check if voter is authorized
	err = k.CheckObserverCanVote(ctx, voter)
	if err != nil {
		return false, false, ballot, "", err
	}

	ballot, isFinalized, isNew, err = k.VoteOnBallot(
		ctx,
		observationChain,
		ballotIndex,
		observertypes.ObservationType_OutboundTx,
		voter,
		observertypes.ConvertReceiveStatusToVoteType(receiveStatus),
	)
	if err != nil {
		return false, false, ballot, "", sdkerrors.Wrap(err, msgVoteOnBallot)
	}

	return isFinalized, isNew, ballot, observationChain.String(), nil
}

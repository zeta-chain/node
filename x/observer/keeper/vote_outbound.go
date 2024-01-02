package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
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
	receiveStatus common.ReceiveStatus,
	voter string,
) (isFinalized bool, isNew bool, ballotStatus observertypes.Ballot, observationChainName string, err error) {
	// Observer Chain already checked then inbound is created
	/* EDGE CASE : Params updated in during the finalization process
	   i.e Inbound has been finalized but outbound is still pending
	*/
	observationChain := k.GetParams(ctx).GetChainFromChainID(outTxChainID)
	if observationChain == nil {
		return false, false, ballotStatus, "", observertypes.ErrSupportedChains
	}
	if observertypes.CheckReceiveStatus(receiveStatus) != nil {
		return false, false, ballotStatus, "", err
	}

	// check if voter is authorized
	if ok := k.IsAuthorized(ctx, voter, observationChain); !ok {
		return false, false, ballotStatus, "", observertypes.ErrNotAuthorizedPolicy
	}

	// fetch or create ballot
	ballot, isNew, err := k.FindBallot(ctx, ballotIndex, observationChain, observertypes.ObservationType_OutBoundTx)
	if err != nil {
		return false, false, ballotStatus, "", err
	}

	// add vote to ballot
	ballot, err = k.AddVoteToBallot(ctx, ballot, voter, observertypes.ConvertReceiveStatusToVoteType(receiveStatus))
	if err != nil {
		return false, false, ballotStatus, "", err
	}

	ballot, isFinalizedInThisBlock := k.CheckIfFinalizingVote(ctx, ballot)
	return isFinalizedInThisBlock, isNew, ballot, observationChain.String(), nil
}

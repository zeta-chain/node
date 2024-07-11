package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// VoteOnBallot finds a ballot or creates a new one if not found,
// and casts a vote on it. Then proceed to check if the vote has been finalized.
// This function holds generic logic for all types of votes.
func (k Keeper) VoteOnBallot(
	ctx sdk.Context,
	chain chains.Chain,
	ballotIndex string,
	observationType types.ObservationType,
	voter string,
	voteType types.VoteType,
) (
	ballot types.Ballot,
	isFinalized bool,
	isNew bool,
	err error) {
	ballot, isNew, err = k.FindBallot(ctx, ballotIndex, chain, observationType)
	if err != nil {
		return ballot, false, false, err
	}

	ballot, err = k.AddVoteToBallot(ctx, ballot, voter, voteType)
	if err != nil {
		return ballot, false, isNew, err
	}

	ballot, isFinalized = k.CheckIfFinalizingVote(ctx, ballot)

	return ballot, isFinalized, isNew, nil
}

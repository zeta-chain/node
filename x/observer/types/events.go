package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	BallotIdentifier       = "BallotIdentifier"
	CCTXIndex              = "CCTXIndex"
	BallotObservationHash  = "BallotObservationHash"
	BallotObservationChain = "BallotObservationChain"
	BallotType             = "BallotType"
	BallotCreated          = "BallotCreated"
)

func EmitEventBallotCreated(ctx sdk.Context, ballot Ballot, observationHash, obserVationChain string) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(BallotCreated,
			sdk.NewAttribute(BallotIdentifier, ballot.BallotIdentifier),
			sdk.NewAttribute(CCTXIndex, ballot.BallotIdentifier),
			sdk.NewAttribute(BallotObservationHash, observationHash),
			sdk.NewAttribute(BallotObservationChain, obserVationChain),
			sdk.NewAttribute(BallotType, ballot.ObservationType.String()),
		),
	)
}

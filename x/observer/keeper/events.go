package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	types "github.com/zeta-chain/zetacore/x/observer/types"
)

func EmitEventBallotCreated(ctx sdk.Context, ballot types.Ballot, observationHash, obserVationChain string) {
	err := ctx.EventManager().EmitTypedEvent(&types.EventBallotCreated{
		MsgTypeUrl:       "/zetachain.zetacore.observer.internal.BallotCreated",
		BallotIdentifier: ballot.BallotIdentifier,
		BallotType:       ballot.ObservationType.String(),
		ObservationHash:  observationHash,
		ObservationChain: obserVationChain,
	})
	if err != nil {
		ctx.Logger().Error("failed to emit EventBallotCreated : %s", err.Error())
	}
}

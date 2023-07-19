package keeper

import (
	types2 "github.com/coinbase/rosetta-sdk-go/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"

	types "github.com/zeta-chain/zetacore/x/observer/types"
)

func EmitEventBallotCreated(ctx sdk.Context, ballot types.Ballot, observationHash, obserVationChain string) {
	err := ctx.EventManager().EmitTypedEvent(&types.EventBallotCreated{
		BallotIdentifier: ballot.BallotIdentifier,
		BallotType:       ballot.ObservationType.String(),
		ObservationHash:  observationHash,
		ObservationChain: obserVationChain,
	})
	if err != nil {
		ctx.Logger().Error("failed to emit EventBallotCreated : %s", err.Error())
	}
}

func EmitEventKeyGenBlockUpdated(ctx sdk.Context, keygen *types.Keygen) {
	err := ctx.EventManager().EmitTypedEvents(&types.EventKeygenBlockUpdated{
		MsgTypeUrl:    sdk.MsgTypeURL(&types.MsgUpdateKeygen{}),
		KeygenBlock:   strconv.Itoa(int(keygen.BlockNumber)),
		KeygenPubkeys: types2.PrettyPrintStruct(keygen.GranteePubkeys),
	})
	if err != nil {
		ctx.Logger().Error("Error emitting EventKeygenBlockUpdated :", err)
	}
}

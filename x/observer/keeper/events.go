package keeper

import (
	"strconv"

	types2 "github.com/coinbase/rosetta-sdk-go/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	types "github.com/zeta-chain/zetacore/x/observer/types"
)

func EmitEventBallotCreated(ctx sdk.Context, ballot types.Ballot, observationHash, observationChain string) {
	err := ctx.EventManager().EmitTypedEvent(&types.EventBallotCreated{
		BallotIdentifier: ballot.BallotIdentifier,
		BallotType:       ballot.ObservationType.String(),
		ObservationHash:  observationHash,
		ObservationChain: observationChain,
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

func EmitEventAddObserver(ctx sdk.Context, observerCount uint64, operatorAddress, zetaclientGranteeAddress, zetaclientGranteePubkey string) {
	err := ctx.EventManager().EmitTypedEvents(&types.EventNewObserverAdded{
		MsgTypeUrl:               sdk.MsgTypeURL(&types.MsgAddObserver{}),
		ObserverAddress:          operatorAddress,
		ZetaclientGranteeAddress: zetaclientGranteeAddress,
		ZetaclientGranteePubkey:  zetaclientGranteePubkey,
		ObserverLastBlockCount:   observerCount,
	})
	if err != nil {
		ctx.Logger().Error("Error emitting EmitEventAddObserver :", err)
	}
}

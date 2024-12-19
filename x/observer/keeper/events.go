package keeper

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	types "github.com/zeta-chain/node/x/observer/types"
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

// vendor this code from github.com/coinbase/rosetta-sdk-go/types
func prettyPrintStruct(val interface{}) string {
	prettyStruct, err := json.MarshalIndent(
		val,
		"",
		" ",
	)
	if err != nil {
		log.Fatal(err)
	}

	return string(prettyStruct)
}

func EmitEventBallotDeleted(ctx sdk.Context, ballot types.Ballot) {
	var voterList []types.VoterList
	voterList, err := ballot.GenerateVoterList()
	if err != nil {
		ctx.Logger().
			Error(fmt.Sprintf("failed to generate voter list for ballot %s", ballot.BallotIdentifier), err.Error())
	}
	err = ctx.EventManager().EmitTypedEvent(&types.EventBallotDeleted{
		MsgTypeUrl:       "zetachain.zetacore.observer.internal.BallotDeleted",
		BallotIdentifier: ballot.BallotIdentifier,
		BallotType:       ballot.ObservationType.String(),
		Voters:           voterList,
	})
	if err != nil {
		ctx.Logger().Error("failed to emit EventBallotDeleted : %s", err.Error())
	}
}

func EmitEventKeyGenBlockUpdated(ctx sdk.Context, keygen *types.Keygen) {
	err := ctx.EventManager().EmitTypedEvents(&types.EventKeygenBlockUpdated{
		MsgTypeUrl:    sdk.MsgTypeURL(&types.MsgUpdateKeygen{}),
		KeygenBlock:   strconv.Itoa(int(keygen.BlockNumber)),
		KeygenPubkeys: prettyPrintStruct(keygen.GranteePubkeys),
	})
	if err != nil {
		ctx.Logger().Error("Error emitting EventKeygenBlockUpdated :", err)
	}
}

func EmitEventAddObserver(
	ctx sdk.Context,
	observerCount uint64,
	operatorAddress, zetaclientGranteeAddress, zetaclientGranteePubkey string,
) {
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

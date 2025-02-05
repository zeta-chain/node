package simulation

import (
	"fmt"
	"strings"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/cosmos/gogoproto/proto"
)

// GenAndDeliverTxWithRandFees generates a transaction with a random fee and delivers it.
func GenAndDeliverTxWithRandFees(
	txCtx simulation.OperationInput, errorOnDeliver bool,
) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
	account := txCtx.AccountKeeper.GetAccount(txCtx.Context, txCtx.SimAccount.Address)
	spendable := txCtx.Bankkeeper.SpendableCoins(txCtx.Context, account.GetAddress())

	var fees sdk.Coins
	var err error

	coins, hasNeg := spendable.SafeSub(txCtx.CoinsSpentInMsg...)
	if hasNeg {
		return simtypes.NoOpMsg(
			txCtx.ModuleName,
			sdk.MsgTypeURL(txCtx.Msg),
			"message doesn't leave room for fees",
		), nil, err
	}

	fees, err = simtypes.RandomFees(txCtx.R, txCtx.Context, coins)
	if err != nil {
		return simtypes.NoOpMsg(txCtx.ModuleName, sdk.MsgTypeURL(txCtx.Msg), "unable to generate fees"), nil, err
	}
	if errorOnDeliver {
		return GenAndDeliverTxWithError(txCtx, fees)
	}
	return GenAndDeliverTxNoError(txCtx, fees)
}

// GenAndDeliverTxNoError generates a transactions and delivers it with the provided fees.
// This function does not return an error if the transaction fails to deliver.
func GenAndDeliverTxNoError(
	txCtx simulation.OperationInput,
	fees sdk.Coins,
) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
	account := txCtx.AccountKeeper.GetAccount(txCtx.Context, txCtx.SimAccount.Address)
	tx, err := simtestutil.GenSignedMockTx(
		txCtx.R,
		txCtx.TxGen,
		[]sdk.Msg{txCtx.Msg},
		fees,
		simtestutil.DefaultGenTxGas,
		txCtx.Context.ChainID(),
		[]uint64{account.GetAccountNumber()},
		[]uint64{account.GetSequence()},
		txCtx.SimAccount.PrivKey,
	)
	if err != nil {
		return simtypes.NoOpMsg(txCtx.ModuleName, sdk.MsgTypeURL(txCtx.Msg), "unable to generate mock tx"), nil, err
	}

	// Do not return error here, this method of delivery is used when we are not certain if the tx can fail or succeed.
	// Example:
	//When adding future operations, for votes,
	//we don't know if the vote will pass or fail.This is
	//because the observer set might change in the next block and cause the vote to fail.
	_, _, err = txCtx.App.SimDeliver(txCtx.TxGen.TxEncoder(), tx)
	if err != nil {
		return simtypes.NoOpMsg(txCtx.ModuleName, sdk.MsgTypeURL(txCtx.Msg), "unable to deliver tx"), nil, nil
	}

	return OperationMessage(txCtx.Msg), nil, nil
}

func GenAndDeliverTxWithError(
	txCtx simulation.OperationInput,
	fees sdk.Coins,
) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
	account := txCtx.AccountKeeper.GetAccount(txCtx.Context, txCtx.SimAccount.Address)
	tx, err := simtestutil.GenSignedMockTx(
		txCtx.R,
		txCtx.TxGen,
		[]sdk.Msg{txCtx.Msg},
		fees,
		simtestutil.DefaultGenTxGas,
		txCtx.Context.ChainID(),
		[]uint64{account.GetAccountNumber()},
		[]uint64{account.GetSequence()},
		txCtx.SimAccount.PrivKey,
	)
	if err != nil {
		return simtypes.NoOpMsg(txCtx.ModuleName, sdk.MsgTypeURL(txCtx.Msg), "unable to generate mock tx"), nil, err
	}

	// Return error here, we don't expect this tx to fail.
	_, _, err = txCtx.App.SimDeliver(txCtx.TxGen.TxEncoder(), tx)
	if err != nil {
		return simtypes.NoOpMsg(txCtx.ModuleName, sdk.MsgTypeURL(txCtx.Msg), "unable to deliver tx"), nil, err
	}

	return OperationMessage(txCtx.Msg), nil, nil
}

// OperationMessage parses the msg type to extract the module name and returns a new OperationMsg
// We use this instead of the default provided by the sdk because , the default implementation uses the second element of the split
// default implementation: "cosmos.bank.v1beta1.MsgSend" => "bank" [Second element,this position is hardcoded]
// This is not accurate for us , since our urls are of the format "zetachain.zetacore.crosschain.MsgVoteInbound"
// and we need to fetch the third element
func OperationMessage(msg proto.Message) simtypes.OperationMsg {
	msgType := sdk.MsgTypeURL(msg)
	parts := strings.Split(msgType, ".")
	moduleName := "unknown"
	if len(parts) >= 3 {
		moduleName = parts[2]
	}

	// Its better to panic here than return an error. This code is only used by simulation tests
	// Returning and error would need additional error handling in the tests itself , without providing any benefits
	protoBz, err := proto.Marshal(msg)
	if err != nil {
		panic(fmt.Sprintf("unable to marshal proto message %v", err))
	}

	return simtypes.NewOperationMsgBasic(moduleName, msgType, "", true, protoBz)
}

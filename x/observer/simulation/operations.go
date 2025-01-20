package simulation

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/codec"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/observer/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

// Simulation operation weights constants
// Operation weights are used by the simulation program to simulate the weight of different operations.
// This decides what percentage of a certain type of operation is part of a block.
// Based on the weights assigned in the cosmos sdk modules , 100 seems to the max weight used , and therefore guarantees that at least one operation of that type is present in a block.
// Operation weights are used by the `SimulateFromSeed`
// function to pick a random operation based on the weights.The functions with higher weights are more likely to be picked.

// Therefore, this decides the percentage of a certain operation that is part of a block.

// Based on the weights assigned in the cosmos sdk modules,
// 100 seems to the max weight used,and we should use relative weights
// to signify the number of each operation in a block.
const (
	OpWeightMsgTypeMsgEnableCCTX                  = "op_weight_msg_enable_crosschain_flags"         // #nosec G101 not a hardcoded credential
	OpWeightMsgTypeMsgDisableCCTX                 = "op_weight_msg_disable_crosschain_flags"        // #nosec G101 not a hardcoded credential
	OpWeightMsgTypeMsgVoteTSS                     = "op_weight_msg_vote_tss"                        // #nosec G101 not a hardcoded credential
	OpWeightMsgTypeMsgUpdateKeygen                = "op_weight_msg_update_keygen"                   // #nosec G101 not a hardcoded credential
	OpWeightMsgTypeMsgUpdateObserver              = "op_weight_msg_update_observer"                 // #nosec G101 not a hardcoded credential
	OpWeightMsgTypeMsgUpdateChainParams           = "op_weight_msg_update_chain_params"             // #nosec G101 not a hardcoded credential
	OpWeightMsgTypeMsgRemoveChainParams           = "op_weight_msg_remove_chain_params"             // #nosec G101 not a hardcoded credential
	OpWeightMsgTypeMsgResetChainNonces            = "op_weight_msg_reset_chain_nonces"              // #nosec G101 not a hardcoded credential
	OpWeightMsgTypeMsgUpdateGasPriceIncreaseFlags = "op_weight_msg_update_gas_price_increase_flags" // #nosec G101 not a hardcoded credential
	OpWeightMsgTypeMsgAddObserver                 = "op_weight_msg_add_observer"                    // #nosec G101 not a hardcoded credential

	// DefaultWeightMsgTypeMsgEnableCCTX We use a high weight for this operation
	// to ensure that it is present in the block more number of times than any operation that changes the validator set
	// Arrived at this number based on the weights used in the cosmos sdk staking module and through some trial and error
	DefaultWeightMsgTypeMsgEnableCCTX                  = 100
	DefaultWeightMsgTypeMsgDisableCCTX                 = 10
	DefaultWeightMsgTypeMsgVoteTSS                     = 10
	DefaultWeightMsgTypeMsgUpdateKeygen                = 10
	DefaultWeightMsgTypeMsgUpdateObserver              = 10
	DefaultWeightMsgTypeMsgUpdateChainParams           = 10
	DefaultWeightMsgTypeMsgRemoveChainParams           = 10
	DefaultWeightMsgTypeMsgResetChainNonces            = 5
	DefaultWeightMsgTypeMsgUpdateGasPriceIncreaseFlags = 10
	DefaultWeightMsgTypeMsgAddObserver                 = 5

	DefaultRetryCount = 10
)

// WeightedOperations for observer module
func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONCodec, k keeper.Keeper,
) simulation.WeightedOperations {
	var (
		weightMsgTypeMsgEnableCCTX                  int
		weightMsgTypeMsgDisableCCTX                 int
		weightMsgTypeMsgVoteTSS                     int
		weightMsgTypeMsgUpdateKeygen                int
		weightMsgTypeMsgUpdateObserver              int
		weightMsgTypeMsgUpdateChainParams           int
		weightMsgTypeMsgRemoveChainParams           int
		weightMsgTypeMsgResetChainNonces            int
		weightMsgTypeMsgUpdateGasPriceIncreaseFlags int
		weightMsgTypeMsgAddObserver                 int
	)

	appParams.GetOrGenerate(OpWeightMsgTypeMsgEnableCCTX, &weightMsgTypeMsgEnableCCTX, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgEnableCCTX = DefaultWeightMsgTypeMsgEnableCCTX
		})

	appParams.GetOrGenerate(OpWeightMsgTypeMsgDisableCCTX, &weightMsgTypeMsgDisableCCTX, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgDisableCCTX = DefaultWeightMsgTypeMsgDisableCCTX
		})

	appParams.GetOrGenerate(OpWeightMsgTypeMsgVoteTSS, &weightMsgTypeMsgVoteTSS, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgVoteTSS = DefaultWeightMsgTypeMsgVoteTSS
		})

	appParams.GetOrGenerate(OpWeightMsgTypeMsgUpdateKeygen, &weightMsgTypeMsgUpdateKeygen, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgUpdateKeygen = DefaultWeightMsgTypeMsgUpdateKeygen
		})

	appParams.GetOrGenerate(OpWeightMsgTypeMsgUpdateObserver, &weightMsgTypeMsgUpdateObserver, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgUpdateObserver = DefaultWeightMsgTypeMsgUpdateObserver
		})

	appParams.GetOrGenerate(OpWeightMsgTypeMsgUpdateChainParams, &weightMsgTypeMsgUpdateChainParams, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgUpdateChainParams = DefaultWeightMsgTypeMsgUpdateChainParams
		})

	appParams.GetOrGenerate(OpWeightMsgTypeMsgRemoveChainParams, &weightMsgTypeMsgRemoveChainParams, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgRemoveChainParams = DefaultWeightMsgTypeMsgRemoveChainParams
		})

	appParams.GetOrGenerate(OpWeightMsgTypeMsgResetChainNonces, &weightMsgTypeMsgResetChainNonces, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgResetChainNonces = DefaultWeightMsgTypeMsgResetChainNonces
		})

	appParams.GetOrGenerate(
		OpWeightMsgTypeMsgUpdateGasPriceIncreaseFlags,
		&weightMsgTypeMsgUpdateGasPriceIncreaseFlags,
		nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgUpdateGasPriceIncreaseFlags = DefaultWeightMsgTypeMsgUpdateGasPriceIncreaseFlags
		},
	)

	appParams.GetOrGenerate(OpWeightMsgTypeMsgAddObserver, &weightMsgTypeMsgAddObserver, nil,
		func(_ *rand.Rand) {
			weightMsgTypeMsgAddObserver = DefaultWeightMsgTypeMsgAddObserver
		})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgTypeMsgEnableCCTX,
			SimulateEnableCCTX(k),
		),

		simulation.NewWeightedOperation(
			weightMsgTypeMsgDisableCCTX,
			SimulateDisableCCTX(k),
		),

		simulation.NewWeightedOperation(
			weightMsgTypeMsgUpdateKeygen,
			SimulateUpdateKeygen(k),
		),
		//
		simulation.NewWeightedOperation(
			weightMsgTypeMsgUpdateChainParams,
			SimulateUpdateChainParams(k),
		),
		//
		//simulation.NewWeightedOperation(
		//	weightMsgTypeMsgRemoveChainParams,
		//	SimulateMsgRemoveChainParams(k),
		//),

		simulation.NewWeightedOperation(
			weightMsgTypeMsgResetChainNonces,
			SimulateResetChainNonces(k),
		),

		simulation.NewWeightedOperation(
			weightMsgTypeMsgUpdateGasPriceIncreaseFlags,
			SimulateUpdateGasPriceIncreaseFlags(k),
		),

		simulation.NewWeightedOperation(
			weightMsgTypeMsgAddObserver,
			SimulateUpdateObserver(k),
		),

		simulation.NewWeightedOperation(
			weightMsgTypeMsgAddObserver,
			SimulateAddObserverNodeAccount(k),
		),

		simulation.NewWeightedOperation(
			weightMsgTypeMsgAddObserver,
			SimulateAddObserver(k),
		),

		simulation.NewWeightedOperation(
			weightMsgTypeMsgVoteTSS,
			SimulateMsgVoteTSS(k),
		),
	}
}

func GetPolicyAccount(ctx sdk.Context, k types.AuthorityKeeper, accounts []simtypes.Account) (simtypes.Account, error) {
	policies, found := k.GetPolicies(ctx)
	if !found {
		return simtypes.Account{}, fmt.Errorf("policies object not found")
	}
	if len(policies.Items) == 0 {
		return simtypes.Account{}, fmt.Errorf("no policies found")
	}

	admin := policies.Items[0].Address
	address, err := types.GetOperatorAddressFromAccAddress(admin)
	if err != nil {
		return simtypes.Account{}, err
	}
	simAccount, found := simtypes.FindAccount(accounts, address)
	if !found {
		return simtypes.Account{}, fmt.Errorf("admin account not found in list of simulation accounts")
	}
	return simAccount, nil
}

func GetExternalChain(ctx sdk.Context, k keeper.Keeper, r *rand.Rand) (chains.Chain, error) {
	supportedChains := k.GetSupportedChains(ctx)
	if len(supportedChains) == 0 {
		return chains.Chain{}, fmt.Errorf("no supported chains found")
	}
	externalChain := chains.Chain{}
	foundExternalChain := RepeatCheck(func() bool {
		c := supportedChains[r.Intn(len(supportedChains))]
		if !c.IsZetaChain() {
			externalChain = c
			return true
		}
		return false
	})

	if !foundExternalChain {
		return chains.Chain{}, fmt.Errorf("no external chain found")
	}
	return externalChain, nil
}

// GetRandomAccountAndObserver returns a random account and the associated observer address
func GetRandomAccountAndObserver(
	r *rand.Rand,
	ctx sdk.Context,
	k keeper.Keeper,
	accounts []simtypes.Account,
) (simtypes.Account, string, []string, error) {
	observerList := []string{}
	observers, found := k.GetObserverSet(ctx)
	if !found {
		return simtypes.Account{}, "", observerList, fmt.Errorf("observer set not found")
	}

	observerList = observers.ObserverList

	if len(observers.ObserverList) == 0 {
		return simtypes.Account{}, "", observerList, fmt.Errorf("no observers present in observer set found")
	}

	randomObserver := ""
	foundObserver := RepeatCheck(func() bool {
		randomObserver = GetRandomObserver(r, observerList)
		_, foundNodeAccount := k.GetNodeAccount(ctx, randomObserver)
		if !foundNodeAccount {
			return false
		}
		ok := k.IsNonTombstonedObserver(ctx, randomObserver)
		if ok {
			return true
		}
		return false
	})

	if !foundObserver {
		return simtypes.Account{}, "", nil, fmt.Errorf("no observer found")
	}

	simAccount, err := GetSimAccount(randomObserver, accounts)
	if err != nil {
		return simtypes.Account{}, "", observerList, err
	}
	return simAccount, randomObserver, observerList, nil
}

func GetRandomNodeAccount(
	r *rand.Rand,
	ctx sdk.Context,
	k keeper.Keeper,
	accounts []simtypes.Account,
) (simtypes.Account, string, error) {
	nodeAccounts := k.GetAllNodeAccount(ctx)

	if len(nodeAccounts) == 0 {
		return simtypes.Account{}, "", fmt.Errorf("no node accounts present")
	}

	randomNodeAccount := nodeAccounts[r.Intn(len(nodeAccounts))].Operator

	simAccount, err := GetSimAccount(randomNodeAccount, accounts)
	if err != nil {
		return simtypes.Account{}, "", err
	}
	return simAccount, randomNodeAccount, nil
}

func GetRandomObserver(r *rand.Rand, observerList []string) string {
	idx := r.Intn(len(observerList))
	return observerList[idx]
}

// GetSimAccount returns the account associated with the observer address from the list of accounts provided
// GetSimAccount can fail if all the observers are removed from the observer set ,this can happen
//if the other modules create transactions which affect the validator
//and triggers any of the staking hooks defined in the observer modules

func GetSimAccount(observerAddress string, accounts []simtypes.Account) (simtypes.Account, error) {
	operatorAddress, err := types.GetOperatorAddressFromAccAddress(observerAddress)
	if err != nil {
		return simtypes.Account{}, fmt.Errorf("validator not found for observer ")
	}

	simAccount, found := simtypes.FindAccount(accounts, operatorAddress)
	if !found {
		return simtypes.Account{}, fmt.Errorf("operator account not found")
	}
	return simAccount, nil
}

func RepeatCheck(fn func() bool) bool {
	for i := 0; i < DefaultRetryCount; i++ {
		if fn() {
			return true
		}
	}
	return false
}

func ObserverVotesSimulationMatrix() (simtypes.TransitionMatrix, []float64, int) {
	observerVotesTransitionMatrix, _ := simulation.CreateTransitionMatrix([][]int{
		{20, 10, 0, 0, 0, 0},
		{55, 50, 20, 10, 0, 0},
		{25, 25, 30, 25, 30, 15},
		{0, 15, 30, 25, 30, 30},
		{0, 0, 20, 30, 30, 30},
		{0, 0, 0, 10, 10, 25},
	})
	// The states are:
	// column 1: All observers vote
	// column 2: 90% vote
	// column 3: 75% vote
	// column 4: 40% vote
	// column 5: 15% vote
	// column 6: noone votes
	// All columns sum to 100 for simplicity, but this is arbitrary and can be changed
	statePercentageArray := []float64{1, .9, .75, .4, .15, 0}
	curNumVotesState := 1
	return observerVotesTransitionMatrix, statePercentageArray, curNumVotesState
}

func BallotVoteSimulationMatrix() (simtypes.TransitionMatrix, []float64, int) {
	ballotTransitionMatrix, _ := simulation.CreateTransitionMatrix([][]int{
		{70, 10},
		{30, 10},
	})
	// The states are:
	// column 1: 100% vote yes
	// column 2: 0% vote yes
	// For all conditions we assume if the vote is not a yes
	// then it is a no .
	yesVoteArray := []float64{1, 0}
	ballotVotesState := 1
	return ballotTransitionMatrix, yesVoteArray, ballotVotesState
}

// GenAndDeliverTxWithRandFees generates a transaction with a random fee and delivers it.
func GenAndDeliverTxWithRandFees(
	txCtx simulation.OperationInput,
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
	return GenAndDeliverTx(txCtx, fees)
}

// GenAndDeliverTx generates a transactions and delivers it with the provided fees.
// This function does not return an error if the transaction fails to deliver.
func GenAndDeliverTx(
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

	_, _, err = txCtx.App.SimDeliver(txCtx.TxGen.TxEncoder(), tx)
	if err != nil {
		return simtypes.NoOpMsg(txCtx.ModuleName, sdk.MsgTypeURL(txCtx.Msg), "unable to deliver tx"), nil, nil
	}

	return simtypes.NewOperationMsg(txCtx.Msg, true, ""), nil, nil
}

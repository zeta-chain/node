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
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
	observerTypes "github.com/zeta-chain/node/x/observer/types"
)

// Simulation operation weights constants
// Operation weights are used by the simulation program to simulate the weight of different operations.
// This decides what percentage of a certain type of operation is part of a block.
// Based on the weights assigned in the cosmos sdk modules , 100 seems to the max weight used , and therefore guarantees that at least one operation of that type is present in a block.
// TODO Add more details to comment based on what the number represents in terms of percentage of operations in a block
// https://github.com/zeta-chain/node/issues/3100
const (
	DefaultWeightAddOutboundTracker            = 100
	DefaultWeightAddInboundTracker             = 20
	DefaultWeightRemoveOutboundTracker         = 10
	DefaultWeightVoteGasPrice                  = 100
	DefaultWeightVoteOutbound                  = 100
	DefaultWeightVoteInbound                   = 100
	DefaultWeightWhitelistERC20                = 10
	DefaultWeightMigrateTssFunds               = 1
	DefaultWeightUpdateTssAddress              = 1
	DefaultWeightAbortStuckCCTX                = 10
	DefaultWeightUpdateRateLimiterFlags        = 10
	DefaultWeightRefundAbortedCCTX             = 10
	DefaultWeightUpdateERC20CustodyPauseStatus = 10

	OpWeightMsgAddOutboundTracker          = "op_weight_msg_add_outbound_tracker"              // #nosec G101 not a hardcoded credential
	OpWeightAddInboundTracker              = "op_weight_msg_add_inbound_tracker"               // #nosec G101 not a hardcoded credential
	OpWeightRemoveOutboundTracker          = "op_weight_msg_remove_outbound_tracker"           // #nosec G101 not a hardcoded credential
	OpWeightVoteGasPrice                   = "op_weight_msg_vote_gas_price"                    // #nosec G101 not a hardcoded credential
	OpWeightVoteOutbound                   = "op_weight_msg_vote_outbound"                     // #nosec G101 not a hardcoded credential
	OpWeightVoteInbound                    = "op_weight_msg_vote_inbound"                      // #nosec G101 not a hardcoded credential
	OpWeightWhitelistERC20                 = "op_weight_msg_whitelist_erc20"                   // #nosec G101 not a hardcoded credential
	OpWeightMigrateTssFunds                = "op_weight_msg_migrate_tss_funds"                 // #nosec G101 not a hardcoded credential
	OpWeightUpdateTssAddress               = "op_weight_msg_update_tss_address"                // #nosec G101 not a hardcoded credential
	OpWeightAbortStuckCCTX                 = "op_weight_msg_abort_stuck_cctx"                  // #nosec G101 not a hardcoded credential
	OpWeightUpdateRateLimiterFlags         = "op_weight_msg_update_rate_limiter_flags"         // #nosec G101 not a hardcoded credential
	OpWeightRefundAbortedCCTX              = "op_weight_msg_refund_aborted_cctx"               // #nosec G101 not a hardcoded credential
	OppWeightUpdateERC20CustodyPauseStatus = "op_weight_msg_update_erc20_custody_pause_status" // #nosec G101 not a hardcoded credential

)

func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONCodec, k keeper.Keeper) simulation.WeightedOperations {
	var (
		weightAddOutboundTracker            int
		weightAddInboundTracker             int
		weightRemoveOutboundTracker         int
		weightVoteGasPrice                  int
		weightVoteOutbound                  int
		weightVoteInbound                   int
		weightWhitelistERC20                int
		weightMigrateTssFunds               int
		weightUpdateTssAddress              int
		weightAbortStuckCCTX                int
		weightUpdateRateLimiterFlags        int
		weightRefundAbortedCCTX             int
		weightUpdateERC20CustodyPauseStatus int
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgAddOutboundTracker, &weightAddOutboundTracker, nil,
		func(_ *rand.Rand) {
			weightAddOutboundTracker = DefaultWeightAddOutboundTracker
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightAddInboundTracker, &weightAddInboundTracker, nil,
		func(_ *rand.Rand) {
			weightAddInboundTracker = DefaultWeightAddInboundTracker
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightRemoveOutboundTracker, &weightRemoveOutboundTracker, nil,
		func(_ *rand.Rand) {
			weightRemoveOutboundTracker = DefaultWeightRemoveOutboundTracker
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightVoteGasPrice, &weightVoteGasPrice, nil,
		func(_ *rand.Rand) {
			weightVoteGasPrice = DefaultWeightVoteGasPrice
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightVoteOutbound, &weightVoteOutbound, nil,
		func(_ *rand.Rand) {
			weightVoteOutbound = DefaultWeightVoteOutbound
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightVoteInbound, &weightVoteInbound, nil,
		func(_ *rand.Rand) {
			weightVoteInbound = DefaultWeightVoteInbound
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightWhitelistERC20, &weightWhitelistERC20, nil,
		func(_ *rand.Rand) {
			weightWhitelistERC20 = DefaultWeightWhitelistERC20
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMigrateTssFunds, &weightMigrateTssFunds, nil,
		func(_ *rand.Rand) {
			weightMigrateTssFunds = DefaultWeightMigrateTssFunds
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightUpdateTssAddress, &weightUpdateTssAddress, nil,
		func(_ *rand.Rand) {
			weightUpdateTssAddress = DefaultWeightUpdateTssAddress
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightAbortStuckCCTX, &weightAbortStuckCCTX, nil,
		func(_ *rand.Rand) {
			weightAbortStuckCCTX = DefaultWeightAbortStuckCCTX
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightUpdateRateLimiterFlags, &weightUpdateRateLimiterFlags, nil,
		func(_ *rand.Rand) {
			weightUpdateRateLimiterFlags = DefaultWeightUpdateRateLimiterFlags
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightRefundAbortedCCTX, &weightRefundAbortedCCTX, nil,
		func(_ *rand.Rand) {
			weightRefundAbortedCCTX = DefaultWeightRefundAbortedCCTX
		},
	)

	appParams.GetOrGenerate(cdc, OppWeightUpdateERC20CustodyPauseStatus, &weightUpdateERC20CustodyPauseStatus, nil,
		func(_ *rand.Rand) {
			weightUpdateERC20CustodyPauseStatus = DefaultWeightUpdateERC20CustodyPauseStatus
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightVoteGasPrice,
			SimulateMsgVoteGasPrice(k),
		),
		simulation.NewWeightedOperation(
			weightVoteInbound,
			SimulateVoteInbound(k),
		),
		simulation.NewWeightedOperation(
			weightVoteOutbound,
			SimulateVoteOutbound(k),
		),
		simulation.NewWeightedOperation(
			weightAddInboundTracker,
			SimulateMsgAddInboundTracker(k),
		),
		simulation.NewWeightedOperation(
			weightAddOutboundTracker,
			SimulateMsgAddOutboundTracker(k),
		),
		simulation.NewWeightedOperation(
			weightRemoveOutboundTracker,
			SimulateMsgRemoveOutboundTracker(k),
		),
		simulation.NewWeightedOperation(
			weightWhitelistERC20,
			SimulateMsgWhitelistERC20(k),
		),
		simulation.NewWeightedOperation(
			weightAbortStuckCCTX,
			SimulateMsgAbortStuckCCTX(k),
		),
		simulation.NewWeightedOperation(
			weightRefundAbortedCCTX,
			SimulateMsgRefundAbortedCCTX(k),
		),
		simulation.NewWeightedOperation(
			weightUpdateRateLimiterFlags,
			SimulateMsgUpdateRateLimiterFlags(k),
		),
		simulation.NewWeightedOperation(
			weightUpdateERC20CustodyPauseStatus,
			SimulateUpdateERC20CustodyPauseStatus(k),
		),
	}
}

func GetRandomObserver(r *rand.Rand, observerList []string) string {
	idx := r.Intn(len(observerList))
	return observerList[idx]
}

func GetRandomChainID(r *rand.Rand, chains []chains.Chain) int64 {
	idx := r.Intn(len(chains))
	return chains[idx].ChainId
}

// GetRandomAccountAndObserver returns a random account and the associated observer address
func GetRandomAccountAndObserver(
	r *rand.Rand,
	ctx sdk.Context,
	k keeper.Keeper,
	accounts []simtypes.Account,
) (simtypes.Account, string, error) {
	observers, found := k.GetObserverKeeper().GetObserverSet(ctx)
	if !found {
		return simtypes.Account{}, "", fmt.Errorf("observer set not found")
	}

	if len(observers.ObserverList) == 0 {
		return simtypes.Account{}, "", fmt.Errorf("no observers present in observer set found")
	}

	randomObserver := GetRandomObserver(r, observers.ObserverList)
	simAccount, err := GetObserverAccount(randomObserver, accounts)
	if err != nil {
		return simtypes.Account{}, "", err
	}
	return simAccount, randomObserver, nil
}

// GetObserverAccount returns the account associated with the observer address from the list of accounts provided
// GetObserverAccount can fail if all the observers are removed from the observer set ,this can happen
//if the other modules create transactions which affect the validator
//and triggers any of the staking hooks defined in the observer modules

func GetObserverAccount(observerAddress string, accounts []simtypes.Account) (simtypes.Account, error) {
	operatorAddress, err := observerTypes.GetOperatorAddressFromAccAddress(observerAddress)
	if err != nil {
		return simtypes.Account{}, fmt.Errorf("validator not found for observer ")
	}

	simAccount, found := simtypes.FindAccount(accounts, operatorAddress)
	if !found {
		return simtypes.Account{}, fmt.Errorf("operator account not found")
	}
	return simAccount, nil
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
		return simtypes.NoOpMsg(txCtx.ModuleName, txCtx.MsgType, "message doesn't leave room for fees"), nil, err
	}

	fees, err = simtypes.RandomFees(txCtx.R, txCtx.Context, coins)
	if err != nil {
		return simtypes.NoOpMsg(txCtx.ModuleName, txCtx.MsgType, "unable to generate fees"), nil, err
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
		return simtypes.NoOpMsg(txCtx.ModuleName, txCtx.MsgType, "unable to generate mock tx"), nil, err
	}

	_, _, err = txCtx.App.SimDeliver(txCtx.TxGen.TxEncoder(), tx)
	if err != nil {
		return simtypes.NoOpMsg(txCtx.ModuleName, txCtx.MsgType, "unable to deliver tx"), nil, nil
	}

	return simtypes.NewOperationMsg(txCtx.Msg, true, "", txCtx.Cdc), nil, nil
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
		{70, 10, 20},
		{20, 30, 30},
		{10, 60, 50},
	})
	// The states are:
	// column 1: 100% vote yes
	// column 2: 50% vote yes
	// column 3: 0% vote yes
	// For all conditions we assume if the the vote is not a yes.
	// then it is a no .Not voting condtion is handled by the ObserverVotesSimulationMatrix matrix
	yesVoteArray := []float64{1, .5, 0}
	ballotVotesState := 1
	return ballotTransitionMatrix, yesVoteArray, ballotVotesState
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
	address, err := observerTypes.GetOperatorAddressFromAccAddress(admin)
	if err != nil {
		return simtypes.Account{}, err
	}
	simAccount, found := simtypes.FindAccount(accounts, address)
	if !found {
		return simtypes.Account{}, fmt.Errorf("admin account not found in list of simulation accounts")
	}
	return simAccount, nil
}

func GetAsset(ctx sdk.Context, k types.FungibleKeeper, chainID int64) (string, error) {
	foreignCoins := k.GetAllForeignCoins(ctx)
	asset := ""

	for _, coin := range foreignCoins {
		if coin.ForeignChainId == chainID {
			return coin.Asset, nil
		}
	}

	return asset, fmt.Errorf("asset not found for chain %d", chainID)
}

package simulation

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/zeta-chain/node/pkg/authz"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
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
	DefaultWeightMsgAddOutboundTracker  = 50
	DefaultWeightAddInboundTracker      = 50
	DefaultWeightRemoveOutboundTracker  = 5
	DefaultWeightVoteGasPrice           = 100
	DefaultWeightVoteOutbound           = 50
	DefaultWeightVoteInbound            = 100
	DefaultWeightWhitelistERC20         = 1
	DefaultWeightMigrateTssFunds        = 1
	DefaultWeightUpdateTssAddress       = 1
	DefaultWeightAbortStuckCCTX         = 10
	DefaultWeightUpdateRateLimiterFlags = 1

	OpWeightMsgAddOutboundTracker  = "op_weight_msg_add_outbound_tracker"      // #nosec G101 not a hardcoded credential
	OpWeightAddInboundTracker      = "op_weight_msg_add_inbound_tracker"       // #nosec G101 not a hardcoded credential
	OpWeightRemoveOutboundTracker  = "op_weight_msg_remove_outbound_tracker"   // #nosec G101 not a hardcoded credential
	OpWeightVoteGasPrice           = "op_weight_msg_vote_gas_price"            // #nosec G101 not a hardcoded credential
	OpWeightVoteOutbound           = "op_weight_msg_vote_outbound"             // #nosec G101 not a hardcoded credential
	OpWeightVoteInbound            = "op_weight_msg_vote_inbound"              // #nosec G101 not a hardcoded credential
	OpWeightWhitelistERC20         = "op_weight_msg_whitelist_erc20"           // #nosec G101 not a hardcoded credential
	OpWeightMigrateTssFunds        = "op_weight_msg_migrate_tss_funds"         // #nosec G101 not a hardcoded credential
	OpWeightUpdateTssAddress       = "op_weight_msg_update_tss_address"        // #nosec G101 not a hardcoded credential
	OpWeightAbortStuckCCTX         = "op_weight_msg_abort_stuck_cctx"          // #nosec G101 not a hardcoded credential
	OpWeightUpdateRateLimiterFlags = "op_weight_msg_update_rate_limiter_flags" // #nosec G101 not a hardcoded credential

)

func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONCodec, k keeper.Keeper) simulation.WeightedOperations {
	var (
		weightMsgAddOutboundTracker  int
		weightAddInboundTracker      int
		weightRemoveOutboundTracker  int
		weightVoteGasPrice           int
		weightVoteOutbound           int
		weightVoteInbound            int
		weightWhitelistERC20         int
		weightMigrateTssFunds        int
		weightUpdateTssAddress       int
		weightAbortStuckCCTX         int
		weightUpdateRateLimiterFlags int
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgAddOutboundTracker, &weightMsgAddOutboundTracker, nil,
		func(_ *rand.Rand) {
			weightMsgAddOutboundTracker = DefaultWeightMsgAddOutboundTracker
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

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightVoteGasPrice,
			SimulateMsgVoteGasPrice(k),
		),
		simulation.NewWeightedOperation(
			weightVoteInbound,
			SimulateVoteInbound(k),
		),
	}
}

// operationSimulateVoteInbound generates a MsgVoteInbound with a random vote and delivers it.
func operationSimulateVoteInbound(
	k keeper.Keeper,
	msg types.MsgVoteInbound,
	simAccount simtypes.Account,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, _ []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		// Fetch the account from the auth keeper which can then be used to fetch spendable coins
		authAccount := k.GetAuthKeeper().GetAccount(ctx, simAccount.Address)
		spendable := k.GetBankKeeper().SpendableCoins(ctx, authAccount.GetAddress())

		// Generate a transaction with a random fee and deliver it
		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           moduletestutil.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             &msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   k.GetAuthKeeper(),
			Bankkeeper:      k.GetBankKeeper(),
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		// Generate and deliver the transaction using the function defined by us instead of using the default function provided by the cosmos-sdk
		// The main difference between the two functions is that the one defined by us does not error out if the vote fails.
		// We need this behaviour as the votes are assigned to future operations, i.e., they are scheduled to be executed in a future block. We do not know at the time of scheduling if the vote will be successful or not.
		// There might be multiple reasons for a vote to fail , like the observer not being present in the observer set, the observer not being an observer, etc.
		return GenAndDeliverTxWithRandFees(txCtx)
	}
}

func SimulateVoteInbound(k keeper.Keeper) simtypes.Operation {
	// The states are:
	// column 1: All observers vote
	// column 2: 90% vote
	// column 3: 75% vote
	// column 4: 40% vote
	// column 5: 15% vote
	// column 6: noone votes
	// All columns sum to 100 for simplicity, but this is arbitrary and can be changed
	numVotesTransitionMatrix, _ := simulation.CreateTransitionMatrix([][]int{
		{20, 10, 0, 0, 0, 0},
		{55, 50, 20, 10, 0, 0},
		{25, 25, 30, 25, 30, 15},
		{0, 15, 30, 25, 30, 30},
		{0, 0, 20, 30, 30, 30},
		{0, 0, 0, 10, 10, 25},
	})

	statePercentageArray := []float64{1, .9, .75, .4, .15, 0}
	curNumVotesState := 1

	return func(
		r *rand.Rand,
		app *baseapp.BaseApp,
		ctx sdk.Context,
		accs []simtypes.Account,
		chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// TODO : randomize these values
		// Right now we use a constant value for cctx creation , this is the same as the one used in unit tests for the successful condition.
		// TestKeeper_VoteInbound/successfully vote on evm deposit
		// But this can improved by adding more randomization

		//https://github.com/zeta-chain/node/issues/3101
		to, from := int64(1337), int64(101)
		supportedChains := k.GetObserverKeeper().GetSupportedChains(ctx)
		for _, chain := range supportedChains {
			if chains.IsEVMChain(chain.ChainId, []chains.Chain{}) {
				from = chain.ChainId
			}
			if chains.IsZetaChain(chain.ChainId, []chains.Chain{}) {
				to = chain.ChainId
			}
		}
		msg := sample.InboundVoteSim(0, from, to, r)

		// Pick a random observer to create the ballot
		// If this returns an error, it is likely that the entire observer set has been removed
		simAccount, firstVoter, err := GetRandomAccountAndObserver(r, ctx, k, accs)
		if err != nil {
			return simtypes.OperationMsg{}, nil, nil
		}

		txGen := moduletestutil.MakeTestEncodingConfig().TxConfig
		account := k.GetAuthKeeper().GetAccount(ctx, simAccount.Address)
		firstMsg := msg
		firstMsg.Creator = firstVoter

		err = firstMsg.ValidateBasic()
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to validate first inbound vote"), nil, err
		}

		tx, err := simtestutil.GenSignedMockTx(
			r,
			txGen,
			[]sdk.Msg{&firstMsg},
			sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)},
			simtestutil.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		// We can return error here as we  can guarantee that the first vote will be successful.
		// Since we query the observer set before adding votes
		_, _, err = app.SimDeliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		opMsg := simtypes.NewOperationMsg(&msg, true, "", nil)

		// Add subsequent votes
		observerSet, found := k.GetObserverKeeper().GetObserverSet(ctx)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, authz.InboundVoter.String(), "observer set not found"), nil, nil
		}

		// 1) Schedule operations for votes
		// 1.1) first pick a number of people to vote.
		curNumVotesState = numVotesTransitionMatrix.NextState(r, curNumVotesState)
		numVotes := int(math.Ceil(float64(len(observerSet.ObserverList)) * statePercentageArray[curNumVotesState]))

		// 1.2) select who votes
		whoVotes := r.Perm(len(observerSet.ObserverList))
		whoVotes = whoVotes[:numVotes]

		var fops []simtypes.FutureOperation

		for _, observerIdx := range whoVotes {
			observerAddress := observerSet.ObserverList[observerIdx]
			// firstVoter has already voted.
			if observerAddress == firstVoter {
				continue
			}
			observerAccount, err := GetObserverAccount(observerAddress, accs)
			if err != nil {
				continue
			}
			// 1.3) schedule the vote
			votingMsg := msg
			votingMsg.Creator = observerAddress

			e := votingMsg.ValidateBasic()
			if e != nil {
				return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to validate voting msg"), nil, e
			}

			fops = append(fops, simtypes.FutureOperation{
				// Submit all subsequent votes in the next block.
				// We can consider adding a random block height between 1 and ballot maturity blocks in the future.
				BlockHeight: int(ctx.BlockHeight() + 1),
				Op:          operationSimulateVoteInbound(k, votingMsg, observerAccount),
			})
		}
		return opMsg, fops, nil
	}
}

// SimulateMsgVoteGasPrice generates a MsgVoteGasPrice and delivers it
func SimulateMsgVoteGasPrice(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		// Get a random account and observer
		// If this returns an error, it is likely that the entire observer set has been removed
		simAccount, randomObserver, err := GetRandomAccountAndObserver(r, ctx, k, accounts)
		if err != nil {
			return simtypes.OperationMsg{}, nil, nil
		}
		authAccount := k.GetAuthKeeper().GetAccount(ctx, simAccount.Address)
		spendable := k.GetBankKeeper().SpendableCoins(ctx, authAccount.GetAddress())

		supportedChains := k.GetObserverKeeper().GetSupportedChains(ctx)
		if len(supportedChains) == 0 {
			return simtypes.NoOpMsg(
				types.ModuleName,
				authz.GasPriceVoter.String(),
				"no supported chains found",
			), nil, nil
		}
		randomChainID := GetRandomChainID(r, supportedChains)

		// Vote for random gas price. Gas prices do not use a ballot system, so we can vote directly without having to schedule future operations.
		// The random nature of the price might create weird gas prices for the chain, but it is fine for now. We can remove the randomness if needed
		msg := types.MsgVoteGasPrice{
			Creator:     randomObserver,
			ChainId:     randomChainID,
			Price:       r.Uint64(),
			PriorityFee: r.Uint64(),
			BlockNumber: r.Uint64(),
			Supply:      fmt.Sprintf("%d", r.Int63()),
		}

		// System contracts are deployed on the first block, so we cannot vote on gas prices before that
		if ctx.BlockHeight() <= 1 {
			return simtypes.NewOperationMsg(&msg, true, "block height less than 1", nil), nil, nil
		}

		err = msg.ValidateBasic()
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to validate vote gas price  msg"), nil, err
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           moduletestutil.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             &msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   k.GetAuthKeeper(),
			Bankkeeper:      k.GetBankKeeper(),
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
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

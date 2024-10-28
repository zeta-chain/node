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

const (
	DefaultWeightMsgAddOutboundTracker  = 50
	DefaultWeightAddInboundTracker      = 50
	DefaultWeightRemoveOutboundTracker  = 5
	DefaultWeightVoteGasPrice           = 100
	DefaultWeightVoteOutbound           = 50
	DefaultWeightVoteInbound            = 10
	DefaultWeightWhitelistERC20         = 1
	DefaultWeightMigrateTssFunds        = 1
	DefaultWeightUpdateTssAddress       = 1
	DefaultWeightAbortStuckCCTX         = 10
	DefaultWeightUpdateRateLimiterFlags = 1

	OpWeightMsgAddOutboundTracker  = "op_weight_msg_add_outbound_tracker"
	OpWeightAddInboundTracker      = "op_weight_msg_add_inbound_tracker"
	OpWeightRemoveOutboundTracker  = "op_weight_msg_remove_outbound_tracker"
	OpWeightVoteGasPrice           = "op_weight_msg_vote_gas_price"
	OpWeightVoteOutbound           = "op_weight_msg_vote_outbound"
	OpWeightVoteInbound            = "op_weight_msg_vote_inbound"
	OpWeightWhitelistERC20         = "op_weight_msg_whitelist_erc20"
	OpWeightMigrateTssFunds        = "op_weight_msg_migrate_tss_funds"
	OpWeightUpdateTssAddress       = "op_weight_msg_update_tss_address"
	OpWeightAbortStuckCCTX         = "op_weight_msg_abort_stuck_cctx"
	OpWeightUpdateRateLimiterFlags = "op_weight_msg_update_rate_limiter_flags"
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
		//simulation.NewWeightedOperation(
		//	weightVoteGasPrice,
		//	SimulateMsgVoteGasPrice(k),
		//),
		simulation.NewWeightedOperation(
			weightVoteInbound,
			SimulateVoteInbound(k),
		),
	}
}

func operationSimulateVoteInbound(k keeper.Keeper, msg types.MsgVoteInbound, simAccount simtypes.Account) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, chainID string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {

		account := k.GetAuthKeeper().GetAccount(ctx, simAccount.Address)
		spendable := k.GetBankKeeper().SpendableCoins(ctx, account.GetAddress())

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
	// All columns sum to 100 for simplicity, but this is arbitrary and
	// feel free to change.
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

		//var (
		//	randomChainSender   = GetRandomChainID(r, supportedChains)
		//	randomChainReceiver = GetRandomChainID(r, supportedChains)
		//	sender              = sample.EthAddress()
		//	receiver            = sample.EthAddress()
		//	amount              = r.Uint64()
		//	coinType            = coin.CoinType_Gas
		//	hash                = sample.Hash().String()
		//	gasLimit            = r.Uint64()
		//)
		// TODO : randomize these values
		to, from := int64(101), int64(101)
		msg := sample.InboundVote(0, from, to)

		simAccount, firstVoter, err := GetRandomAccountAndObserver(r, ctx, k, accs)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, authz.InboundVoter.String(), "unable to get random account and observer"), nil, nil
		}

		txGen := moduletestutil.MakeTestEncodingConfig().TxConfig
		account := k.GetAuthKeeper().GetAccount(ctx, simAccount.Address)
		firstMsg := msg
		firstMsg.Creator = firstVoter

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

		// We can return error here as we we can guarantee that the first vote will be successful. Since we query the observer set before adding votes
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

		// 1.2) select who votes and when
		whoVotes := r.Perm(len(observerSet.ObserverList))
		// didntVote := whoVotes[numVotes:]
		whoVotes = whoVotes[:numVotes]
		//ballotMaturityPeriod := 20
		var fops []simtypes.FutureOperation

		//fmt.Printf("\nScheduling %d votes for ballot : %s at height : %d", numVotes, msg.Digest(), ctx.BlockHeight()+1)
		//if numVotes > 3 {
		//	fmt.Printf("More than 3 votes : %s", msg.Digest())
		//}
		for _, observerIdx := range whoVotes {
			observerAddress := observerSet.ObserverList[observerIdx]
			if observerAddress == firstVoter {
				continue
			}
			observerAccount, err := GetObserverAccount(observerAddress, accs)
			if err != nil {
				panic(err)
			}
			// 1.3) schedule the vote
			votingMsg := msg
			votingMsg.Creator = observerAddress
			fops = append(fops, simtypes.FutureOperation{
				BlockHeight: int(ctx.BlockHeight() + 1),
				Op:          operationSimulateVoteInbound(k, votingMsg, observerAccount),
			})
		}
		//fmt.Println("\nActual votes : ", len(fops))

		return opMsg, fops, nil
	}
}

func SimulateMsgVoteGasPrice(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, chainID string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {

		simAccount, randomObserver, err := GetRandomAccountAndObserver(r, ctx, k, accounts)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, authz.GasPriceVoter.String(), "unable to get random account and observer"), nil, err
		}

		supportedChains := k.GetObserverKeeper().GetSupportedChains(ctx)
		if len(supportedChains) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, authz.GasPriceVoter.String(), "no supported chains found"), nil, nil
		}
		randomChainID := GetRandomChainID(r, supportedChains)

		msg := types.MsgVoteGasPrice{
			Creator:     randomObserver,
			ChainId:     randomChainID,
			Price:       r.Uint64(),
			PriorityFee: r.Uint64(),
			BlockNumber: r.Uint64(),
			Supply:      fmt.Sprintf("%d", r.Int63()),
		}

		txCtx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         moduletestutil.MakeTestEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           &msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    simAccount,
			AccountKeeper: k.GetAuthKeeper(),
			Bankkeeper:    k.GetBankKeeper(),
			ModuleName:    types.ModuleName,
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

func GetRandomAccountAndObserver(r *rand.Rand, ctx sdk.Context, k keeper.Keeper, accounts []simtypes.Account) (simtypes.Account, string, error) {
	observers, found := k.GetObserverKeeper().GetObserverSet(ctx)
	if !found {
		return simtypes.Account{}, "", fmt.Errorf("observer set not found")
	}

	if len(observers.ObserverList) == 0 {
		return simtypes.Account{}, "", fmt.Errorf("no observers present in observer set found")
	}

	randomObserver := GetRandomObserver(r, observers.ObserverList)

	// TODO : use GetObserverAccount
	operatorAddress, err := observerTypes.GetOperatorAddressFromAccAddress(randomObserver)
	if err != nil {
		return simtypes.Account{}, "", fmt.Errorf("validator not found for observer ")
	}

	simAccount, found := simtypes.FindAccount(accounts, operatorAddress)
	if !found {
		return simtypes.Account{}, "", fmt.Errorf("operator account not found")
	}
	return simAccount, randomObserver, nil
}

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
func GenAndDeliverTxWithRandFees(txCtx simulation.OperationInput) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
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

// GenAndDeliverTx generates a transactions and delivers it.
func GenAndDeliverTx(txCtx simulation.OperationInput, fees sdk.Coins) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
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

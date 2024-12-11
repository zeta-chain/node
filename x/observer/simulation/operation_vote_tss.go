package simulation

import (
	"math"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

func operationSimulateVoteTss(
	k keeper.Keeper,
	msg types.MsgVoteTSS,
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

// SimulateVoteOutbound generates a MsgVoteOutbound with random values
// This is the only operation which saves a cctx directly to the store.
func SimulateMsgVoteTSS(k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand,
		app *baseapp.BaseApp,
		ctx sdk.Context,
		accs []simtypes.Account,
		chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		yesVote := chains.ReceiveStatus_success
		noVote := chains.ReceiveStatus_failed
		ballotVotesTransitionMatrix, yesVotePercentageArray, ballotVotesState := BallotVoteSimulationMatrix()
		nodeAccounts := k.GetAllNodeAccount(ctx)
		numVotes := len(nodeAccounts)
		ballotVotesState = ballotVotesTransitionMatrix.NextState(r, ballotVotesState)
		yesVotePercentage := yesVotePercentageArray[ballotVotesState]
		numberOfYesVotes := int(math.Ceil(float64(numVotes) * yesVotePercentage))

		vote := yesVote
		if numberOfYesVotes == 0 {
			vote = noVote
		}

		newTss, err := sample.TSSFromRand(r)
		if err != nil {
			return simtypes.OperationMsg{}, nil, err
		}

		keygen, found := k.GetKeygen(ctx)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgVoteTSS, "keygen not found"), nil, nil
		}

		msg := types.MsgVoteTSS{
			Creator:          "",
			TssPubkey:        newTss.TssPubkey,
			KeygenZetaHeight: keygen.BlockNumber,
			Status:           vote,
		}

		// Pick a random observer to create the ballot
		// If this returns an error, it is likely that the entire observer set has been removed
		simAccount, firstVoter, err := GetRandomNodeAccount(r, ctx, k, accs)
		if err != nil {
			return simtypes.OperationMsg{}, nil, nil
		}

		txGen := moduletestutil.MakeTestEncodingConfig().TxConfig
		account := k.GetAuthKeeper().GetAccount(ctx, simAccount.Address)

		firstMsg := msg
		firstMsg.Creator = firstVoter

		// THe first vote should always create a new ballot
		_, found = k.GetBallot(ctx, firstMsg.Digest())
		if found {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "ballot already exists"), nil, nil
		}

		err = firstMsg.ValidateBasic()
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to validate first tss vote"), nil, err
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

		var fops []simtypes.FutureOperation

		for voteCount, nodeAccount := range nodeAccounts {
			if vote == yesVote && voteCount == numberOfYesVotes {
				vote = noVote
			}
			// firstVoter has already voted.
			if nodeAccount.Operator == firstVoter {
				continue
			}
			observerAccount, err := GetSimAccount(nodeAccount.Operator, accs)
			if err != nil {
				continue
			}
			// 1.3) schedule the vote
			votingMsg := msg
			votingMsg.Creator = nodeAccount.Operator
			votingMsg.Status = vote

			e := votingMsg.ValidateBasic()
			if e != nil {
				return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to validate voting msg"), nil, e
			}

			fops = append(fops, simtypes.FutureOperation{
				// Submit all subsequent votes in the next block.
				// We can consider adding a random block height between 1 and ballot maturity blocks in the future.
				BlockHeight: int(ctx.BlockHeight() + 1),
				Op:          operationSimulateVoteTss(k, votingMsg, observerAccount),
			})
		}
		return opMsg, fops, nil
	}
}

package simulation

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	zetasimulation "github.com/zeta-chain/node/testutil/simulation"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func operationSimulateVoteOutbound(
	k keeper.Keeper,
	msg types.MsgVoteOutbound,
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
		return zetasimulation.GenAndDeliverTxWithRandFees(txCtx, false)
	}
}

// SimulateVoteOutbound generates a MsgVoteOutbound with random values and delivers it.
//It also schedules future operations for subsequent votes.

func SimulateVoteOutbound(k keeper.Keeper) simtypes.Operation {
	defaultVote := chains.ReceiveStatus_success
	alternativeVote := chains.ReceiveStatus_failed
	observerVotesTransitionMatrix, statePercentageArray, curNumVotesState := zetasimulation.ObserverVotesSimulationMatrix()
	ballotVotesTransitionMatrix, yesVotePercentageArray, ballotVotesState := zetasimulation.OutboundVoteStatusSimulationMatrix()
	return func(
		r *rand.Rand,
		app *baseapp.BaseApp,
		ctx sdk.Context,
		accs []simtypes.Account,
		chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		to, from := chains.GoerliLocalnet.ChainId, chains.ZetaChainPrivnet.ChainId
		supportedChains := k.GetObserverKeeper().GetSupportedChains(ctx)
		for _, chain := range supportedChains {
			if chains.IsEthereumChain(chain.ChainId, []chains.Chain{}) {
				to = chain.ChainId
			}
			if chains.IsZetaChain(chain.ChainId, []chains.Chain{}) {
				from = chain.ChainId
			}
		}

		_, creator, _, err := zetasimulation.GetRandomAccountAndObserver(r, ctx, k.GetObserverKeeper(), accs)
		if err != nil {
			return simtypes.OperationMsg{}, nil, nil
		}
		index := ethcrypto.Keccak256Hash(fmt.Appendf(nil, "%d", r.Int63())).Hex()

		tss, found := k.GetObserverKeeper().GetTSS(ctx)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgVoteOutbound, "tss not found"), nil, nil
		}

		asset, err := zetasimulation.GetAsset(ctx, k.GetFungibleKeeper(), to)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgVoteOutbound, "unable to get asset"), nil, err
		}

		// Generate a new cctx and save it , which can be used to finalize the outbound
		cctx := sample.CCTXfromRand(r, creator, index, to, from, tss.TssPubkey, asset)
		msg := types.MsgVoteOutbound{
			CctxHash:                          cctx.Index,
			OutboundTssNonce:                  cctx.GetCurrentOutboundParam().TssNonce,
			OutboundChain:                     cctx.GetCurrentOutboundParam().ReceiverChainId,
			Status:                            defaultVote,
			Creator:                           cctx.Creator,
			ObservedOutboundHash:              ethcommon.BytesToHash(sample.EthAddressFromRand(r).Bytes()).String(),
			ValueReceived:                     cctx.GetCurrentOutboundParam().Amount,
			ObservedOutboundBlockHeight:       cctx.GetCurrentOutboundParam().ObservedExternalHeight,
			ObservedOutboundEffectiveGasPrice: cctx.GetCurrentOutboundParam().EffectiveGasPrice,
			ObservedOutboundGasUsed:           cctx.GetCurrentOutboundParam().GasUsed,
			CoinType:                          cctx.InboundParams.CoinType,
			ConfirmationMode:                  cctx.GetCurrentOutboundParam().ConfirmationMode,
		}

		err = k.SetObserverOutboundInfo(ctx, to, &cctx)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgVoteOutbound,
				"unable to set observer outbound info",
			), nil, err
		}

		msg.OutboundTssNonce = cctx.GetCurrentOutboundParam().TssNonce
		k.SaveCCTXUpdate(ctx, cctx, tss.TssPubkey)

		// Pick a random observer to create the ballot
		// If this returns an error, it is likely that the entire observer set has been removed
		simAccount, firstVoter, _, err := zetasimulation.GetRandomAccountAndObserver(
			r,
			ctx,
			k.GetObserverKeeper(),
			accs,
		)
		if err != nil {
			return simtypes.OperationMsg{}, nil, nil
		}

		txGen := moduletestutil.MakeTestEncodingConfig().TxConfig
		account := k.GetAuthKeeper().GetAccount(ctx, simAccount.Address)
		firstMsg := msg
		firstMsg.Creator = firstVoter

		// THe first vote should always create a new ballot
		_, found = k.GetObserverKeeper().GetBallot(ctx, firstMsg.Digest())
		if found {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgVoteOutbound, "ballot already exists"), nil, nil
		}

		err = firstMsg.ValidateBasic()
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgVoteOutbound,
				"unable to validate first outbound vote",
			), nil, err
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
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgVoteOutbound, "unable to generate mock tx"), nil, err
		}

		// We can return error here as we can guarantee that the first vote will be successful.
		// Since we query the observer set before adding votes
		_, _, err = app.SimDeliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgVoteOutbound, "unable to deliver tx"), nil, err
		}

		opMsg := zetasimulation.OperationMessage(&msg)

		// Add subsequent votes
		observerSet, found := k.GetObserverKeeper().GetObserverSet(ctx)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgVoteOutbound, "observer set not found"), nil, nil
		}

		// 1) Schedule operations for votes
		// 1.1) first pick a number of people to vote.
		curNumVotesState = observerVotesTransitionMatrix.NextState(r, curNumVotesState)
		numVotes := int(math.Ceil(float64(len(observerSet.ObserverList)) * statePercentageArray[curNumVotesState]))

		// 1.2) select who votes
		whoVotes := r.Perm(len(observerSet.ObserverList))
		whoVotes = whoVotes[:numVotes]

		var fops []simtypes.FutureOperation

		ballotVotesState = ballotVotesTransitionMatrix.NextState(r, ballotVotesState)
		yesVotePercentage := yesVotePercentageArray[ballotVotesState]
		numberOfYesVotes := int(math.Ceil(float64(numVotes) * yesVotePercentage))
		vote := defaultVote

		for voteCount, observerIdx := range whoVotes {
			if voteCount == numberOfYesVotes {
				vote = alternativeVote
			}
			observerAddress := observerSet.ObserverList[observerIdx]
			// firstVoter has already voted.
			if observerAddress == firstVoter {
				continue
			}
			observerAccount, err := zetasimulation.GetObserverAccount(observerAddress, accs)
			if err != nil {
				continue
			}
			// 1.3) schedule the vote
			votingMsg := msg
			votingMsg.Creator = observerAddress
			votingMsg.Status = vote

			e := votingMsg.ValidateBasic()
			if e != nil {
				return simtypes.NoOpMsg(types.ModuleName, TypeMsgVoteOutbound, "unable to validate voting msg"), nil, e
			}

			fops = append(fops, simtypes.FutureOperation{
				// Submit all subsequent votes in the next block.
				// We can consider adding a random block height between 1 and ballot maturity blocks in the future.
				BlockHeight: int(ctx.BlockHeight() + 1),
				Op:          operationSimulateVoteOutbound(k, votingMsg, observerAccount),
			})
		}
		return opMsg, fops, nil
	}
}

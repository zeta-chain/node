package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

// SimulateAddObserver generates a TypeMsgAddObserver and delivers it. This message sets AddNodeAccountOnly to false;
// Therefore, it adds the observer to the observer set
func SimulateAddObserver(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		policyAccount, err := GetPolicyAccount(ctx, k.GetAuthorityKeeper(), accounts)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgAddObserver, err.Error()), nil, nil
		}

		authAccount := k.GetAuthKeeper().GetAccount(ctx, policyAccount.Address)
		spendable := k.GetBankKeeper().SpendableCoins(ctx, authAccount.GetAddress())

		observerSet, found := k.GetObserverSet(ctx)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgAddObserver,
				"no observer set found",
			), nil, nil
		}

		observerMap := make(map[string]bool)
		for _, observer := range observerSet.ObserverList {
			observerMap[observer] = true
		}

		nodeAccounts := k.GetAllNodeAccount(ctx)

		// Pick a random observer which part of the node account but not in the observer set
		// New accounts are added to the node account list via SimulateAddObserverNodeAccount
		var newObserver string
		foundNA := RepeatCheck(func() bool {
			newObserver = nodeAccounts[r.Intn(len(nodeAccounts))].Operator
			if _, found := observerMap[newObserver]; !found {
				return true
			}
			return false
		})
		if !foundNA {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgAddObserver,
				"no node accounts available which can be added as observer",
			), nil, nil
		}
		pubkey, err := sample.PubkeyStringFromRand(r)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgAddObserver, err.Error()), nil, nil
		}
		msg := types.MsgAddObserver{
			Creator:                 policyAccount.Address.String(),
			ObserverAddress:         newObserver,
			ZetaclientGranteePubkey: pubkey,
			AddNodeAccountOnly:      false,
		}

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, err
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           moduletestutil.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             &msg,
			Context:         ctx,
			SimAccount:      policyAccount,
			AccountKeeper:   k.GetAuthKeeper(),
			Bankkeeper:      k.GetBankKeeper(),
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

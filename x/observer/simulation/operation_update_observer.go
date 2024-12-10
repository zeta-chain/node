package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/zeta-chain/node/x/observer/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

// SimulateMsgUpdateObserver generates a TypeMsgUpdateObserver and delivers it.
func SimulateMsgUpdateObserver(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		policyAccount, err := GetPolicyAccount(ctx, k.GetAuthorityKeeper(), accounts)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateObserver, err.Error()), nil, nil
		}

		authAccount := k.GetAuthKeeper().GetAccount(ctx, policyAccount.Address)
		spendable := k.GetBankKeeper().SpendableCoins(ctx, authAccount.GetAddress())

		_, randomObserver, err, observerList := GetRandomAccountAndObserver(r, ctx, k, accounts)

		observerMap := make(map[string]bool)
		for _, observer := range observerList {
			observerMap[observer] = true
		}

		validators := k.GetStakingKeeper().GetAllValidators(ctx)
		if len(validators) == 0 {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgUpdateObserver,
				"no validators found",
			), nil, nil
		}
		newObserver := ""
		for {
			randomValidator := validators[r.Intn(len(validators))]
			nO, err := types.GetAccAddressFromOperatorAddress(randomValidator.OperatorAddress)
			if err != nil {
				continue
			}
			newObserver = nO.String()
			err = k.IsValidator(ctx, newObserver)
			if err != nil {
				continue
			}
			if _, ok := observerMap[newObserver]; !ok {
				break
			}

		}

		lastBlockCount, found := k.GetLastObserverCount(ctx)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgUpdateObserver,
				"no last block count found",
			), nil, nil
		}
		if int(lastBlockCount.Count) != len(observerList) {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgUpdateObserver,
				"observer count mismatch",
			), nil, nil
		}

		msg := types.MsgUpdateObserver{
			Creator:            policyAccount.Address.String(),
			OldObserverAddress: randomObserver,
			NewObserverAddress: newObserver,
			UpdateReason:       types.ObserverUpdateReason_AdminUpdate,
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
			MsgType:         msg.Type(),
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

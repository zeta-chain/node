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

// SimulateUpdateObserver generates a TypeMsgUpdateObserver and delivers it.
func SimulateUpdateObserver(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		policyAccount, err := GetPolicyAccount(ctx, k.GetAuthorityKeeper(), accounts)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateObserver, err.Error()), nil, nil
		}

		authAccount := k.GetAuthKeeper().GetAccount(ctx, policyAccount.Address)
		spendable := k.GetBankKeeper().SpendableCoins(ctx, authAccount.GetAddress())

		_, randomObserver, observerList, err := GetRandomAccountAndObserver(r, ctx, k, accounts)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateObserver, err.Error()), nil, nil
		}

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
		foundNewObserver := RepeatCheck(func() bool {
			randomValidator := validators[r.Intn(len(validators))]
			randomValidatorAddress, err := types.GetAccAddressFromOperatorAddress(randomValidator.OperatorAddress)
			if err != nil {
				return false
			}
			newObserver = randomValidatorAddress.String()
			err = k.IsValidator(ctx, newObserver)
			if err != nil {
				return false
			}
			if _, ok := observerMap[newObserver]; !ok {
				return true
			}
			return false
		})

		if !foundNewObserver {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgUpdateObserver,
				"no new observer found",
			), nil, nil
		}

		lastBlockCount, found := k.GetLastObserverCount(ctx)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgUpdateObserver,
				"no last block count found",
			), nil, nil
		}
		// #nosec G115 - overflow is not a concern here
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

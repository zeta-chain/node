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

// SimulateMsgAddObserverNodeAccount generates a TypeMsgAddObserver and delivers it.
func SimulateMsgAddObserverNodeAccount(k keeper.Keeper) simtypes.Operation {
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
			randomValidatorOperatorAddress, err := types.GetAccAddressFromOperatorAddress(
				randomValidator.OperatorAddress,
			)
			if err != nil {
				continue
			}
			newObserver = randomValidatorOperatorAddress.String()
			err = k.IsValidator(ctx, newObserver)
			if err != nil {
				continue
			}
			if _, ok := observerMap[newObserver]; !ok {
				break
			}
		}

		msg := types.MsgAddObserver{
			Creator:                 policyAccount.Address.String(),
			ObserverAddress:         newObserver,
			ZetaclientGranteePubkey: sample.PubkeyStringFromRand(r),
			AddNodeAccountOnly:      true,
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

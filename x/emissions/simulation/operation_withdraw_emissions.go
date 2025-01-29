package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	zetasimulation "github.com/zeta-chain/node/testutil/simulation"
	"github.com/zeta-chain/node/x/emissions/keeper"
	"github.com/zeta-chain/node/x/emissions/types"
)

// SimulateMsgWithdrawEmissions generates a MsgWithdrawEmission with a random amount of emissions to withdraw
func SimulateMsgWithdrawEmissions(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		// Get a random account and observer
		// If this returns an error, it is likely that the entire observer set has been removed
		simAccount, randomObserver, _, err := zetasimulation.GetRandomAccountAndObserver(
			r,
			ctx,
			k.GetObserverKeeper(),
			accounts,
		)
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgWithdrawEmission,
				"cant to fetch a observer account",
			), nil, err
		}
		authAccount := k.GetAuthKeeper().GetAccount(ctx, simAccount.Address)
		spendable := k.GetBankKeeper().SpendableCoins(ctx, authAccount.GetAddress())

		availableEmissions, found := k.GetWithdrawableEmission(ctx, randomObserver)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgWithdrawEmission,
				"no emissions found",
			), nil, nil
		}

		if availableEmissions.Amount.IsZero() {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgWithdrawEmission,
				"no emissions available to withdraw",
			), nil, nil
		}

		// Pick a random amount of emissions to withdraw between 0 and availableEmissions
		amount, err := simtypes.RandPositiveInt(r, availableEmissions.Amount)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName,
				TypeMsgWithdrawEmission,
				"unable to generate random amount",
			), nil, err
		}

		// Withdraw emissions
		msg := types.MsgWithdrawEmission{
			Creator: randomObserver,
			Amount:  amount,
		}

		err = msg.ValidateBasic()
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgWithdrawEmission,
				"unable to validate MsgWithdrawEmission",
			), nil, err
		}

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

		return zetasimulation.GenAndDeliverTxWithRandFees(txCtx, true)
	}
}

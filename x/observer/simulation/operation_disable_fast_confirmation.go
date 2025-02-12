package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	zetasimulation "github.com/zeta-chain/node/testutil/simulation"
	"github.com/zeta-chain/node/x/observer/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

// SimulateDisableFastConfirmation generates a MsgDisableFastConfirmation and delivers it.
func SimulateDisableFastConfirmation(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		policyAccount, err := zetasimulation.GetPolicyAccount(ctx, k.GetAuthorityKeeper(), accounts)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgDisableFastConfirmation, err.Error()), nil, nil
		}
		authAccount := k.GetAuthKeeper().GetAccount(ctx, policyAccount.Address)
		spendable := k.GetBankKeeper().SpendableCoins(ctx, authAccount.GetAddress())

		supportedChains := k.GetSupportedChains(ctx)
		if len(supportedChains) == 0 {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgDisableFastConfirmation,
				"no supported chains found",
			), nil, nil
		}

		// Pick a random external chain
		randomExternalChain := int64(0)
		foundExternalChain := zetasimulation.RepeatCheck(func() bool {
			c := supportedChains[r.Intn(len(supportedChains))]
			if !c.IsZetaChain() {
				randomExternalChain = c.ChainId
				return true
			}
			return false
		})

		if !foundExternalChain {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgDisableFastConfirmation,
				"no external chain found",
			), nil, nil
		}

		msg := types.MsgDisableFastConfirmation{
			Creator: policyAccount.Address.String(),
			ChainId: randomExternalChain,
		}

		err = msg.ValidateBasic()
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgDisableFastConfirmation, err.Error()), nil, err
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

		return zetasimulation.GenAndDeliverTxWithRandFees(txCtx, true)
	}
}

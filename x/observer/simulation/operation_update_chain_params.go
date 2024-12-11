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

// SimulateUpdateChainParams generates a MsgUpdateChainParams and delivers it.
func SimulateUpdateChainParams(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		policyAccount, err := GetPolicyAccount(ctx, k.GetAuthorityKeeper(), accounts)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateChainParams, err.Error()), nil, nil
		}

		authAccount := k.GetAuthKeeper().GetAccount(ctx, policyAccount.Address)
		spendable := k.GetBankKeeper().SpendableCoins(ctx, authAccount.GetAddress())

		randomChain, err := GetExternalChain(ctx, k, r)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateChainParams, err.Error()), nil, nil
		}

		cp := sample.ChainParamsFromRand(r, randomChain.ChainId)

		msg := types.MsgUpdateChainParams{
			Creator:     policyAccount.Address.String(),
			ChainParams: cp,
		}

		err = msg.ValidateBasic()
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

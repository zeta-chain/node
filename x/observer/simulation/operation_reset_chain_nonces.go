package simulation

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/zeta-chain/node/x/observer/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

// SimulateResetChainNonces generates a MsgResetChainNonces and delivers it.
func SimulateResetChainNonces(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		policyAccount, err := GetPolicyAccount(ctx, k.GetAuthorityKeeper(), accounts)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgResetChainNonces, err.Error()), nil, nil
		}

		authAccount := k.GetAuthKeeper().GetAccount(ctx, policyAccount.Address)
		spendable := k.GetBankKeeper().SpendableCoins(ctx, authAccount.GetAddress())

		randomChain, err := GetExternalChain(ctx, k, r)
		if err != nil {
			return simtypes.NoOpMsg(
					types.ModuleName,
					types.TypeMsgResetChainNonces,
					err.Error(),
				), nil, fmt.Errorf(
					"error getting external chain",
				)
		}

		tss, found := k.GetTSS(ctx)
		if !found {
			return simtypes.NoOpMsg(
					types.ModuleName,
					types.TypeMsgResetChainNonces,
					"TSS not found",
				), nil, fmt.Errorf(
					"TSS not found",
				)
		}
		pendingNonces, found := k.GetPendingNonces(ctx, tss.TssPubkey, randomChain.ChainId)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgResetChainNonces, "Pending nonces not found"), nil,
				fmt.Errorf("pending nonces not found for chain %d %s", randomChain.ChainId, randomChain.ChainName)
		}

		nonceIncrement := int64(r.Intn(99)) + 1

		msg := types.MsgResetChainNonces{
			Creator:        policyAccount.Address.String(),
			ChainId:        randomChain.ChainId,
			ChainNonceHigh: pendingNonces.NonceHigh + nonceIncrement,
			ChainNonceLow:  pendingNonces.NonceLow + nonceIncrement,
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

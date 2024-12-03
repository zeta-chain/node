package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// SimulateMsgRefundAbortedCCTX generates a MsgRefundAbortedCCTX with random values
func SimulateMsgRefundAbortedCCTX(k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		// Fetch the account from the auth keeper which can then be used to fetch spendable coins}
		policyAccount, err := GetPolicyAccount(ctx, k.GetAuthorityKeeper(), accounts)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgAbortStuckCCTX, err.Error()), nil, nil
		}

		authAccount := k.GetAuthKeeper().GetAccount(ctx, policyAccount.Address)
		spendable := k.GetBankKeeper().SpendableCoins(ctx, authAccount.GetAddress())

		cctxList := k.GetAllCrossChainTx(ctx)
		abortedCctx := types.CrossChainTx{}
		abortedCctxFound := false

		for _, cctx := range cctxList {
			if cctx.CctxStatus.Status == types.CctxStatus_Aborted {
				abortedCctx = cctx
				abortedCctxFound = true
				break
			}
		}
		if !abortedCctxFound {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgAbortStuckCCTX, "no aborted cctx found"), nil, nil
		}

		if abortedCctx.CctxStatus.IsAbortRefunded {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgAbortStuckCCTX, "aborted cctx already refunded"), nil, nil
		}

		msg := types.MsgRefundAbortedCCTX{
			Creator:       policyAccount.Address.String(),
			CctxIndex:     abortedCctx.Index,
			RefundAddress: sample.EthAddressFromRand(r).String(),
		}

		err = msg.ValidateBasic()
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to validate MsgRefundAbortedCCTX msg"), nil, err
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

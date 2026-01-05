package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/zeta-chain/node/testutil/sample"
	zetasimulation "github.com/zeta-chain/node/testutil/simulation"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// SimulateMsgUpdateTssAddress generates a MsgUpdateTssAddress with random values and delivers it
func SimulateMsgUpdateTssAddress(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		policyAccount, err := zetasimulation.GetPolicyAccount(ctx, k.GetAuthorityKeeper(), accounts)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgUpdateTssAddress, err.Error()), nil, nil
		}

		authAccount := k.GetAuthKeeper().GetAccount(ctx, policyAccount.Address)
		spendable := k.GetBankKeeper().SpendableCoins(ctx, authAccount.GetAddress())

		supportedChains := k.TSSFundsMigrationChains(ctx)
		if len(supportedChains) == 0 {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgUpdateTssAddress,
				"no chains found which support tss migration",
			), nil, nil
		}

		cctxList := k.GetAllCrossChainTx(ctx)
		if len(cctxList) == 0 {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgUpdateTssAddress,
				"no cross chain txs found",
			), nil, nil
		}

		// Pick any cctx with status OutboundMined, and use its index for the migration
		// We set the fund migrator directly as we are not simulating MsgMigrateTssFunds
		minedCCTX := types.CrossChainTx{}
		foundMined := false
		for _, cctx := range cctxList {
			if cctx.CctxStatus.Status == types.CctxStatus_OutboundMined {
				minedCCTX = cctx
				foundMined = true
				break
			}
		}
		if !foundMined {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgUpdateTssAddress,
				"no mined cross chain txs found in mined state",
			), nil, nil
		}

		// Thee tss migrator is set for all chains supporting tss migration
		for _, chain := range supportedChains {
			tssMigrator := observertypes.TssFundMigratorInfo{
				ChainId:            chain.ChainId,
				MigrationCctxIndex: minedCCTX.Index,
			}
			k.GetObserverKeeper().SetFundMigrator(ctx, tssMigrator)
		}

		oldTss, found := k.GetObserverKeeper().GetTSS(ctx)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgUpdateTssAddress,
				"no TSS found",
			), nil, nil
		}

		// Set the new TSS to state
		newTss, err := sample.TSSFromRand(r)
		newTss.FinalizedZetaHeight = oldTss.FinalizedZetaHeight + 10
		newTss.KeyGenZetaHeight = oldTss.KeyGenZetaHeight + 10
		if err != nil {
			return simtypes.NoOpMsg(
					types.ModuleName,
					TypeMsgUpdateTssAddress,
					err.Error()),
				nil, nil
		}
		k.GetObserverKeeper().SetTSSHistory(ctx, newTss)

		msg := types.MsgUpdateTssAddress{
			Creator:   policyAccount.Address.String(),
			TssPubkey: newTss.TssPubkey,
		}

		err = msg.ValidateBasic()
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgUpdateTssAddress,
				"unable to validate MsgUpdateTssAddress msg",
			), nil, err
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

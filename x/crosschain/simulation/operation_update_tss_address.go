package simulation

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// SimulateMsgUpdateTssAddress generates a MsgUpdateTssAddress with random values and delivers it
func SimulateMsgUpdateTssAddress(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		policyAccount, err := GetPolicyAccount(ctx, k.GetAuthorityKeeper(), accounts)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWhitelistERC20, err.Error()), nil, nil
		}

		authAccount := k.GetAuthKeeper().GetAccount(ctx, policyAccount.Address)
		spendable := k.GetBankKeeper().SpendableCoins(ctx, authAccount.GetAddress())

		supportedChains := k.GetChainsSupportingTSSMigration(ctx)
		if len(supportedChains) == 0 {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgUpdateTssAddress,
				"no chains found which support tss migration",
			), nil, nil
		}

		for _, chain := range supportedChains {
			index := ethcrypto.Keccak256Hash([]byte(fmt.Sprintf("%d", r.Int63()))).Hex()
			cctx := types.CrossChainTx{Index: index,
				CctxStatus: &types.Status{Status: types.CctxStatus_OutboundMined}}
			tssmigrator := observertypes.TssFundMigratorInfo{
				ChainId:            chain.ChainId,
				MigrationCctxIndex: index,
			}
			k.SetCrossChainTx(ctx, cctx)
			k.GetObserverKeeper().SetFundMigrator(ctx, tssmigrator)
		}

		oldTss, found := k.GetObserverKeeper().GetTSS(ctx)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgUpdateTssAddress,
				"no TSS found",
			), nil, nil
		}
		newTss, err := sample.TSSFromRand(r)
		newTss.FinalizedZetaHeight = oldTss.FinalizedZetaHeight + 10
		newTss.KeyGenZetaHeight = oldTss.KeyGenZetaHeight + 10
		if err != nil {
			return simtypes.NoOpMsg(
					types.ModuleName,
					types.TypeMsgUpdateTssAddress,
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
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to validate MsgRemoveOutboundTracker msg"), nil, err
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

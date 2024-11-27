package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func SimulateMsgWhitelistERC20(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		policyAccount, err := GetPolicyAccount(ctx, k.GetAuthorityKeeper(), accounts)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWhitelistERC20, err.Error()), nil, nil
		}

		authAccount := k.GetAuthKeeper().GetAccount(ctx, policyAccount.Address)
		spendable := k.GetBankKeeper().SpendableCoins(ctx, authAccount.GetAddress())

		supportedChains := k.GetObserverKeeper().GetSupportedChains(ctx)
		if len(supportedChains) == 0 {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgWhitelistERC20,
				"no supported chains found",
			), nil, nil
		}

		filteredChains := chains.FilterChains(supportedChains, chains.FilterByVM(chains.Vm_evm))

		//pick a random chain
		randomChain := supportedChains[r.Intn(len(filteredChains))]
		var tokenAddress string
		switch {
		case randomChain.IsEVMChain():
			tokenAddress = sample.EthAddressFromRand(r).String()
			//updatedTokenAddress = sample.EthAddressFromRand(r).String()
		default:
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgWhitelistERC20, "unsupported chain"), nil, nil
		}

		msg := types.MsgWhitelistERC20{
			Creator:      policyAccount.Address.String(),
			ChainId:      randomChain.ChainId,
			Erc20Address: tokenAddress,
			GasLimit:     100000,
			Decimals:     18,
			Name:         "Test",
			Symbol:       "TST",
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

package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// SimulateUpdateERC20CustodyPauseStatus generates a MsgUpdateERC20CustodyPauseStatus with random values and delivers it
func SimulateUpdateERC20CustodyPauseStatus(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		policyAccount, err := GetPolicyAccount(ctx, k.GetAuthorityKeeper(), accounts)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeUpdateERC20CustodyPauseStatus, err.Error()), nil, nil
		}

		authAccount := k.GetAuthKeeper().GetAccount(ctx, policyAccount.Address)
		spendable := k.GetBankKeeper().SpendableCoins(ctx, authAccount.GetAddress())

		supportedChains := k.GetObserverKeeper().GetSupportedChains(ctx)
		if len(supportedChains) == 0 {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeUpdateERC20CustodyPauseStatus,
				"no supported chains found",
			), nil, nil
		}

		filteredChains := chains.FilterChains(supportedChains, chains.FilterExternalChains)

		//pick a random chain
		randomChain := filteredChains[r.Intn(len(filteredChains))]

		_, found := k.GetObserverKeeper().GetChainNonces(ctx, randomChain.ChainId)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeUpdateERC20CustodyPauseStatus,
				"no chain nonces found",
			), nil, nil
		}

		_, found = k.GetObserverKeeper().GetTSS(ctx)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeUpdateERC20CustodyPauseStatus,
				"no TSS found",
			), nil, nil
		}

		_, found = k.GetObserverKeeper().GetChainParamsByChainID(ctx, randomChain.ChainId)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeUpdateERC20CustodyPauseStatus,
				"no chain params found",
			), nil, nil
		}
		medianGasPrice, priorityFee, found := k.GetMedianGasValues(ctx, randomChain.ChainId)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeUpdateERC20CustodyPauseStatus,
				"no median gas values found",
			), nil, nil
		}
		medianGasPrice = medianGasPrice.MulUint64(types.ERC20CustodyPausingGasMultiplierEVM)
		priorityFee = priorityFee.MulUint64(types.ERC20CustodyPausingGasMultiplierEVM)

		if priorityFee.GT(medianGasPrice) {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeUpdateERC20CustodyPauseStatus,
				"priorityFee is greater than median gasPrice",
			), nil, nil
		}

		msg := types.MsgUpdateERC20CustodyPauseStatus{
			Creator: policyAccount.Address.String(),
			ChainId: randomChain.ChainId,
			Pause:   r.Intn(2) == 0,
		}

		err = msg.ValidateBasic()
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				msg.Type(),
				"unable to validate MsgUpdateERC20CustodyPauseStatus msg",
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

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

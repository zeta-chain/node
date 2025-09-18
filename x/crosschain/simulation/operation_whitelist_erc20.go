package simulation

import (
	"math/rand"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	zetasimulation "github.com/zeta-chain/node/testutil/simulation"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// SimulateMsgWhitelistERC20 generates a MsgWhitelistERC20 with random values and delivers it
func SimulateMsgWhitelistERC20(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		policyAccount, err := zetasimulation.GetPolicyAccount(ctx, k.GetAuthorityKeeper(), accounts)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgWhitelistERC20, err.Error()), nil, nil
		}

		authAccount := k.GetAuthKeeper().GetAccount(ctx, policyAccount.Address)
		spendable := k.GetBankKeeper().SpendableCoins(ctx, authAccount.GetAddress())

		supportedChains := k.GetObserverKeeper().GetSupportedChains(ctx)
		if len(supportedChains) == 0 {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgWhitelistERC20,
				"no supported chains found",
			), nil, nil
		}

		filteredChains := chains.FilterChains(supportedChains, chains.FilterByVM(chains.Vm_evm))
		if len(filteredChains) == 0 {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgWhitelistERC20,
				"no EVM-compatible chains found",
			), nil, nil
		}

		//pick a random chain
		// Keep the switch case to add solana support in future
		// TODO : https://github.com/zeta-chain/node/issues/3287
		randomChain := filteredChains[r.Intn(len(filteredChains))]
		var tokenAddress string
		switch {
		case randomChain.IsEVMChain():
			tokenAddress = sample.EthAddressFromRand(r).String()
		default:
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgWhitelistERC20, "unsupported chain"), nil, nil
		}

		_, found := k.GetObserverKeeper().GetTSS(ctx)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgWhitelistERC20,
				"no TSS found",
			), nil, nil
		}

		_, found = k.GetObserverKeeper().GetChainParamsByChainID(ctx, randomChain.ChainId)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgWhitelistERC20,
				"no chain params found",
			), nil, nil
		}

		medianGasPrice, priorityFee, isFound := k.GetMedianGasValues(ctx, randomChain.ChainId)
		if !isFound {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgWhitelistERC20,
				"median gas price not found",
			), nil, nil
		}

		medianGasPrice = medianGasPrice.MulUint64(types.ERC20CustodyWhitelistGasMultiplierEVM)
		priorityFee = priorityFee.MulUint64(types.ERC20CustodyWhitelistGasMultiplierEVM)

		if priorityFee.GT(medianGasPrice) {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgWhitelistERC20,
				"priorityFee is greater than median gasPrice",
			), nil, nil
		}

		foreignCoins := k.GetFungibleKeeper().GetAllForeignCoins(ctx)
		for _, fCoin := range foreignCoins {
			if fCoin.Asset == tokenAddress && fCoin.ForeignChainId == randomChain.ChainId {
				return simtypes.NoOpMsg(
					types.ModuleName,
					TypeMsgWhitelistERC20,
					"ERC20 already whitelisted",
				), nil, nil
			}
		}

		gasLimit := r.Int63n(1000000000) + 1
		nameLength := r.Intn(97) + 3
		msg := types.MsgWhitelistERC20{
			Creator:      policyAccount.Address.String(),
			ChainId:      randomChain.ChainId,
			Erc20Address: tokenAddress,
			GasLimit:     gasLimit,
			Decimals:     18,
			Name:         sample.StringRandom(r, nameLength),
			Symbol:       sample.StringRandom(r, 3),
			LiquidityCap: sdkmath.NewUint(sample.Uint64InRangeFromRand(r, 1, 1000000000000000000)),
		}

		err = msg.ValidateBasic()
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgWhitelistERC20,
				"unable to validate MsgWhitelistERC20",
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

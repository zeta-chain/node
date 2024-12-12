package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/zeta-chain/node/pkg/authz"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// SimulateMsgVoteGasPrice generates a MsgVoteGasPrice and delivers it
func SimulateMsgVoteGasPrice(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		// Get a random account and observer
		// If this returns an error, it is likely that the entire observer set has been removed
		simAccount, randomObserver, err := GetRandomAccountAndObserver(r, ctx, k, accounts)
		if err != nil {
			return simtypes.OperationMsg{}, nil, nil
		}
		authAccount := k.GetAuthKeeper().GetAccount(ctx, simAccount.Address)
		spendable := k.GetBankKeeper().SpendableCoins(ctx, authAccount.GetAddress())

		supportedChains := k.GetObserverKeeper().GetSupportedChains(ctx)
		if len(supportedChains) == 0 {
			return simtypes.NoOpMsg(
				types.ModuleName,
				authz.GasPriceVoter.String(),
				"no supported chains found",
			), nil, nil
		}
		randomChainID := GetRandomChainID(r, supportedChains)

		// Vote for random gas price. Gas prices do not use a ballot system, so we can vote directly without having to schedule future operations.

		var price uint64
		for price == 0 {
			maxGasPrice := uint64(1000 * 1e9) // 1000 Gwei
			price = uint64(1e9) + r.Uint64()%maxGasPrice
		}
		// Select priority fee between 0 and price
		priorityFee := r.Uint64() % price
		msg := types.MsgVoteGasPrice{
			Creator:     randomObserver,
			ChainId:     randomChainID,
			Price:       price,
			PriorityFee: priorityFee,
			BlockNumber: uint64(ctx.BlockHeight()) + r.Uint64()%1000, // nosec G115 - overflow is not a issue here
			Supply:      sdk.NewInt(r.Int63n(1e18)).String(),
		}

		// System contracts are deployed on the first block, so we cannot vote on gas prices before that
		if ctx.BlockHeight() <= 1 {
			return simtypes.NewOperationMsg(&msg, true, "block height less than 1", nil), nil, nil
		}

		err = msg.ValidateBasic()
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to validate vote gas price  msg"), nil, err
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           moduletestutil.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             &msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   k.GetAuthKeeper(),
			Bankkeeper:      k.GetBankKeeper(),
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

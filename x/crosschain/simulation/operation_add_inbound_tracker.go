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
)

// SimulateMsgAddInboundTracker generates a MsgAddInboundTracker with random values
func SimulateMsgAddInboundTracker(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		// Get a random account and observer
		// If this returns an error, it is likely that the entire observer set has been removed
		simAccount, randomObserver, _, err := zetasimulation.GetRandomAccountAndObserver(
			r,
			ctx,
			k.GetObserverKeeper(),
			accounts,
		)
		if err != nil {
			return simtypes.OperationMsg{}, nil, nil
		}
		authAccount := k.GetAuthKeeper().GetAccount(ctx, simAccount.Address)
		spendable := k.GetBankKeeper().SpendableCoins(ctx, authAccount.GetAddress())

		supportedChains := k.GetObserverKeeper().GetSupportedChains(ctx)
		if len(supportedChains) == 0 {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgAddInboundTracker,
				"no supported chains found",
			), nil, nil
		}
		randomChainID := zetasimulation.GetRandomChainID(r, supportedChains)
		txHash := sample.HashFromRand(r)
		coinType := sample.CoinTypeFromRand(r)

		// Add a new inbound Tracker
		msg := types.MsgAddInboundTracker{
			Creator:  randomObserver,
			ChainId:  randomChainID,
			TxHash:   txHash.String(),
			CoinType: coinType,
		}

		err = msg.ValidateBasic()
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgAddInboundTracker,
				"unable to validate MsgAddInboundTracker",
			), nil, err
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           moduletestutil.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             &msg,
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

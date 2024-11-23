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

// SimulateMsgAddOutboundTracker generates a MsgAddOutboundTracker with random values
func SimulateMsgAddOutboundTracker(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {

		chainID := int64(1337)
		supportedChains := k.GetObserverKeeper().GetSupportedChains(ctx)
		for _, chain := range supportedChains {
			if chains.IsEthereumChain(chain.ChainId, []chains.Chain{}) {
				chainID = chain.ChainId
			}

		}
		// Get a random account and observer
		// If this returns an error, it is likely that the entire observer set has been removed
		simAccount, randomObserver, err := GetRandomAccountAndObserver(r, ctx, k, accounts)
		if err != nil {
			return simtypes.OperationMsg{}, nil, nil
		}
		authAccount := k.GetAuthKeeper().GetAccount(ctx, simAccount.Address)
		spendable := k.GetBankKeeper().SpendableCoins(ctx, authAccount.GetAddress())

		txHash := sample.HashFromRand(r)
		tss, found := k.GetObserverKeeper().GetTSS(ctx)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgAddOutboundTracker,
				"no TSS found",
			), nil, nil
		}

		pendingNonces, found := k.GetObserverKeeper().GetPendingNonces(ctx, tss.TssPubkey, chainID)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgAddOutboundTracker,
				"no pending nonces found",
			), nil, nil
		}

		// pick a random nonce from the pending nonces between 0 and nonceLow
		//fmt.Printf("pendingNonces.NonceLow: %d | pendingNonces.NonceHigh: %d \n",
		//	pendingNonces.NonceLow, pendingNonces.NonceHigh)
		if pendingNonces.NonceLow == pendingNonces.NonceHigh {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgAddOutboundTracker,
				"no pending nonces found",
			), nil, nil
		}

		nonce := pendingNonces.NonceLow

		tracker, found := k.GetOutboundTracker(ctx, chainID, uint64(nonce))
		if found && tracker.IsMaxed() {
			return simtypes.NoOpMsg(
				types.ModuleName,
				types.TypeMsgAddOutboundTracker,
				"tracker is maxed",
			), nil, nil
		}
		// Add a new inbound Tracker
		msg := types.MsgAddOutboundTracker{
			Creator:   randomObserver,
			ChainId:   chainID,
			Nonce:     uint64(nonce),
			TxHash:    txHash.String(),
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
		}

		// System contracts are deployed on the first block, so we cannot vote on gas prices before that
		if ctx.BlockHeight() <= 2 {
			return simtypes.NewOperationMsg(&msg, true, "block height less than 1", nil), nil, nil
		}

		err = msg.ValidateBasic()
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to validate MsgAddOutboundTracker msg"), nil, err
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

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
	zetasimulation "github.com/zeta-chain/node/testutil/simulation"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// SimulateMsgAddOutboundTracker generates a MsgAddOutboundTracker with random values
func SimulateMsgAddOutboundTracker(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		chainID := chains.GoerliLocalnet.ChainId
		supportedChains := k.GetObserverKeeper().GetSupportedChains(ctx)
		if len(supportedChains) == 0 {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgAddOutboundTracker,
				"no supported chains found",
			), nil, nil
		}

		for _, chain := range supportedChains {
			if chains.IsEthereumChain(chain.ChainId, []chains.Chain{}) {
				chainID = chain.ChainId
			}
		}
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

		txHash := sample.HashFromRand(r)
		tss, found := k.GetObserverKeeper().GetTSS(ctx)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgAddOutboundTracker,
				"no TSS found",
			), nil, nil
		}

		pendingNonces, found := k.GetObserverKeeper().GetPendingNonces(ctx, tss.TssPubkey, chainID)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgAddOutboundTracker,
				"pending nonces object not found",
			), nil, nil
		}

		// pick a random nonce from the pending nonces between 0 and nonceLow
		// If nonce low is the same as nonce high, it means that there are no pending nonces to add trackers for
		if pendingNonces.NonceLow == pendingNonces.NonceHigh {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgAddOutboundTracker,
				"no pending nonces found",
			), nil, nil
		}
		// Pick a random pending nonce
		nonce := 0
		switch {
		case pendingNonces.NonceHigh <= 1:
			nonce = int(pendingNonces.NonceLow)
		case pendingNonces.NonceLow == 0:
			nonce = r.Intn(int(pendingNonces.NonceHigh))
		default:
			nonce = r.Intn(int(pendingNonces.NonceHigh)-int(pendingNonces.NonceLow)) + int(pendingNonces.NonceLow)
		}

		// Verify if the tracker is maxed
		tracker, found := k.GetOutboundTracker(
			ctx,
			chainID,
			uint64(nonce),
		) // #nosec G115 - overflow is not an issue here
		if found && tracker.MaxReached() {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgAddOutboundTracker,
				"tracker is maxed",
			), nil, nil
		}

		// Verify the nonceToCCTX exists
		nonceToCCTX, found := k.GetObserverKeeper().GetNonceToCctx(ctx, tss.TssPubkey, chainID, int64(nonce))
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgAddOutboundTracker,
				"no nonce to cctx found",
			), nil, nil
		}

		// Verify the cctx exists
		_, found = k.GetCrossChainTx(ctx, nonceToCCTX.CctxIndex)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgAddOutboundTracker,
				"no cctx found for nonce",
			), nil, nil
		}
		// Add a new inbound Tracker
		msg := types.MsgAddOutboundTracker{
			Creator: randomObserver,
			ChainId: chainID,
			Nonce:   uint64(nonce), // #nosec G115 - overflow is not an issue here
			TxHash:  txHash.String(),
		}

		err = msg.ValidateBasic()
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgAddOutboundTracker,
				"unable to validate MsgAddOutboundTracker msg",
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

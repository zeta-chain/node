package keeper_test

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"testing"

	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func createNCctxWithStatus(
	keeper *keeper.Keeper,
	ctx sdk.Context,
	n int,
	status types.CctxStatus,
	tssPubkey string,
) []types.CrossChainTx {
	items := make([]types.CrossChainTx, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Index = fmt.Sprintf("%d-%d", i, status)
		items[i].CctxStatus = &types.Status{
			Status:              status,
			StatusMessage:       "",
			LastUpdateTimestamp: 0,
		}
		items[i].ZetaFees = math.OneUint()
		items[i].InboundParams = &types.InboundParams{ObservedHash: fmt.Sprintf("%d", i), Amount: math.OneUint()}
		items[i].OutboundParams = []*types.OutboundParams{{Amount: math.ZeroUint(), CallOptions: &types.CallOptions{}}}
		items[i].RevertOptions = types.NewEmptyRevertOptions()

		keeper.SaveCCTXUpdate(ctx, items[i], tssPubkey)
	}
	return items
}

// Keeper Tests
func createNCctx(keeper *keeper.Keeper, ctx sdk.Context, n int, tssPubkey string) []types.CrossChainTx {
	items := make([]types.CrossChainTx, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].InboundParams = &types.InboundParams{
			Sender:                 fmt.Sprintf("%d", i),
			SenderChainId:          int64(i),
			TxOrigin:               fmt.Sprintf("%d", i),
			Asset:                  fmt.Sprintf("%d", i),
			CoinType:               coin.CoinType_Zeta,
			ObservedHash:           fmt.Sprintf("%d", i),
			ObservedExternalHeight: uint64(i),
			FinalizedZetaHeight:    uint64(i),
		}
		items[i].OutboundParams = []*types.OutboundParams{{
			Receiver:        fmt.Sprintf("%d", i),
			ReceiverChainId: int64(i),
			Hash:            fmt.Sprintf("%d", i),
			TssNonce:        uint64(i),
			CallOptions: &types.CallOptions{
				GasLimit: uint64(i),
			},
			GasPrice:               fmt.Sprintf("%d", i),
			BallotIndex:            fmt.Sprintf("%d", i),
			ObservedExternalHeight: uint64(i),
			CoinType:               coin.CoinType_Zeta,
		}}
		items[i].CctxStatus = &types.Status{
			Status:              types.CctxStatus_PendingInbound,
			StatusMessage:       "any",
			LastUpdateTimestamp: 0,
		}
		items[i].InboundParams.Amount = math.OneUint()

		items[i].ZetaFees = math.OneUint()
		items[i].Index = sample.GetCctxIndexFromString(fmt.Sprintf("%d", i))
		items[i].RevertOptions = types.NewEmptyRevertOptions()

		keeper.SaveCCTXUpdate(ctx, items[i], tssPubkey)
	}
	return items
}

func TestCCTXs(t *testing.T) {
	cctxsTest := []struct {
		TestName        string
		PendingInbound  int
		PendingOutbound int
		OutboundMined   int
		Confirmed       int
		PendingRevert   int
		Reverted        int
		Aborted         int
	}{
		{
			TestName:        "test pending",
			PendingInbound:  10,
			PendingOutbound: 10,
			Confirmed:       10,
			PendingRevert:   10,
			Aborted:         10,
			OutboundMined:   10,
			Reverted:        10,
		},
		{
			TestName:        "test pending random",
			PendingInbound:  rand.Intn(300-10) + 10,
			PendingOutbound: rand.Intn(300-10) + 10,
			Confirmed:       rand.Intn(300-10) + 10,
			PendingRevert:   rand.Intn(300-10) + 10,
			Aborted:         rand.Intn(300-10) + 10,
			OutboundMined:   rand.Intn(300-10) + 10,
			Reverted:        rand.Intn(300-10) + 10,
		},
	}
	for _, tt := range cctxsTest {
		t.Run(tt.TestName, func(t *testing.T) {
			keeper, ctx, _, zk := keepertest.CrosschainKeeper(t)
			keeper.SetZetaAccounting(ctx, types.ZetaAccounting{AbortedZetaAmount: math.ZeroUint()})
			var sends []types.CrossChainTx
			tss := sample.Tss()
			zk.ObserverKeeper.SetTSS(ctx, tss)
			sends = append(
				sends,
				createNCctxWithStatus(
					keeper,
					ctx,
					tt.PendingInbound,
					types.CctxStatus_PendingInbound,
					tss.TssPubkey,
				)...)
			sends = append(
				sends,
				createNCctxWithStatus(
					keeper,
					ctx,
					tt.PendingOutbound,
					types.CctxStatus_PendingOutbound,
					tss.TssPubkey,
				)...)
			sends = append(
				sends,
				createNCctxWithStatus(keeper, ctx, tt.PendingRevert, types.CctxStatus_PendingRevert, tss.TssPubkey)...)
			sends = append(
				sends,
				createNCctxWithStatus(keeper, ctx, tt.Aborted, types.CctxStatus_Aborted, tss.TssPubkey)...)
			sends = append(
				sends,
				createNCctxWithStatus(keeper, ctx, tt.OutboundMined, types.CctxStatus_OutboundMined, tss.TssPubkey)...)
			sends = append(
				sends,
				createNCctxWithStatus(keeper, ctx, tt.Reverted, types.CctxStatus_Reverted, tss.TssPubkey)...)

			require.Equal(t, len(sends), len(keeper.GetAllCrossChainTx(ctx)))
			for _, s := range sends {
				send, found := keeper.GetCrossChainTx(ctx, s.Index)
				require.True(t, found)
				require.Equal(t, s, send)
			}
		})
	}
}

func compareCctx(l types.CrossChainTx, r types.CrossChainTx) int {
	return strings.Compare(l.Index, r.Index)
}

func TestCCTXGetAll(t *testing.T) {
	keeper, ctx, _, zk := keepertest.CrosschainKeeper(t)
	tss := sample.Tss()
	zk.ObserverKeeper.SetTSS(ctx, tss)
	items := createNCctx(keeper, ctx, 10, tss.TssPubkey)
	cctx := keeper.GetAllCrossChainTx(ctx)

	slices.SortFunc(items, compareCctx)
	slices.SortFunc(cctx, compareCctx)

	require.Equal(t, items, cctx)
}

// Querier Tests

func TestCCTXQuerySingle(t *testing.T) {
	keeper, ctx, _, zk := keepertest.CrosschainKeeper(t)
	tss := sample.Tss()
	zk.ObserverKeeper.SetTSS(ctx, tss)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNCctx(keeper, ctx, 2, tss.TssPubkey)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetCctxRequest
		response *types.QueryGetCctxResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetCctxRequest{Index: msgs[0].Index},
			response: &types.QueryGetCctxResponse{CrossChainTx: &msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetCctxRequest{Index: msgs[1].Index},
			response: &types.QueryGetCctxResponse{CrossChainTx: &msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetCctxRequest{Index: "missing"},
			err:     status.Error(codes.InvalidArgument, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Cctx(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.Equal(t, tc.response, response)
			}
		})
	}
}

func TestCCTXQueryPaginated(t *testing.T) {
	keeper, ctx, _, zk := keepertest.CrosschainKeeper(t)
	tss := sample.Tss()
	zk.ObserverKeeper.SetTSS(ctx, sample.Tss())
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNCctx(keeper, ctx, 5, tss.TssPubkey)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllCctxRequest {
		return &types.QueryAllCctxRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
				Reverse:    true,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.CctxAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			for j := i; j < len(msgs) && j < i+step; j++ {
				require.Equal(t, &msgs[j], resp.CrossChainTx[j-i])
			}
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.CctxAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			for j := i; j < len(msgs) && j < i+step; j++ {
				require.Equal(t, &msgs[j], resp.CrossChainTx[j-i])
			}
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.CctxAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.CctxAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

func TestKeeper_RemoveCrossChainTx(t *testing.T) {
	keeper, ctx, _, zk := keepertest.CrosschainKeeper(t)
	tss := sample.Tss()
	zk.ObserverKeeper.SetTSS(ctx, tss)
	txs := createNCctx(keeper, ctx, 5, tss.TssPubkey)

	keeper.RemoveCrossChainTx(ctx, txs[0].Index)
	txs = keeper.GetAllCrossChainTx(ctx)
	require.Equal(t, 4, len(txs))
}

func TestCrossChainTx_AddOutbound(t *testing.T) {
	t.Run("successfully get outbound tx", func(t *testing.T) {
		_, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		hash := sample.Hash().String()

		err := cctx.UpdateCurrentOutbound(ctx, types.MsgVoteOutbound{
			ValueReceived:                     cctx.GetCurrentOutboundParam().Amount,
			ObservedOutboundHash:              hash,
			ObservedOutboundBlockHeight:       10,
			ObservedOutboundGasUsed:           100,
			ObservedOutboundEffectiveGasPrice: sdkmath.NewInt(100),
			ObservedOutboundEffectiveGasLimit: 20,
			ConfirmationMode:                  types.ConfirmationMode_SAFE,
		}, observertypes.BallotStatus_BallotFinalized_SuccessObservation)
		require.NoError(t, err)
		require.Equal(t, cctx.GetCurrentOutboundParam().Hash, hash)
		require.Equal(t, cctx.GetCurrentOutboundParam().GasUsed, uint64(100))
		require.Equal(t, cctx.GetCurrentOutboundParam().EffectiveGasPrice, sdkmath.NewInt(100))
		require.Equal(t, cctx.GetCurrentOutboundParam().EffectiveGasLimit, uint64(20))
		require.Equal(t, cctx.GetCurrentOutboundParam().ObservedExternalHeight, uint64(10))
		require.Equal(t, cctx.GetCurrentOutboundParam().ConfirmationMode, types.ConfirmationMode_SAFE)
	})

	t.Run("successfully get outbound tx for failed ballot without amount check", func(t *testing.T) {
		_, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		hash := sample.Hash().String()

		err := cctx.UpdateCurrentOutbound(ctx, types.MsgVoteOutbound{
			ObservedOutboundHash:              hash,
			ObservedOutboundBlockHeight:       10,
			ObservedOutboundGasUsed:           100,
			ObservedOutboundEffectiveGasPrice: sdkmath.NewInt(100),
			ObservedOutboundEffectiveGasLimit: 20,
			ConfirmationMode:                  types.ConfirmationMode_SAFE,
		}, observertypes.BallotStatus_BallotFinalized_FailureObservation)
		require.NoError(t, err)
		require.Equal(t, cctx.GetCurrentOutboundParam().Hash, hash)
		require.Equal(t, cctx.GetCurrentOutboundParam().GasUsed, uint64(100))
		require.Equal(t, cctx.GetCurrentOutboundParam().EffectiveGasPrice, sdkmath.NewInt(100))
		require.Equal(t, cctx.GetCurrentOutboundParam().EffectiveGasLimit, uint64(20))
		require.Equal(t, cctx.GetCurrentOutboundParam().ObservedExternalHeight, uint64(10))
		require.Equal(t, cctx.GetCurrentOutboundParam().ConfirmationMode, types.ConfirmationMode_SAFE)
	})

	t.Run("failed to get outbound tx if amount does not match value received", func(t *testing.T) {
		_, ctx, _, _ := keepertest.CrosschainKeeper(t)

		cctx := sample.CrossChainTx(t, "test")
		hash := sample.Hash().String()

		err := cctx.UpdateCurrentOutbound(ctx, types.MsgVoteOutbound{
			ValueReceived:                     sdkmath.NewUint(100),
			ObservedOutboundHash:              hash,
			ObservedOutboundBlockHeight:       10,
			ObservedOutboundGasUsed:           100,
			ObservedOutboundEffectiveGasPrice: sdkmath.NewInt(100),
			ObservedOutboundEffectiveGasLimit: 20,
		}, observertypes.BallotStatus_BallotFinalized_SuccessObservation)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
	})
}

func Test_NewCCTX(t *testing.T) {
	t.Run("should return a cctx with correct values", func(t *testing.T) {
		_, ctx, _, _ := keepertest.CrosschainKeeper(t)
		senderChain := chains.Goerli
		sender := sample.EthAddress()
		receiverChain := chains.Goerli
		receiver := sample.EthAddress()
		creator := sample.AccAddress()
		amount := sdkmath.NewUint(42)
		message := "test"
		inboundBlockHeight := uint64(420)
		inboundHash := sample.Hash()
		gasLimit := uint64(100)
		asset := "test-asset"
		eventIndex := uint64(1)
		cointType := coin.CoinType_ERC20
		tss := sample.Tss()
		msg := types.MsgVoteInbound{
			Creator:            creator,
			Sender:             sender.String(),
			SenderChainId:      senderChain.ChainId,
			Receiver:           receiver.String(),
			ReceiverChain:      receiverChain.ChainId,
			Amount:             amount,
			Message:            message,
			InboundHash:        inboundHash.String(),
			InboundBlockHeight: inboundBlockHeight,
			CallOptions: &types.CallOptions{
				GasLimit: gasLimit,
			},
			CoinType:                cointType,
			TxOrigin:                sender.String(),
			Asset:                   asset,
			EventIndex:              eventIndex,
			ProtocolContractVersion: types.ProtocolContractVersion_V2,
			ConfirmationMode:        types.ConfirmationMode_FAST,
			Status:                  types.InboundStatus_INSUFFICIENT_DEPOSITOR_FEE,
			ErrorMessage:            "deposited amount is less than depositor fee",
		}
		cctx, err := types.NewCCTX(ctx, msg, tss.TssPubkey)
		require.NoError(t, err)
		require.Equal(t, receiver.String(), cctx.GetCurrentOutboundParam().Receiver)
		require.Equal(t, receiverChain.ChainId, cctx.GetCurrentOutboundParam().ReceiverChainId)
		require.Equal(t, sender.String(), cctx.GetInboundParams().Sender)
		require.Equal(t, senderChain.ChainId, cctx.GetInboundParams().SenderChainId)
		require.Equal(t, amount, cctx.GetInboundParams().Amount)
		require.Equal(t, message, cctx.RelayedMessage)
		require.Equal(t, inboundHash.String(), cctx.GetInboundParams().ObservedHash)
		require.Equal(t, inboundBlockHeight, cctx.GetInboundParams().ObservedExternalHeight)
		require.Equal(t, gasLimit, cctx.GetCurrentOutboundParam().CallOptions.GasLimit)
		require.Equal(t, asset, cctx.GetInboundParams().Asset)
		require.Equal(t, cointType, cctx.InboundParams.CoinType)
		require.Equal(t, uint64(0), cctx.GetCurrentOutboundParam().TssNonce)
		require.Equal(t, sdkmath.ZeroUint(), cctx.GetCurrentOutboundParam().Amount)
		require.Equal(t, types.CctxStatus_PendingInbound, cctx.CctxStatus.Status)
		require.Equal(t, false, cctx.CctxStatus.IsAbortRefunded)
		require.Equal(t, types.ProtocolContractVersion_V2, cctx.ProtocolContractVersion)
		require.Equal(t, types.ConfirmationMode_FAST, cctx.GetInboundParams().ConfirmationMode)
		require.Equal(t, types.ConfirmationMode_SAFE, cctx.GetCurrentOutboundParam().ConfirmationMode)
		require.Equal(t, types.InboundStatus_INSUFFICIENT_DEPOSITOR_FEE, cctx.GetInboundParams().Status)
		require.Equal(t, "deposited amount is less than depositor fee", cctx.GetInboundParams().ErrorMessage)
	})

	t.Run("should return an error if the cctx is invalid", func(t *testing.T) {
		_, ctx, _, _ := keepertest.CrosschainKeeper(t)
		senderChain := chains.Goerli
		sender := sample.EthAddress()
		receiverChain := chains.Goerli
		receiver := sample.EthAddress()
		creator := sample.AccAddress()
		amount := sdkmath.NewUint(42)
		message := "test"
		inboundBlockHeight := uint64(420)
		inboundHash := sample.Hash()
		gasLimit := uint64(100)
		asset := "test-asset"
		eventIndex := uint64(1)
		cointType := coin.CoinType_ERC20
		tss := sample.Tss()
		msg := types.MsgVoteInbound{
			Creator:            creator,
			Sender:             "",
			SenderChainId:      senderChain.ChainId,
			Receiver:           receiver.String(),
			ReceiverChain:      receiverChain.ChainId,
			Amount:             amount,
			Message:            message,
			InboundHash:        inboundHash.String(),
			InboundBlockHeight: inboundBlockHeight,
			CallOptions: &types.CallOptions{
				GasLimit: gasLimit,
			},
			CoinType:   cointType,
			TxOrigin:   sender.String(),
			Asset:      asset,
			EventIndex: eventIndex,
		}
		_, err := types.NewCCTX(ctx, msg, tss.TssPubkey)
		require.ErrorContains(t, err, "sender cannot be empty")
	})

	t.Run("zero value for protocol contract version gives V1", func(t *testing.T) {
		cctx := types.CrossChainTx{}
		require.Equal(t, types.ProtocolContractVersion_V1, cctx.ProtocolContractVersion)
	})
}

func TestKeeper_UpdateNonceToCCTX(t *testing.T) {
	t.Run("should set nonce to cctx if status is PendingOutbound", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		chainID := chains.Ethereum.ChainId
		nonce := uint64(10)

		cctx := types.CrossChainTx{Index: "test",
			OutboundParams: []*types.OutboundParams{{ReceiverChainId: chainID, TssNonce: nonce}},
			CctxStatus:     &types.Status{Status: types.CctxStatus_PendingOutbound},
		}
		tssPubkey := "test-tss-pubkey"

		// Act
		k.SetNonceToCCTX(ctx, cctx, tssPubkey)

		// Assert
		nonceToCctx, found := k.GetObserverKeeper().GetNonceToCctx(ctx, tssPubkey, chainID, int64(nonce))
		require.True(t, found)
		require.Equal(t, cctx.Index, nonceToCctx.CctxIndex)
		require.Equal(t, tssPubkey, nonceToCctx.Tss)
		require.Equal(t, chainID, nonceToCctx.ChainId)
	})

	t.Run("should set nonce to cctx if status is PendingRevert", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		chainID := chains.Ethereum.ChainId
		nonce := uint64(10)

		cctx := types.CrossChainTx{Index: "test",
			OutboundParams: []*types.OutboundParams{{ReceiverChainId: chainID, TssNonce: nonce}},
			CctxStatus:     &types.Status{Status: types.CctxStatus_PendingRevert},
		}
		tssPubkey := "test-tss-pubkey"

		// Act
		k.SetNonceToCCTX(ctx, cctx, tssPubkey)

		// Assert
		nonceToCctx, found := k.GetObserverKeeper().GetNonceToCctx(ctx, tssPubkey, chainID, int64(nonce))
		require.True(t, found)
		require.Equal(t, cctx.Index, nonceToCctx.CctxIndex)
		require.Equal(t, tssPubkey, nonceToCctx.Tss)
		require.Equal(t, chainID, nonceToCctx.ChainId)
	})

	t.Run("should not set nonce to cctx if status is not PendingOutbound or PendingRevert", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		chainID := chains.Ethereum.ChainId
		nonce := uint64(10)

		cctx := types.CrossChainTx{Index: "test",
			OutboundParams: []*types.OutboundParams{{ReceiverChainId: chainID, TssNonce: nonce}},
			CctxStatus:     &types.Status{Status: types.CctxStatus_Aborted},
		}
		tssPubkey := "test-tss-pubkey"

		// Act
		k.SetNonceToCCTX(ctx, cctx, tssPubkey)

		// Assert
		_, found := k.GetObserverKeeper().GetNonceToCctx(ctx, tssPubkey, chainID, int64(nonce))
		require.False(t, found)
	})
}

func TestKeeper_UpdateInboundHashToCCTX(t *testing.T) {
	t.Run(
		"should update inbound hash to cctx mapping if new cctx index is found for the same inbound hash",
		func(t *testing.T) {
			// Arrange
			k, ctx, _, _ := keepertest.CrosschainKeeper(t)
			inboundHash := sample.Hash().String()
			index1 := sample.ZetaIndex(t)
			index2 := sample.ZetaIndex(t)

			inboundHashToCctx := types.InboundHashToCctx{
				InboundHash: inboundHash,
				CctxIndex:   []string{index1},
			}
			k.SetInboundHashToCctx(ctx, inboundHashToCctx)
			cctx := types.CrossChainTx{Index: index2, InboundParams: &types.InboundParams{ObservedHash: inboundHash}}

			// Act
			k.UpdateInboundHashToCCTX(ctx, cctx)

			// Assert
			inboundHashToCctx, found := k.GetInboundHashToCctx(ctx, inboundHash)
			require.True(t, found)
			require.Equal(t, inboundHash, inboundHashToCctx.InboundHash)
			require.Equal(t, 2, len(inboundHashToCctx.CctxIndex))
			require.Contains(t, inboundHashToCctx.CctxIndex, index1)
			require.Contains(t, inboundHashToCctx.CctxIndex, index2)
		},
	)

	t.Run("should do nothing if the cctx index is already in the mapping", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		inboundHash := sample.Hash().String()
		index := sample.ZetaIndex(t)

		inboundHashToCctx := types.InboundHashToCctx{
			InboundHash: inboundHash,
			CctxIndex:   []string{index},
		}
		k.SetInboundHashToCctx(ctx, inboundHashToCctx)
		cctx := types.CrossChainTx{Index: index, InboundParams: &types.InboundParams{ObservedHash: inboundHash}}

		// Act
		k.UpdateInboundHashToCCTX(ctx, cctx)

		// Assert
		inboundHashToCctx, found := k.GetInboundHashToCctx(ctx, inboundHash)
		require.True(t, found)
		require.Equal(t, inboundHash, inboundHashToCctx.InboundHash)
		require.Equal(t, 1, len(inboundHashToCctx.CctxIndex))
		require.Contains(t, inboundHashToCctx.CctxIndex, index)
	})

	t.Run("should add cctx index to mapping if InboundHashToCctx is not found", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		inboundHash := sample.Hash().String()
		index := sample.ZetaIndex(t)

		cctx := types.CrossChainTx{Index: index, InboundParams: &types.InboundParams{ObservedHash: inboundHash}}

		// Act
		k.UpdateInboundHashToCCTX(ctx, cctx)

		// Assert
		inboundHashToCctx, found := k.GetInboundHashToCctx(ctx, inboundHash)
		require.True(t, found)
		require.Equal(t, inboundHash, inboundHashToCctx.InboundHash)
		require.Equal(t, 1, len(inboundHashToCctx.CctxIndex))
		require.Contains(t, inboundHashToCctx.CctxIndex, index)
	})
}

func TestKeeper_UpdateZetaAccounting(t *testing.T) {
	t.Run("should update zeta accounting if cctx is aborted and coin type is zeta", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		amount := sdkmath.NewUint(100)
		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{CoinType: coin.CoinType_Zeta},
			CctxStatus: &types.Status{
				IsAbortRefunded: false,
				Status:          types.CctxStatus_Aborted},
			OutboundParams: []*types.OutboundParams{{Amount: amount}},
		}
		k.SetZetaAccounting(ctx, types.ZetaAccounting{AbortedZetaAmount: math.ZeroUint()})

		// Act
		k.UpdateZetaAccounting(ctx, cctx)

		// Assert
		zetaAccounting, found := k.GetZetaAccounting(ctx)
		require.True(t, found)
		require.Equal(t, amount, zetaAccounting.AbortedZetaAmount)
	})

	t.Run("should not update zeta accounting if cctx is not aborted", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		amount := sdkmath.NewUint(100)
		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{CoinType: coin.CoinType_Zeta},
			CctxStatus: &types.Status{
				IsAbortRefunded: false,
				Status:          types.CctxStatus_PendingOutbound},
			OutboundParams: []*types.OutboundParams{{Amount: amount}},
		}
		k.SetZetaAccounting(ctx, types.ZetaAccounting{AbortedZetaAmount: math.ZeroUint()})

		// Act
		k.UpdateZetaAccounting(ctx, cctx)

		// Assert
		zetaAccounting, found := k.GetZetaAccounting(ctx)
		require.True(t, found)
		require.Equal(t, math.ZeroUint(), zetaAccounting.AbortedZetaAmount)
	})

	t.Run("should update to amount if zeta accounting is not set", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		amount := sdkmath.NewUint(100)
		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{CoinType: coin.CoinType_Zeta},
			CctxStatus: &types.Status{
				IsAbortRefunded: false,
				Status:          types.CctxStatus_Aborted},
			OutboundParams: []*types.OutboundParams{{Amount: amount}},
		}

		// Act
		k.UpdateZetaAccounting(ctx, cctx)

		// Assert
		zetaAccounting, found := k.GetZetaAccounting(ctx)
		require.True(t, found)
		require.Equal(t, amount, zetaAccounting.AbortedZetaAmount)
	})

	t.Run("should not update zeta accounting if the cctx is already refunded", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		amount := sdkmath.NewUint(100)
		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{CoinType: coin.CoinType_Zeta},
			CctxStatus: &types.Status{
				IsAbortRefunded: true,
				Status:          types.CctxStatus_Aborted},
			OutboundParams: []*types.OutboundParams{{Amount: amount}},
		}
		k.SetZetaAccounting(ctx, types.ZetaAccounting{AbortedZetaAmount: math.ZeroUint()})

		// Act
		k.UpdateZetaAccounting(ctx, cctx)

		// Assert
		zetaAccounting, found := k.GetZetaAccounting(ctx)
		require.True(t, found)
		require.Equal(t, math.ZeroUint(), zetaAccounting.AbortedZetaAmount)
	})
}

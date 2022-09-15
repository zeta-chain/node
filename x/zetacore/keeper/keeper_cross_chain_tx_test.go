package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

// Keeper Tests
func createNSend(keeper *Keeper, ctx sdk.Context, n int) []types.CrossChainTx {
	items := make([]types.CrossChainTx, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].InBoundTxParams = &types.InBoundTxParams{
			Sender:                   fmt.Sprintf("%d", i),
			SenderChain:              fmt.Sprintf("%d", i),
			InBoundTxObservedHash:    fmt.Sprintf("%d", i),
			InBoundTxObservedHeight:  uint64(i),
			InBoundTxFinalizedHeight: uint64(i),
		}
		items[i].OutBoundTxParams = &types.OutBoundTxParams{
			Receiver:               fmt.Sprintf("%d", i),
			ReceiverChain:          fmt.Sprintf("%d", i),
			Broadcaster:            uint64(i),
			OutBoundTxHash:         fmt.Sprintf("%d", i),
			OutBoundTxTSSNonce:     uint64(i),
			OutBoundTxGasLimit:     uint64(i),
			OutBoundTxGasPrice:     fmt.Sprintf("%d", i),
			OutBoundTXReceiveIndex: fmt.Sprintf("%d", i),
		}
		items[i].CctxStatus = &types.Status{
			Status:              types.CctxStatus_PendingInbound,
			StatusMessage:       "any",
			LastUpdateTimestamp: 0,
		}
		items[i].ZetaBurnt = sdk.OneUint()
		items[i].ZetaMint = sdk.OneUint()
		items[i].Index = fmt.Sprintf("%d", i)
		keeper.SetCrossChainTx(ctx, items[i])
	}
	return items
}

func TestSendGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNSend(keeper, ctx, 1)
	for _, item := range items {
		rst, found := keeper.GetCrossChainTx(ctx, item.Index)
		assert.True(t, found)
		assert.Equal(t, item, rst)
	}
}
func TestSendRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNSend(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveCrossChainTx(ctx, item.Index)
		_, found := keeper.GetCrossChainTx(ctx, item.Index)
		assert.False(t, found)
	}
}

func TestSendGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNSend(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllCrossChainTx(ctx))
}

// Querier Tests

func TestSendQuerySingle(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNSend(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetSendRequest
		response *types.QueryGetSendResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetSendRequest{Index: msgs[0].Index},
			response: &types.QueryGetSendResponse{CrossChainTx: &msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetSendRequest{Index: msgs[1].Index},
			response: &types.QueryGetSendResponse{CrossChainTx: &msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetSendRequest{Index: "missing"},
			err:     status.Error(codes.InvalidArgument, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Send(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.Equal(t, tc.response, response)
			}
		})
	}
}

func TestSendQueryPaginated(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNSend(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllSendRequest {
		return &types.QueryAllSendRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.SendAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			for j := i; j < len(msgs) && j < i+step; j++ {
				assert.Equal(t, &msgs[j], resp.CrossChainTx[j-i])
			}
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.SendAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			for j := i; j < len(msgs) && j < i+step; j++ {
				assert.Equal(t, &msgs[j], resp.CrossChainTx[j-i])
			}
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.SendAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.SendAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

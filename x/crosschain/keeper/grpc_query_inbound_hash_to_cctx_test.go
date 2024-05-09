package keeper_test

import (
	"fmt"
	"strconv"
	"testing"

	"cosmossdk.io/math"
	"github.com/zeta-chain/zetacore/testutil/sample"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	crosschainkeeper "github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestInboundHashToCctxQuerySingle(t *testing.T) {
	keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNInboundHashToCctx(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetInboundHashToCctxRequest
		response *types.QueryGetInboundHashToCctxResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetInboundHashToCctxRequest{
				InboundHash: msgs[0].InboundHash,
			},
			response: &types.QueryGetInboundHashToCctxResponse{InboundHashToCctx: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetInboundHashToCctxRequest{
				InboundHash: msgs[1].InboundHash,
			},
			response: &types.QueryGetInboundHashToCctxResponse{InboundHashToCctx: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetInboundHashToCctxRequest{
				InboundHash: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.InboundHashToCctx(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestInTxHashToCctxQueryPaginated(t *testing.T) {
	keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNInboundHashToCctx(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllInboundHashToCctxRequest {
		return &types.QueryAllInboundHashToCctxRequest{
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
			resp, err := keeper.InboundHashToCctxAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.InboundHashToCctx), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.InboundHashToCctx),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.InboundHashToCctxAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.InboundHashToCctx), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.InboundHashToCctx),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.InboundHashToCctxAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.InboundHashToCctx),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.InboundHashToCctxAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

func createInTxHashToCctxWithCctxs(keeper *crosschainkeeper.Keeper, ctx sdk.Context) ([]types.CrossChainTx,
	types.InboundHashToCctx) {
	cctxs := make([]types.CrossChainTx, 5)
	for i := range cctxs {
		cctxs[i].Creator = "any"
		cctxs[i].Index = fmt.Sprintf("0x123%d", i)
		cctxs[i].ZetaFees = math.OneUint()
		cctxs[i].InboundParams = &types.InboundParams{ObservedHash: fmt.Sprintf("%d", i), Amount: math.OneUint()}
		cctxs[i].CctxStatus = &types.Status{Status: types.CctxStatus_PendingInbound}
		keeper.SetCctxAndNonceToCctxAndInboundHashToCctx(ctx, cctxs[i])
	}

	var inboundHashToCctx types.InboundHashToCctx
	inboundHashToCctx.InboundHash = fmt.Sprintf("0xabc")
	for i := range cctxs {
		inboundHashToCctx.CctxIndex = append(inboundHashToCctx.CctxIndex, cctxs[i].Index)
	}
	keeper.SetInboundHashToCctx(ctx, inboundHashToCctx)

	return cctxs, inboundHashToCctx
}

func TestKeeper_InTxHashToCctxDataQuery(t *testing.T) {
	keeper, ctx, _, zk := keepertest.CrosschainKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	zk.ObserverKeeper.SetTSS(ctx, sample.Tss())
	t.Run("can query all cctxs data with in tx hash", func(t *testing.T) {
		cctxs, inboundHashToCctx := createInTxHashToCctxWithCctxs(keeper, ctx)
		req := &types.QueryInboundHashToCctxDataRequest{
			InboundHash: inboundHashToCctx.InboundHash,
		}
		res, err := keeper.InboundHashToCctxData(wctx, req)
		require.NoError(t, err)
		require.Equal(t, len(cctxs), len(res.CrossChainTxs))
		for i := range cctxs {
			require.Equal(t, nullify.Fill(cctxs[i]), nullify.Fill(res.CrossChainTxs[i]))
		}
	})
	t.Run("in tx hash not found", func(t *testing.T) {
		req := &types.QueryInboundHashToCctxDataRequest{
			InboundHash: "notfound",
		}
		_, err := keeper.InboundHashToCctxData(wctx, req)
		require.ErrorIs(t, err, status.Error(codes.NotFound, "not found"))
	})
	t.Run("cctx not indexed return internal error", func(t *testing.T) {
		keeper.SetInboundHashToCctx(ctx, types.InboundHashToCctx{
			InboundHash: "nocctx",
			CctxIndex:   []string{"notfound"},
		})

		req := &types.QueryInboundHashToCctxDataRequest{
			InboundHash: "nocctx",
		}
		_, err := keeper.InboundHashToCctxData(wctx, req)
		require.ErrorIs(t, err, status.Error(codes.Internal, "cctx indexed notfound doesn't exist"))
	})
}

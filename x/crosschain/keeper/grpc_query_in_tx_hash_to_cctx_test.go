package keeper_test

import (
	"cosmossdk.io/math"
	"fmt"
	"strconv"
	"testing"

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

func TestInTxHashToCctxQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.CrosschainKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNInTxHashToCctx(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetInTxHashToCctxRequest
		response *types.QueryGetInTxHashToCctxResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetInTxHashToCctxRequest{
				InTxHash: msgs[0].InTxHash,
			},
			response: &types.QueryGetInTxHashToCctxResponse{InTxHashToCctx: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetInTxHashToCctxRequest{
				InTxHash: msgs[1].InTxHash,
			},
			response: &types.QueryGetInTxHashToCctxResponse{InTxHashToCctx: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetInTxHashToCctxRequest{
				InTxHash: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.InTxHashToCctx(wctx, tc.request)
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
	keeper, ctx := keepertest.CrosschainKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNInTxHashToCctx(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllInTxHashToCctxRequest {
		return &types.QueryAllInTxHashToCctxRequest{
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
			resp, err := keeper.InTxHashToCctxAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.InTxHashToCctx), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.InTxHashToCctx),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.InTxHashToCctxAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.InTxHashToCctx), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.InTxHashToCctx),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.InTxHashToCctxAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.InTxHashToCctx),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.InTxHashToCctxAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

func createInTxHashToCctxWithCctxs(keeper *crosschainkeeper.Keeper, ctx sdk.Context) ([]types.CrossChainTx,
	types.InTxHashToCctx) {
	cctxs := make([]types.CrossChainTx, 5)
	for i := range cctxs {
		cctxs[i].Creator = "any"
		cctxs[i].Index = fmt.Sprintf("0x123%d", i)
		cctxs[i].ZetaFees = math.OneUint()
		cctxs[i].InboundTxParams = &types.InboundTxParams{InboundTxObservedHash: fmt.Sprintf("%d", i), Amount: math.OneUint()}
		keeper.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, cctxs[i])
	}

	var inTxHashToCctx types.InTxHashToCctx
	inTxHashToCctx.InTxHash = fmt.Sprintf("0xabc")
	for i := range cctxs {
		inTxHashToCctx.CctxIndex = append(inTxHashToCctx.CctxIndex, cctxs[i].Index)
	}
	keeper.SetInTxHashToCctx(ctx, inTxHashToCctx)

	return cctxs, inTxHashToCctx
}

func TestKeeper_InTxHashToCctxDataQuery(t *testing.T) {
	keeper, ctx := keepertest.CrosschainKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)

	t.Run("can query all cctxs data with in tx hash", func(t *testing.T) {
		cctxs, inTxHashToCctx := createInTxHashToCctxWithCctxs(keeper, ctx)
		req := &types.QueryInTxHashToCctxDataRequest{
			InTxHash: inTxHashToCctx.InTxHash,
		}
		res, err := keeper.InTxHashToCctxData(wctx, req)
		require.NoError(t, err)
		require.Equal(t, len(cctxs), len(res.CrossChainTxs))
		for i := range cctxs {
			require.Equal(t, nullify.Fill(cctxs[i]), nullify.Fill(res.CrossChainTxs[i]))
		}
	})
	t.Run("in tx hash not found", func(t *testing.T) {
		req := &types.QueryInTxHashToCctxDataRequest{
			InTxHash: "notfound",
		}
		_, err := keeper.InTxHashToCctxData(wctx, req)
		require.ErrorIs(t, err, status.Error(codes.NotFound, "not found"))
	})
	t.Run("cctx not indexed return internal error", func(t *testing.T) {
		keeper.SetInTxHashToCctx(ctx, types.InTxHashToCctx{
			InTxHash:  "nocctx",
			CctxIndex: []string{"notfound"},
		})

		req := &types.QueryInTxHashToCctxDataRequest{
			InTxHash: "nocctx",
		}
		_, err := keeper.InTxHashToCctxData(wctx, req)
		require.ErrorIs(t, err, status.Error(codes.Internal, "cctx indexed notfound doesn't exist"))
	})
}

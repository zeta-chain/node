package keeper_test

import (
	"math/rand"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/crypto"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestTSSQuerySingle(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)
	tss := sample.Tss()
	wctx := sdk.WrapSDKContext(ctx)

	for _, tc := range []struct {
		desc           string
		request        *types.QueryGetTSSRequest
		response       *types.QueryGetTSSResponse
		skipSettingTss bool
		err            error
	}{
		{
			desc:           "Skip setting tss",
			request:        &types.QueryGetTSSRequest{},
			skipSettingTss: true,
			err:            status.Error(codes.InvalidArgument, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
		{
			desc:     "Should return tss",
			request:  &types.QueryGetTSSRequest{},
			response: &types.QueryGetTSSResponse{TSS: tss},
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			if !tc.skipSettingTss {
				k.SetTSS(ctx, tss)
			}
			response, err := k.TSS(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.Equal(t, tc.response, response)
			}
		})
	}
}

func TestTSSQueryHistory(t *testing.T) {
	keeper, ctx, _, _ := keepertest.ObserverKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	for _, tc := range []struct {
		desc          string
		tssCount      int
		foundPrevious bool
		err           error
	}{
		{
			desc:          "1 Tss addresses",
			tssCount:      1,
			foundPrevious: false,
			err:           nil,
		},
		{
			desc:          "10 Tss addresses",
			tssCount:      10,
			foundPrevious: true,
			err:           nil,
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			tssList := sample.TssList(tc.tssCount)
			for _, tss := range tssList {
				keeper.SetTSS(ctx, tss)
				keeper.SetTSSHistory(ctx, tss)
			}
			request := &types.QueryTssHistoryRequest{}
			response, err := keeper.TssHistory(wctx, request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.Equal(t, len(tssList), len(response.TssList))
				prevTss, found := keeper.GetPreviousTSS(ctx)
				require.Equal(t, tc.foundPrevious, found)
				if found {
					require.Equal(t, tssList[len(tssList)-2], prevTss)
				}
			}
		})
	}
}

func TestKeeper_GetTssAddress(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.GetTssAddress(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should error if not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.GetTssAddress(wctx, &types.QueryGetTssAddressRequest{
			BitcoinChainId: 1,
		})
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should error if invalid chain id", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		tss := sample.Tss()
		k.SetTSS(ctx, tss)

		res, err := k.GetTssAddress(wctx, &types.QueryGetTssAddressRequest{
			BitcoinChainId: 987,
		})
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should return if valid chain id", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		tss := sample.Tss()
		k.SetTSS(ctx, tss)

		res, err := k.GetTssAddress(wctx, &types.QueryGetTssAddressRequest{
			BitcoinChainId: chains.BtcRegtestChain.ChainId,
		})
		require.NoError(t, err)
		expectedBitcoinParams, err := chains.BitcoinNetParamsFromChainID(chains.BtcRegtestChain.ChainId)
		require.NoError(t, err)
		expectedBtcAddress, err := crypto.GetTssAddrBTC(tss.TssPubkey, expectedBitcoinParams)
		require.NoError(t, err)
		expectedEthAddress, err := crypto.GetTssAddrEVM(tss.TssPubkey)
		require.NoError(t, err)
		require.NotNil(t, &types.QueryGetTssAddressByFinalizedHeightResponse{
			Eth: expectedEthAddress.String(),
			Btc: expectedBtcAddress,
		}, res)
	})
}

func TestKeeper_GetTssAddressByFinalizedHeight(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.GetTssAddressByFinalizedHeight(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should error if not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.GetTssAddressByFinalizedHeight(wctx, &types.QueryGetTssAddressByFinalizedHeightRequest{
			BitcoinChainId: 1,
		})
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should error if invalid chain id", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		tssList := sample.TssList(100)
		r := rand.Intn((len(tssList)-1)-0) + 0
		for _, tss := range tssList {
			k.SetTSSHistory(ctx, tss)
		}

		res, err := k.GetTssAddressByFinalizedHeight(wctx, &types.QueryGetTssAddressByFinalizedHeightRequest{
			BitcoinChainId:      987,
			FinalizedZetaHeight: tssList[r].FinalizedZetaHeight,
		})
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should return if valid chain id", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		tssList := sample.TssList(100)
		r := rand.Intn((len(tssList)-1)-0) + 0
		for _, tss := range tssList {
			k.SetTSSHistory(ctx, tss)
		}

		res, err := k.GetTssAddressByFinalizedHeight(wctx, &types.QueryGetTssAddressByFinalizedHeightRequest{
			BitcoinChainId:      chains.BtcRegtestChain.ChainId,
			FinalizedZetaHeight: tssList[r].FinalizedZetaHeight,
		})
		require.NoError(t, err)
		expectedBitcoinParams, err := chains.BitcoinNetParamsFromChainID(chains.BtcRegtestChain.ChainId)
		require.NoError(t, err)
		expectedBtcAddress, err := crypto.GetTssAddrBTC(tssList[r].TssPubkey, expectedBitcoinParams)
		require.NoError(t, err)
		expectedEthAddress, err := crypto.GetTssAddrEVM(tssList[r].TssPubkey)
		require.NoError(t, err)
		require.NotNil(t, &types.QueryGetTssAddressByFinalizedHeightResponse{
			Eth: expectedEthAddress.String(),
			Btc: expectedBtcAddress,
		}, res)
	})
}

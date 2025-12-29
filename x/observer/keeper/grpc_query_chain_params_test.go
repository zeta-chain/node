package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestKeeper_GetChainParamsForChain(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.GetChainParamsForChain(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should error if chain params not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.GetChainParamsForChain(wctx, &types.QueryGetChainParamsForChainRequest{
			ChainId: 987,
		})
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should return if chain params found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		list := types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:               chains.ZetaChainPrivnet.ChainId,
					IsSupported:           false,
					BallotThreshold:       sdkmath.LegacyZeroDec(),
					MinObserverDelegation: sdkmath.LegacyZeroDec(),
					GasPriceMultiplier:    types.DefaultGasPriceMultiplier,
				},
			},
		}
		k.SetChainParamsList(ctx, list)

		res, err := k.GetChainParamsForChain(wctx, &types.QueryGetChainParamsForChainRequest{
			ChainId: chains.ZetaChainPrivnet.ChainId,
		})
		require.NoError(t, err)
		require.Equal(t, &types.QueryGetChainParamsForChainResponse{
			ChainParams: list.ChainParams[0],
		}, res)
	})
}

func TestKeeper_GetChainParams(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.GetChainParams(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should error if chain params not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.GetChainParams(wctx, &types.QueryGetChainParamsRequest{})
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should return if chain params found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		list := types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:               chains.ZetaChainPrivnet.ChainId,
					IsSupported:           false,
					BallotThreshold:       sdkmath.LegacyZeroDec(),
					MinObserverDelegation: sdkmath.LegacyZeroDec(),
					GasPriceMultiplier:    types.DefaultGasPriceMultiplier,
				},
			},
		}
		k.SetChainParamsList(ctx, list)

		res, err := k.GetChainParams(wctx, &types.QueryGetChainParamsRequest{})
		require.NoError(t, err)
		require.Equal(t, &types.QueryGetChainParamsResponse{
			ChainParams: &list,
		}, res)
	})
}

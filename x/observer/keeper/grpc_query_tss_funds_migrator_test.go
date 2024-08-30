package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestKeeper_TssFundsMigratorInfo(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.TssFundsMigratorInfo(wctx, nil)
		require.ErrorContains(t, err, "invalid request")
		require.Nil(t, res)
	})

	t.Run("should error if chain id is invalid", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.TssFundsMigratorInfo(wctx, &types.QueryTssFundsMigratorInfoRequest{
			ChainId: 0,
		})
		require.ErrorContains(t, err, "invalid chain id")
		require.Nil(t, res)
	})

	t.Run("should error if not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.TssFundsMigratorInfo(wctx, &types.QueryTssFundsMigratorInfoRequest{
			ChainId: chains.Ethereum.ChainId,
		})
		require.ErrorContains(t, err, "tss fund migrator not found")
		require.Nil(t, res)
	})

	t.Run("should return if found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)
		chainId := chains.Ethereum.ChainId

		fm := types.TssFundMigratorInfo{
			ChainId:            chainId,
			MigrationCctxIndex: sample.ZetaIndex(t),
		}
		k.SetFundMigrator(ctx, fm)

		res, err := k.TssFundsMigratorInfo(wctx, &types.QueryTssFundsMigratorInfoRequest{
			ChainId: chainId,
		})
		require.NoError(t, err)
		require.Equal(t, fm, res.TssFundsMigrator)
	})
}

func TestKeeper_TssFundsMigratorInfoAll(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.TssFundsMigratorInfoAll(wctx, nil)
		require.ErrorContains(t, err, "invalid request")
		require.Nil(t, res)
	})

	t.Run("should return empty list if not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.TssFundsMigratorInfoAll(wctx, &types.QueryTssFundsMigratorInfoAllRequest{})
		require.NoError(t, err)
		require.Equal(t, 0, len(res.TssFundsMigrators))
	})

	t.Run("should return list of infos if found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		migrators := make([]types.TssFundMigratorInfo, 3)
		for i := 0; i < 3; i++ {
			fm := types.TssFundMigratorInfo{
				ChainId:            int64(i),
				MigrationCctxIndex: sample.ZetaIndex(t),
			}
			k.SetFundMigrator(ctx, fm)
			migrators[i] = fm
		}

		res, err := k.TssFundsMigratorInfoAll(wctx, &types.QueryTssFundsMigratorInfoAllRequest{})
		require.NoError(t, err)
		require.Equal(t, 3, len(res.TssFundsMigrators))
		require.ElementsMatch(t, migrators, res.TssFundsMigrators)
	})
}

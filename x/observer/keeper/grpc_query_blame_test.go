package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestKeeper_BlameByIdentifier(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.BlameByIdentifier(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should error if blame not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.BlameByIdentifier(wctx, &types.QueryBlameByIdentifierRequest{
			BlameIdentifier: "test",
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return blame info if found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)
		var chainId int64 = 97
		var nonce uint64 = 101
		digest := sample.ZetaIndex(t)

		index := types.GetBlameIndex(chainId, nonce, digest, 123)
		blame := types.Blame{
			Index:         index,
			FailureReason: "failed to join party",
			Nodes:         nil,
		}
		k.SetBlame(ctx, blame)

		res, err := k.BlameByIdentifier(wctx, &types.QueryBlameByIdentifierRequest{
			BlameIdentifier: index,
		})
		require.NoError(t, err)
		require.Equal(t, &types.QueryBlameByIdentifierResponse{
			BlameInfo: &blame,
		}, res)
	})
}

func TestKeeper_GetAllBlameRecords(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.GetAllBlameRecords(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should return all if found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)
		var chainId int64 = 97
		var nonce uint64 = 101
		digest := sample.ZetaIndex(t)

		index := types.GetBlameIndex(chainId, nonce, digest, 123)
		blame := types.Blame{
			Index:         index,
			FailureReason: "failed to join party",
			Nodes:         nil,
		}
		k.SetBlame(ctx, blame)

		res, err := k.GetAllBlameRecords(wctx, &types.QueryAllBlameRecordsRequest{})
		require.NoError(t, err)
		require.Equal(t, []types.Blame{blame}, res.BlameInfo)
	})
}

func TestKeeper_BlamesByChainAndNonce(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.BlamesByChainAndNonce(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should error if blame not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.BlamesByChainAndNonce(wctx, &types.QueryBlameByChainAndNonceRequest{
			ChainId: 1,
			Nonce:   1,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return blame info if found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)
		var chainId int64 = 97
		var nonce uint64 = 101
		digest := sample.ZetaIndex(t)

		index := types.GetBlameIndex(chainId, nonce, digest, 123)
		blame := types.Blame{
			Index:         index,
			FailureReason: "failed to join party",
			Nodes:         nil,
		}
		k.SetBlame(ctx, blame)

		res, err := k.BlamesByChainAndNonce(wctx, &types.QueryBlameByChainAndNonceRequest{
			ChainId: chainId,
			Nonce:   int64(nonce),
		})
		require.NoError(t, err)
		require.Equal(t, &types.QueryBlameByChainAndNonceResponse{
			BlameInfo: []*types.Blame{&blame},
		}, res)
	})
}

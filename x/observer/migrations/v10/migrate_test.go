package v10_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/keeper"
	v10 "github.com/zeta-chain/node/x/observer/migrations/v10"
	"github.com/zeta-chain/node/x/observer/types"
)

var chainParams = types.ChainParamsList{
	ChainParams: []*types.ChainParams{
		makeChainParamsEmptyConfirmation(1, 14),
		makeChainParamsEmptyConfirmation(56, 20),
		makeChainParamsEmptyConfirmation(8332, 3),
		makeChainParamsEmptyConfirmation(7000, 0),
		makeChainParamsEmptyConfirmation(137, 200),
		makeChainParamsEmptyConfirmation(8453, 90),
		makeChainParamsEmptyConfirmation(900, 32),
	},
}

func TestMigrateStore(t *testing.T) {
	t.Run("can migrate confirmation count", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// set chain params
		setChainParamsList(ctx, *k, chainParams)

		// ensure the chain params are set correctly
		oldChainParams, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.EqualValues(t, chainParams, oldChainParams)

		// migrate the store
		err := v10.MigrateStore(ctx, *k)
		require.NoError(t, err)

		// ensure we still have 7 chain params after migration
		newChainParams, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, newChainParams.ChainParams, len(oldChainParams.ChainParams))

		// compare the old and new chain params
		for i, newParam := range newChainParams.ChainParams {
			oldParam := oldChainParams.ChainParams[i]

			// ensure the confirmation fields are set correctly
			require.Equal(t, newParam.Confirmation.SafeInboundCount, oldParam.ConfirmationCount)
			require.Equal(t, newParam.Confirmation.FastInboundCount, oldParam.ConfirmationCount)
			require.Equal(t, newParam.Confirmation.SafeOutboundCount, oldParam.ConfirmationCount)
			require.Equal(t, newParam.Confirmation.FastOutboundCount, oldParam.ConfirmationCount)

			// ensure nothing else has changed except the confirmation
			oldParam.Confirmation.SafeInboundCount = oldParam.ConfirmationCount
			oldParam.Confirmation.FastInboundCount = oldParam.ConfirmationCount
			oldParam.Confirmation.SafeOutboundCount = oldParam.ConfirmationCount
			oldParam.Confirmation.FastOutboundCount = oldParam.ConfirmationCount
			require.Equal(t, newParam, oldParam)
		}
	})

	t.Run("migrate nothing if chain params not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// ensure no chain params are set
		allChainParams, found := k.GetChainParamsList(ctx)
		require.False(t, found)
		require.Empty(t, allChainParams.ChainParams)

		// migrate the store
		err := v10.MigrateStore(ctx, *k)
		require.ErrorIs(t, err, types.ErrChainParamsNotFound)

		// ensure nothing has changed
		allChainParams, found = k.GetChainParamsList(ctx)
		require.False(t, found)
		require.Empty(t, allChainParams.ChainParams)
	})
}

// makeChainParamsEmptyConfirmation creates a sample chain params with empty confirmation
func makeChainParamsEmptyConfirmation(chainID int64, confirmationCount uint64) *types.ChainParams {
	chainParams := sample.ChainParams(chainID)
	chainParams.ConfirmationCount = confirmationCount
	chainParams.Confirmation = types.Confirmation{}
	return chainParams
}

// setChainParamsList set chain params list in the store
func setChainParamsList(ctx sdk.Context, observerKeeper keeper.Keeper, chainParams types.ChainParamsList) {
	store := ctx.KVStore(observerKeeper.StoreKey())
	b := observerKeeper.Codec().MustMarshal(&chainParams)
	key := types.KeyPrefix(types.AllChainParamsKey)
	store.Set(key, b)
}

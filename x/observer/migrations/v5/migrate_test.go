package v5_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	v5 "github.com/zeta-chain/zetacore/x/observer/migrations/v5"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMigrateObserverMapper(t *testing.T) {
	t.Run("TestMigrateStore", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		legacyObserverMapperStore := prefix.NewStore(ctx.KVStore(k.StoreKey()), types.KeyPrefix(types.ObserverMapperKey))
		legacyObserverMapperList := sample.LegacyObserverMapperList(t, 12, "sample")
		for _, legacyObserverMapper := range legacyObserverMapperList {
			legacyObserverMapperStore.Set(types.KeyPrefix(legacyObserverMapper.Index), k.Codec().MustMarshal(legacyObserverMapper))
		}
		err := v5.MigrateObserverMapper(ctx, k.StoreKey(), k.Codec())
		assert.NoError(t, err)
		observerSet, found := k.GetObserverSet(ctx)
		assert.True(t, found)

		assert.Equal(t, legacyObserverMapperList[0].ObserverList, observerSet.ObserverList)
		iterator := sdk.KVStorePrefixIterator(legacyObserverMapperStore, []byte{})
		defer iterator.Close()

		var observerMappers []*types.ObserverMapper
		for ; iterator.Valid(); iterator.Next() {
			var val types.ObserverMapper
			if !iterator.Valid() {
				k.Codec().MustUnmarshal(iterator.Value(), &val)
				observerMappers = append(observerMappers, &val)
			}
		}
		assert.Equal(t, 0, len(observerMappers))
	})

	t.Run("TestMigrateStoreNoObserverMapper", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		err := v5.MigrateObserverMapper(ctx, k.StoreKey(), k.Codec())
		assert.NoError(t, err)
		_, found := k.GetObserverSet(ctx)
		assert.False(t, found)
	})
}

func TestMigrateObserverParams(t *testing.T) {
	k, ctx := keepertest.ObserverKeeper(t)

	// set chain params
	previousChainParamsList := types.ChainParamsList{
		ChainParams: []*types.ChainParams{
			sample.ChainParams(1),
			sample.ChainParams(2),
			sample.ChainParams(3),
			sample.ChainParams(4),
		},
	}
	k.SetChainParamsList(ctx, previousChainParamsList)

	// set observer params
	dec42, err := sdk.NewDecFromStr("0.42")
	require.NoError(t, err)
	dec43, err := sdk.NewDecFromStr("0.43")
	require.NoError(t, err)
	dec1000, err := sdk.NewDecFromStr("1000.0")
	require.NoError(t, err)
	dec1001, err := sdk.NewDecFromStr("1001.0")
	require.NoError(t, err)
	params := types.Params{
		ObserverParams: []*types.ObserverParams{
			{
				Chain:                 &common.Chain{ChainId: 2},
				BallotThreshold:       dec42,
				MinObserverDelegation: dec1000,
				IsSupported:           true,
			},
			{
				Chain:                 &common.Chain{ChainId: 3},
				BallotThreshold:       dec43,
				MinObserverDelegation: dec1001,
				IsSupported:           true,
			},
		},
	}
	k.SetParams(ctx, params)

	// perform migration
	err = v5.MigrateObserverParams(ctx, *k)
	require.NoError(t, err)

	// check chain params
	newChainParamsList, found := k.GetChainParamsList(ctx)
	require.True(t, found)

	// unchanged values
	require.EqualValues(t, previousChainParamsList.ChainParams[0], newChainParamsList.ChainParams[0])
	require.EqualValues(t, previousChainParamsList.ChainParams[3], newChainParamsList.ChainParams[3])

	// changed values
	require.EqualValues(t, dec42, newChainParamsList.ChainParams[1].BallotThreshold)
	require.EqualValues(t, dec1000, newChainParamsList.ChainParams[1].MinObserverDelegation)
	require.EqualValues(t, dec43, newChainParamsList.ChainParams[2].BallotThreshold)
	require.EqualValues(t, dec1001, newChainParamsList.ChainParams[2].MinObserverDelegation)
	require.True(t, newChainParamsList.ChainParams[1].IsSupported)
	require.True(t, newChainParamsList.ChainParams[2].IsSupported)

	// check remaining values are unchanged
	previousChainParamsList.ChainParams[1].BallotThreshold = dec42
	previousChainParamsList.ChainParams[2].BallotThreshold = dec43
	previousChainParamsList.ChainParams[1].MinObserverDelegation = dec1000
	previousChainParamsList.ChainParams[2].MinObserverDelegation = dec1001
	previousChainParamsList.ChainParams[1].IsSupported = true
	previousChainParamsList.ChainParams[2].IsSupported = true
	require.EqualValues(t, previousChainParamsList.ChainParams[1], newChainParamsList.ChainParams[1])
	require.EqualValues(t, previousChainParamsList.ChainParams[2], newChainParamsList.ChainParams[2])
}

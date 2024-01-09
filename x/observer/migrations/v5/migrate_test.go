package v5_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	v5 "github.com/zeta-chain/zetacore/x/observer/migrations/v5"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMigrateStore(t *testing.T) {
	t.Run("TestMigrateStore", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		legacyObserverMapperStore := prefix.NewStore(ctx.KVStore(k.StoreKey()), types.KeyPrefix(types.ObserverMapperKey))
		legacyObserverMapperList := sample.LegacyObserverMapperList(t, 12, "sample")
		for _, legacyObserverMapper := range legacyObserverMapperList {
			legacyObserverMapperStore.Set(types.KeyPrefix(legacyObserverMapper.Index), k.Codec().MustMarshal(legacyObserverMapper))
		}
		err := v5.MigrateStore(ctx, k.StoreKey(), k.Codec())
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
}

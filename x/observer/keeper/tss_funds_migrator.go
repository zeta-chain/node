package keeper

import (
	"fmt"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/observer/types"
)

func (k Keeper) SetFundMigrator(ctx sdk.Context, fm types.TssFundMigratorInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TssFundMigratorKey))
	b := k.cdc.MustMarshal(&fm)
	store.Set(fmt.Appendf(nil, "%d", fm.ChainId), b)
}

func (k Keeper) GetFundMigrator(ctx sdk.Context, chainID int64) (val types.TssFundMigratorInfo, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TssFundMigratorKey))
	b := store.Get(fmt.Appendf(nil, "%d", chainID))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetAllTssFundMigrators(ctx sdk.Context) (fms []types.TssFundMigratorInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TssFundMigratorKey))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.TssFundMigratorInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		fms = append(fms, val)
	}
	return
}

func (k Keeper) RemoveAllExistingMigrators(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TssFundMigratorKey))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

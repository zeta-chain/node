package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

// SetForeignCoins set a specific foreignCoins in the store from its index
func (k Keeper) SetForeignCoins(ctx sdk.Context, foreignCoins types.ForeignCoins) {
	p := types.KeyPrefix(fmt.Sprintf("%s", types.ForeignCoinsKeyPrefix))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
	b := k.cdc.MustMarshal(&foreignCoins)
	store.Set(types.ForeignCoinsKey(
		foreignCoins.Zrc20ContractAddress,
	), b)
}

// GetForeignCoins returns a foreignCoins from its index
func (k Keeper) GetForeignCoins(
	ctx sdk.Context,
	zrc20Addr string,
) (val types.ForeignCoins, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(fmt.Sprintf("%s", types.ForeignCoinsKeyPrefix)))

	b := store.Get(types.ForeignCoinsKey(
		zrc20Addr,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveForeignCoins removes a foreignCoins from the store
func (k Keeper) RemoveForeignCoins(
	ctx sdk.Context,
	zrc20Addr string,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ForeignCoinsKeyPrefix))
	store.Delete(types.ForeignCoinsKey(
		zrc20Addr,
	))
}

// GetAllForeignCoinsForChain returns all foreignCoins on a given chain
func (k Keeper) GetAllForeignCoinsForChain(ctx sdk.Context, foreignChainID int64) (list []types.ForeignCoins) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(fmt.Sprintf("%s", types.ForeignCoinsKeyPrefix)))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ForeignCoins
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		if val.ForeignChainId == foreignChainID {
			list = append(list, val)
		}
	}
	return
}

// GetAllForeignCoins returns all foreignCoins
func (k Keeper) GetAllForeignCoins(ctx sdk.Context) (list []types.ForeignCoins) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(fmt.Sprintf("%s", types.ForeignCoinsKeyPrefix)))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ForeignCoins
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}
	return
}

func (k Keeper) GetGasCoinForForeignCoin(ctx sdk.Context, chainID int64) (types.ForeignCoins, bool) {
	foreignCoinList := k.GetAllForeignCoinsForChain(ctx, chainID)
	for _, coin := range foreignCoinList {
		if coin.CoinType == common.CoinType_Gas {
			return coin, true
		}
	}
	return types.ForeignCoins{}, false
}

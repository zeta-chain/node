package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

// SetTxList set txList in the store
func (k Keeper) SetTxList(ctx sdk.Context, txList types.TxList) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxListKey))
	b := k.cdc.MustMarshalBinaryBare(&txList)
	store.Set([]byte{0}, b)
}

// GetTxList returns txList
func (k Keeper) GetTxList(ctx sdk.Context) (val types.TxList, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxListKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshalBinaryBare(b, &val)
	return val, true
}

// RemoveTxList removes txList from the store
func (k Keeper) RemoveTxList(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxListKey))
	store.Delete([]byte{0})
}

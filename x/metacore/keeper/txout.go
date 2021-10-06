package keeper

import (
	"encoding/binary"
	"github.com/Meta-Protocol/metacore/x/metacore/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

// GetTxoutCount get the total number of TypeName.LowerCamel
func (k Keeper) GetTxoutCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxoutCountKey))
	byteKey := types.KeyPrefix(types.TxoutCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	count, err := strconv.ParseUint(string(bz), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to uint64
		panic("cannot decode count")
	}

	return count
}

// SetTxoutCount set the total number of txout
func (k Keeper) SetTxoutCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxoutCountKey))
	byteKey := types.KeyPrefix(types.TxoutCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendTxout appends a txout in the store with a new id and update the count
func (k Keeper) AppendTxout(
	ctx sdk.Context,
	txout types.Txout,
) uint64 {
	// Create the txout
	count := k.GetTxoutCount(ctx)

	// Set the ID of the appended value
	txout.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxoutKey))
	appendedValue := k.cdc.MustMarshalBinaryBare(&txout)
	store.Set(GetTxoutIDBytes(txout.Id), appendedValue)

	// Update txout count
	k.SetTxoutCount(ctx, count+1)

	return count
}

// SetTxout set a specific txout in the store
func (k Keeper) SetTxout(ctx sdk.Context, txout types.Txout) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxoutKey))
	b := k.cdc.MustMarshalBinaryBare(&txout)
	store.Set(GetTxoutIDBytes(txout.Id), b)
}

// GetTxout returns a txout from its id
func (k Keeper) GetTxout(ctx sdk.Context, id uint64) types.Txout {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxoutKey))
	var txout types.Txout
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetTxoutIDBytes(id)), &txout)
	return txout
}

// HasTxout checks if the txout exists in the store
func (k Keeper) HasTxout(ctx sdk.Context, id uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxoutKey))
	return store.Has(GetTxoutIDBytes(id))
}

// GetTxoutOwner returns the creator of the
func (k Keeper) GetTxoutOwner(ctx sdk.Context, id uint64) string {
	return k.GetTxout(ctx, id).Creator
}

// RemoveTxout removes a txout from the store
func (k Keeper) RemoveTxout(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxoutKey))
	store.Delete(GetTxoutIDBytes(id))
}

// GetAllTxout returns all txout
func (k Keeper) GetAllTxout(ctx sdk.Context) (list []types.Txout) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxoutKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Txout
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetTxoutIDBytes returns the byte representation of the ID
func GetTxoutIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetTxoutIDFromBytes returns ID in uint64 format from a byte array
func GetTxoutIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}

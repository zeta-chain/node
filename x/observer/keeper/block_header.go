package keeper

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// SetBlockHeader set a specific block header in the store from its index
func (k Keeper) SetBlockHeader(ctx sdk.Context, header common.BlockHeader) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlockHeaderKey))
	b := k.cdc.MustMarshal(&header)
	store.Set(header.Hash, b)
}

// GetBlockHeader returns a block header from its hash
func (k Keeper) GetBlockHeader(ctx sdk.Context, hash []byte) (val common.BlockHeader, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlockHeaderKey))

	b := store.Get(hash)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveBlockHeader removes a block header from the store
func (k Keeper) RemoveBlockHeader(ctx sdk.Context, hash []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlockHeaderKey))
	store.Delete(hash)
}

func (k Keeper) SetBlockHeaderState(ctx sdk.Context, blockHeaderState types.BlockHeaderState) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlockHeaderStateKey))
	b := k.cdc.MustMarshal(&blockHeaderState)
	key := strconv.FormatInt(blockHeaderState.ChainId, 10)
	store.Set(types.KeyPrefix(key), b)
}

func (k Keeper) GetBlockHeaderState(ctx sdk.Context, chainID int64) (val types.BlockHeaderState, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlockHeaderStateKey))

	b := store.Get(types.KeyPrefix(strconv.FormatInt(chainID, 10)))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

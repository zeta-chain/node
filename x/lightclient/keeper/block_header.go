package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

// SetBlockHeader set a specific block header in the store from its index
func (k Keeper) SetBlockHeader(ctx sdk.Context, header proofs.BlockHeader) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlockHeaderKey))
	b := k.cdc.MustMarshal(&header)
	store.Set(header.Hash, b)
}

// GetBlockHeader returns a block header from its hash
func (k Keeper) GetBlockHeader(ctx sdk.Context, hash []byte) (val proofs.BlockHeader, found bool) {
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

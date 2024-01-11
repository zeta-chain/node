package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func (k Keeper) SetFinalizedInbound(ctx sdk.Context, finalizedInboundIndex string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FinalizedInboundsKey))
	store.Set(types.KeyPrefix(finalizedInboundIndex), []byte{1})
}

func (k Keeper) AddFinalizedInbound(ctx sdk.Context, intxHash string, chainID int64, eventIndex uint64) {
	finalizedInboundIndex := types.FinalizedInboundKey(intxHash, chainID, eventIndex)
	k.SetFinalizedInbound(ctx, finalizedInboundIndex)
}
func (k Keeper) IsFinalizedInbound(ctx sdk.Context, intxHash string, chainID int64, eventIndex uint64) bool {
	finalizedInboundIndex := types.FinalizedInboundKey(intxHash, chainID, eventIndex)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FinalizedInboundsKey))
	return store.Has(types.KeyPrefix(finalizedInboundIndex))
}

func (k Keeper) GetAllFinalizedInbound(ctx sdk.Context) (list []string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FinalizedInboundsKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		iterator.Value()
		list = append(list, string(iterator.Key()))
	}
	return
}

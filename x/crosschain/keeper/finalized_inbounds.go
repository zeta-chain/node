package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func (k Keeper) AddFinalizedInbound(ctx sdk.Context, intxHash string, chainID int64, eventIndex uint64) {
	finalizedInboundIndex := fmt.Sprintf("%s-%d-%d", intxHash, chainID, eventIndex)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FinalizedInboundsKey))
	store.Set(types.KeyPrefix(finalizedInboundIndex), []byte{1})
}

func (k Keeper) IsFinalizedInbound(ctx sdk.Context, intxHash string, chainID int64, eventIndex uint64) bool {
	finalizedInboundIndex := fmt.Sprintf("%s-%d-%d", intxHash, chainID, eventIndex)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FinalizedInboundsKey))
	return store.Has(types.KeyPrefix(finalizedInboundIndex))
}

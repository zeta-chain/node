package keeper

import (
	"fmt"
	"sort"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// AppendTss appends a tss to the TSSHistoryKey store and update the current TSS to the latest one
func (k Keeper) AppendTss(ctx sdk.Context, tss types.TSS) {
	k.SetTSS(ctx, tss)
	k.SetTSSHistory(ctx, tss)
}

// SetTSS sets tss information to the store
func (k Keeper) SetTSS(ctx sdk.Context, tss types.TSS) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TSSKey))
	b := k.cdc.MustMarshal(&tss)
	store.Set([]byte{0}, b)
}

// SetTSSHistory Sets a new TSS into the TSS history store
func (k Keeper) SetTSSHistory(ctx sdk.Context, tss types.TSS) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TSSHistoryKey))
	b := k.cdc.MustMarshal(&tss)
	store.Set(types.KeyPrefix(fmt.Sprintf("%d", tss.FinalizedZetaHeight)), b)
}

// GetHistoricalTssByFinalizedHeight Returns the TSS address the specified finalized zeta height
// Finalized zeta height is the zeta block height at which the voting for the generation of a new TSS is finalized
func (k Keeper) GetHistoricalTssByFinalizedHeight(ctx sdk.Context, finalizedZetaHeight int64) (types.TSS, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TSSHistoryKey))
	b := store.Get(types.KeyPrefix(fmt.Sprintf("%d", finalizedZetaHeight)))
	if b == nil {
		return types.TSS{}, false
	}
	var tss types.TSS
	err := k.cdc.Unmarshal(b, &tss)
	if err != nil {
		return types.TSS{}, false
	}
	return tss, true

}

// GetTSS returns the current tss information
func (k Keeper) GetTSS(ctx sdk.Context) (val types.TSS, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TSSKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveTSS removes tss information from the store
func (k Keeper) RemoveTSS(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TSSKey))
	store.Delete([]byte{0})
}

// GetAllTSS returns all tss historical information from the store
func (k Keeper) GetAllTSS(ctx sdk.Context) (list []types.TSS) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TSSHistoryKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.TSS
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}
	return
}

func (k Keeper) GetAllTSSPaginated(ctx sdk.Context, pagination *query.PageRequest) (list []types.TSS, pageRes *query.PageResponse, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TSSHistoryKey))
	pageRes, err = query.Paginate(store, pagination, func(key []byte, value []byte) error {
		var tss types.TSS
		if err := k.cdc.Unmarshal(value, &tss); err != nil {
			return err
		}
		list = append(list, tss)
		return nil
	})
	return
}

// GetPreviousTSS returns the previous tss information
func (k Keeper) GetPreviousTSS(ctx sdk.Context) (val types.TSS, found bool) {
	tssList := k.GetAllTSS(ctx)
	if len(tssList) < 2 {
		return val, false
	}
	// Sort tssList by FinalizedZetaHeight
	sort.SliceStable(tssList, func(i, j int) bool {
		return tssList[i].FinalizedZetaHeight < tssList[j].FinalizedZetaHeight
	})
	return tssList[len(tssList)-2], true
}

func (k Keeper) CheckIfTssPubkeyHasBeenGenerated(ctx sdk.Context, tssPubkey string) (types.TSS, bool) {
	tssList := k.GetAllTSS(ctx)
	for _, tss := range tssList {
		if tss.TssPubkey == tssPubkey {
			return tss, true
		}
	}
	return types.TSS{}, false
}

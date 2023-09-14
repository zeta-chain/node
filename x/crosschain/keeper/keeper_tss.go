package keeper

import (
	"context"
	"fmt"
	"sort"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AppendTss appends a tss to the TSSHistoryKey store and update the current TSS to the latest one
func (k Keeper) AppendTss(ctx sdk.Context, tss types.TSS) {
	k.SetTSS(ctx, tss)
	k.SetTSSHistory(ctx, tss)
}

func (k Keeper) SetTssAndUpdateNonce(ctx sdk.Context, tss types.TSS) {
	k.SetTSS(ctx, tss)
	// initialize the nonces and pending nonces of all enabled chains
	supportedChains := k.zetaObserverKeeper.GetParams(ctx).GetSupportedChains()
	for _, chain := range supportedChains {
		chainNonce := types.ChainNonces{Index: chain.ChainName.String(), ChainId: chain.ChainId, Nonce: 0, FinalizedHeight: uint64(ctx.BlockHeight())}
		k.SetChainNonces(ctx, chainNonce)

		p := types.PendingNonces{
			NonceLow:  0,
			NonceHigh: 0,
			ChainId:   chain.ChainId,
			Tss:       tss.TssPubkey,
		}
		k.SetPendingNonces(ctx, p)
	}
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

// Queries

func (k Keeper) TSS(c context.Context, req *types.QueryGetTSSRequest) (*types.QueryGetTSSResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetTSS(ctx)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetTSSResponse{TSS: &val}, nil
}

// TssHistory Query historical list of TSS information
func (k Keeper) TssHistory(c context.Context, _ *types.QueryTssHistoryRequest) (*types.QueryTssHistoryResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	tssList := k.GetAllTSS(ctx)
	sort.SliceStable(tssList, func(i, j int) bool {
		return tssList[i].FinalizedZetaHeight < tssList[j].FinalizedZetaHeight
	})
	return &types.QueryTssHistoryResponse{TssList: tssList}, nil
}

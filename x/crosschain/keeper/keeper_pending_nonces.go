package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CRUD
func (k Keeper) SetNonceToCctx(ctx sdk.Context, nonceToCctx types.NonceToCctx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NonceToCctxKeyPrefix))
	b := k.cdc.MustMarshal(&nonceToCctx)
	store.Set(types.KeyPrefix(fmt.Sprintf("%s-%d-%d", nonceToCctx.Tss, nonceToCctx.ChainId, nonceToCctx.Nonce)), b)
}

func (k Keeper) GetNonceToCctx(ctx sdk.Context, tss string, chainID int64, nonce int64) (val types.NonceToCctx, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NonceToCctxKeyPrefix))

	b := store.Get(types.KeyPrefix(fmt.Sprintf("%s-%d-%d", tss, chainID, nonce)))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) RemoveNonceToCctx(ctx sdk.Context, nonceToCctx types.NonceToCctx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NonceToCctxKeyPrefix))
	store.Delete(types.KeyPrefix(fmt.Sprintf("%s-%d-%d", nonceToCctx.Tss, nonceToCctx.ChainId, nonceToCctx.Nonce)))
}

func (k Keeper) SetPendingNonces(ctx sdk.Context, pendingNonces types.PendingNonces) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PendingNoncesKeyPrefix))
	b := k.cdc.MustMarshal(&pendingNonces)
	store.Set(types.KeyPrefix(fmt.Sprintf("%s-%d", pendingNonces.Tss, pendingNonces.ChainId)), b)
}

func (k Keeper) GetPendingNonces(ctx sdk.Context, tss string, chainID int64) (val types.PendingNonces, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PendingNoncesKeyPrefix))

	b := store.Get(types.KeyPrefix(fmt.Sprintf("%s-%d", tss, chainID)))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetAllPendingNonces(ctx sdk.Context) (list []*types.PendingNonces) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PendingNoncesKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.PendingNonces
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, &val)
	}

	return
}

func (k Keeper) RemovePendingNonces(ctx sdk.Context, pendingNonces types.PendingNonces) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PendingNoncesKeyPrefix))
	store.Delete(types.KeyPrefix(fmt.Sprintf("%s-%d", pendingNonces.Tss, pendingNonces.ChainId)))
}

// utility
func (k Keeper) RemoveFromPendingNonces(ctx sdk.Context, tssPubkey string, chainID int64, nonce int64) {
	p, found := k.GetPendingNonces(ctx, tssPubkey, chainID)
	if found && nonce >= p.NonceLow && nonce <= p.NonceHigh {
		p.NonceLow = nonce + 1
		k.SetPendingNonces(ctx, p)
	}
}

func (k Keeper) PendingNoncesAll(c context.Context, req *types.QueryAllPendingNoncesRequest) (*types.QueryAllPendingNoncesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	list := k.GetAllPendingNonces(ctx)

	return &types.QueryAllPendingNoncesResponse{
		PendingNonces: list,
	}, nil
}

func (k Keeper) PendingNoncesByChain(c context.Context, req *types.QueryPendingNoncesByChainRequest) (*types.QueryPendingNoncesByChainResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	tss, found := k.GetTSS(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "tss not found")
	}
	list, found := k.GetPendingNonces(ctx, tss.TssPubkey, req.ChainId)
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("pending nonces not found for chain id : %d", req.ChainId))
	}

	return &types.QueryPendingNoncesByChainResponse{
		PendingNonces: list,
	}, nil
}

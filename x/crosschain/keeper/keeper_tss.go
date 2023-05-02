package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	TSSSingletonIndex = "TSSSingletonIndex"
)

// SetTSS set a specific tSS in the store from its index
func (k Keeper) SetTSS(ctx sdk.Context, tSS types.TSS) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TSSKey))
	b := k.cdc.MustMarshal(&tSS)
	store.Set(types.KeyPrefix(TSSSingletonIndex), b)
}

// GetTSS returns a tSS from its index
func (k Keeper) GetTSS(ctx sdk.Context) (val types.TSS, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TSSKey))

	b := store.Get(types.KeyPrefix(TSSSingletonIndex))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveTSS removes a tSS from the store
func (k Keeper) RemoveTSS(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TSSKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllTSS returns all tSS
func (k Keeper) GetAllTSS(ctx sdk.Context) (list []types.TSS) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TSSKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.TSS
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// Queries

func (k Keeper) TSSAll(c context.Context, req *types.QueryAllTSSRequest) (*types.QueryAllTSSResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var tSSs []*types.TSS
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	tSSStore := prefix.NewStore(store, types.KeyPrefix(types.TSSKey))

	pageRes, err := query.Paginate(tSSStore, req.Pagination, func(key []byte, value []byte) error {
		var tSS types.TSS
		if err := k.cdc.Unmarshal(value, &tSS); err != nil {
			return err
		}

		tSSs = append(tSSs, &tSS)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllTSSResponse{TSS: tSSs, Pagination: pageRes}, nil
}

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

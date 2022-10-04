package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SetReceive set a specific receive in the store from its index
func (k Keeper) SetReceive(ctx sdk.Context, receive types.Receive) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ReceiveKey))
	b := k.cdc.MustMarshal(&receive)
	store.Set(types.KeyPrefix(receive.Index), b)
}

// GetReceive returns a receive from its index
func (k Keeper) GetReceive(ctx sdk.Context, index string) (val types.Receive, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ReceiveKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveReceive removes a receive from the store
func (k Keeper) RemoveReceive(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ReceiveKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllReceive returns all receive
func (k Keeper) GetAllReceive(ctx sdk.Context) (list []types.Receive) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ReceiveKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Receive
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// Queries

func (k Keeper) ReceiveAll(c context.Context, req *types.QueryAllReceiveRequest) (*types.QueryAllReceiveResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var receives []*types.Receive
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	receiveStore := prefix.NewStore(store, types.KeyPrefix(types.ReceiveKey))

	pageRes, err := query.Paginate(receiveStore, req.Pagination, func(key []byte, value []byte) error {
		var receive types.Receive
		if err := k.cdc.Unmarshal(value, &receive); err != nil {
			return err
		}

		receives = append(receives, &receive)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllReceiveResponse{Receive: receives, Pagination: pageRes}, nil
}

func (k Keeper) Receive(c context.Context, req *types.QueryGetReceiveRequest) (*types.QueryGetReceiveResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetReceive(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetReceiveResponse{Receive: &val}, nil
}

package keeper

import (
	"context"
	math2 "math"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SetLastBlockHeight set a specific lastBlockHeight in the store from its index
func (k Keeper) SetLastBlockHeight(ctx sdk.Context, lastBlockHeight types.LastBlockHeight) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LastBlockHeightKey))
	b := k.cdc.MustMarshal(&lastBlockHeight)
	store.Set(types.KeyPrefix(lastBlockHeight.Index), b)
}

// GetLastBlockHeight returns a lastBlockHeight from its index
func (k Keeper) GetLastBlockHeight(ctx sdk.Context, index string) (val types.LastBlockHeight, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LastBlockHeightKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveLastBlockHeight removes a lastBlockHeight from the store
func (k Keeper) RemoveLastBlockHeight(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LastBlockHeightKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllLastBlockHeight returns all lastBlockHeight
func (k Keeper) GetAllLastBlockHeight(ctx sdk.Context) (list []types.LastBlockHeight) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LastBlockHeightKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.LastBlockHeight
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// Queries

func (k Keeper) LastBlockHeightAll(c context.Context, req *types.QueryAllLastBlockHeightRequest) (*types.QueryAllLastBlockHeightResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var lastBlockHeights []*types.LastBlockHeight
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	lastBlockHeightStore := prefix.NewStore(store, types.KeyPrefix(types.LastBlockHeightKey))

	pageRes, err := query.Paginate(lastBlockHeightStore, req.Pagination, func(key []byte, value []byte) error {
		var lastBlockHeight types.LastBlockHeight
		if err := k.cdc.Unmarshal(value, &lastBlockHeight); err != nil {
			return err
		}

		lastBlockHeights = append(lastBlockHeights, &lastBlockHeight)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllLastBlockHeightResponse{LastBlockHeight: lastBlockHeights, Pagination: pageRes}, nil
}

func (k Keeper) LastBlockHeight(c context.Context, req *types.QueryGetLastBlockHeightRequest) (*types.QueryGetLastBlockHeightResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetLastBlockHeight(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}
	if val.LastSendHeight < 0 || val.LastSendHeight >= math2.MaxInt64 {
		return nil, status.Error(codes.OutOfRange, "invalid last send height")
	}
	if val.LastReceiveHeight < 0 || val.LastReceiveHeight >= math2.MaxInt64 {
		return nil, status.Error(codes.OutOfRange, "invalid last recv height")
	}

	return &types.QueryGetLastBlockHeightResponse{LastBlockHeight: &val}, nil
}

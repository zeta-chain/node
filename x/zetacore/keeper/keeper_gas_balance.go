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

// SetGasBalance set a specific gasBalance in the store from its index
func (k Keeper) SetGasBalance(ctx sdk.Context, gasBalance types.GasBalance) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GasBalanceKey))
	b := k.cdc.MustMarshal(&gasBalance)
	store.Set(types.KeyPrefix(gasBalance.Index), b)
}

// GetGasBalance returns a gasBalance from its index
func (k Keeper) GetGasBalance(ctx sdk.Context, index string) (val types.GasBalance, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GasBalanceKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveGasBalance removes a gasBalance from the store
func (k Keeper) RemoveGasBalance(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GasBalanceKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllGasBalance returns all gasBalance
func (k Keeper) GetAllGasBalance(ctx sdk.Context) (list []types.GasBalance) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GasBalanceKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.GasBalance
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// Queries

func (k Keeper) GasBalanceAll(c context.Context, req *types.QueryAllGasBalanceRequest) (*types.QueryAllGasBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var gasBalances []*types.GasBalance
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	gasBalanceStore := prefix.NewStore(store, types.KeyPrefix(types.GasBalanceKey))

	pageRes, err := query.Paginate(gasBalanceStore, req.Pagination, func(key []byte, value []byte) error {
		var gasBalance types.GasBalance
		if err := k.cdc.Unmarshal(value, &gasBalance); err != nil {
			return err
		}

		gasBalances = append(gasBalances, &gasBalance)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllGasBalanceResponse{GasBalance: gasBalances, Pagination: pageRes}, nil
}

func (k Keeper) GasBalance(c context.Context, req *types.QueryGetGasBalanceRequest) (*types.QueryGetGasBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetGasBalance(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetGasBalanceResponse{GasBalance: &val}, nil
}

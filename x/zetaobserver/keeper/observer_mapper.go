package keeper

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetaobserver/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SetObserverMapper(ctx sdk.Context, om types.ObserverMapper) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ObserverMapperKey))
	om.Index = fmt.Sprintf("%s-%s", om.ObserverChain.String(), om.ObservationType.String())
	b := k.cdc.MustMarshal(&om)
	store.Set([]byte(om.Index), b)
}

func (k Keeper) GetObserverMapper(ctx sdk.Context, chain, obsType string) (val types.ObserverMapper, found bool) {
	index := fmt.Sprintf("%s-%s", chain, obsType)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ObserverMapperKey))
	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetAllObserverMappers(ctx sdk.Context) (mappers []*types.ObserverMapper) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ObserverMapperKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.ObserverMapper
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		mappers = append(mappers, &val)
	}
	return
}

//Queries

func (k Keeper) ObserversByChainAndType(goCtx context.Context, req *types.QueryObserversByChainAndTypeRequest) (*types.QueryObserversByChainAndTypeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	mapper, _ := k.GetObserverMapper(ctx, req.ObservationChain, req.ObservationType)
	//if !isFound {
	//	return &types.QueryObserversByChainAndTypeResponse{ObserverMapper: "Not Found"}, nil
	//}
	return &types.QueryObserversByChainAndTypeResponse{ObserverMapper: mapper.String()}, nil
}

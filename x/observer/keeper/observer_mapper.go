package keeper

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SetObserverMapper(ctx sdk.Context, om *types.ObserverMapper) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ObserverMapperKey))
	om.Index = fmt.Sprintf("%s-%s", om.ObserverChain.String(), om.ObservationType.String())
	b := k.cdc.MustMarshal(om)
	store.Set([]byte(om.Index), b)
}

func (k Keeper) GetObserverMapper(ctx sdk.Context, chain types.ObserverChain, obsType types.ObservationType) (val types.ObserverMapper, found bool) {
	index := fmt.Sprintf("%s-%s", chain.String(), obsType)
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
func (k Keeper) GetAllObserverMappersForAddress(ctx sdk.Context, address string) (mappers []*types.ObserverMapper) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ObserverMapperKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.ObserverMapper
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		addToList := false
		for _, addr := range val.ObserverList {
			if addr == address {
				addToList = true
			}
		}
		if addToList {
			mappers = append(mappers, &val)
		}
	}
	return
}

// Tx

func (k msgServer) AddObserver(goCtx context.Context, msg *types.MsgAddObserver) (*types.MsgAddObserverResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.IsValidator(ctx, msg.Creator)
	if err != nil {
		return nil, err
	}
	err = k.CheckObserverDelegation(ctx, msg.Creator, msg.ObserverChain, msg.ObservationType)
	if err != nil {
		return nil, err
	}
	if !k.IsChainSupported(ctx, msg.ObserverChain) {
		return nil, errors.Wrap(types.ErrSupportedChains, fmt.Sprintf("chain %s", msg.ObserverChain.String()))
	}
	k.AddObserverToMapper(ctx,
		msg.ObserverChain,
		msg.ObservationType,
		msg.Creator)
	AddObserverEvent(ctx, msg.Creator, msg.ObserverChain, msg.ObservationType)

	return &types.MsgAddObserverResponse{}, nil
}

//Queries

func (k Keeper) ObserversByChainAndType(goCtx context.Context, req *types.QueryObserversByChainAndTypeRequest) (*types.QueryObserversByChainAndTypeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	mapper, _ := k.GetObserverMapper(ctx, types.ParseCommonChaintoObservationChain(req.ObservationChain), types.ParseStringToObservationType(req.ObservationType))
	return &types.QueryObserversByChainAndTypeResponse{Observers: mapper.ObserverList}, nil
}

func (k Keeper) AllObserverMappers(goCtx context.Context, req *types.QueryAllObserverMappersRequest) (*types.QueryAllObserverMappersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	mappers := k.GetAllObserverMappers(ctx)
	return &types.QueryAllObserverMappersResponse{ObserverMappers: mappers}, nil
}

// Utils

func (k Keeper) GetAllObserverAddresses(ctx sdk.Context) []string {
	var val []string
	mappers := k.GetAllObserverMappers(ctx)
	for _, mapper := range mappers {
		val = append(val, mapper.ObserverList...)
	}
	allKeys := make(map[string]bool)
	var dedupedList []string
	for _, item := range val {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			dedupedList = append(dedupedList, item)
		}
	}
	return dedupedList
}

func (k Keeper) AddObserverToMapper(ctx sdk.Context, chain types.ObserverChain, obsType types.ObservationType, address string) {
	mapper, found := k.GetObserverMapper(ctx, chain, obsType)
	if !found {
		k.SetObserverMapper(ctx, &types.ObserverMapper{
			Index:           "",
			ObserverChain:   chain,
			ObservationType: obsType,
			ObserverList:    []string{address},
		})
		return
	}
	for _, addr := range mapper.ObserverList {
		if addr == address {
			return
		}
	}
	mapper.ObserverList = append(mapper.ObserverList, address)
	k.SetObserverMapper(ctx, &mapper)
}

package keeper

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetObserverMapperIndex(chain *common.Chain, obs types.ObservationType) string {
	return fmt.Sprintf("%d-%s", chain.ChainId, obs.String())
}

func (k Keeper) SetObserverMapper(ctx sdk.Context, om *types.ObserverMapper) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ObserverMapperKey))
	om.Index = GetObserverMapperIndex(om.ObserverChain, om.ObservationType)
	b := k.cdc.MustMarshal(om)
	store.Set([]byte(om.Index), b)
}

func (k Keeper) GetObserverMapper(ctx sdk.Context, chain *common.Chain, obsType types.ObservationType) (val types.ObserverMapper, found bool) {
	index := GetObserverMapperIndex(chain, obsType)
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
	chain := k.GetParams(ctx).GetChainFromChainID(msg.ChainId)
	if chain == nil {
		return nil, sdkerrors.Wrap(types.ErrSupportedChains, fmt.Sprintf("ChainID %d", msg.ChainId))
	}
	err = k.CheckObserverDelegation(ctx, msg.Creator, chain, msg.ObservationType)
	if err != nil {
		return nil, err
	}
	k.AddObserverToMapper(ctx,
		chain,
		msg.ObservationType,
		msg.Creator)

	return &types.MsgAddObserverResponse{}, nil
}

//Queries

func (k Keeper) ObserversByChainAndType(goCtx context.Context, req *types.QueryObserversByChainAndTypeRequest) (*types.QueryObserversByChainAndTypeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	// TODO move parsing to client
	chainName := common.ParseStringToObserverChain(req.ObservationChain)
	chain := k.GetParams(ctx).GetChainFromChainName(chainName)
	if chain == nil {
		return &types.QueryObserversByChainAndTypeResponse{}, types.ErrSupportedChains
	}
	mapper, found := k.GetObserverMapper(ctx, chain, types.ParseStringToObservationType(req.ObservationType))
	if !found {
		return &types.QueryObserversByChainAndTypeResponse{}, types.ErrObserverNotPresent
	}
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

func (k Keeper) AddObserverToMapper(ctx sdk.Context, chain *common.Chain, obsType types.ObservationType, address string) {
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
	// Return if duplicate
	for _, addr := range mapper.ObserverList {
		if addr == address {
			return
		}
	}
	mapper.ObserverList = append(mapper.ObserverList, address)
	k.SetObserverMapper(ctx, &mapper)
}

package keeper

import (
	"context"
	"fmt"
	"math"

	cosmoserrors "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetObserverMapperIndex(chain *common.Chain) string {
	return fmt.Sprintf("%d", chain.ChainId)
}

func (k Keeper) SetLastObserverCount(ctx sdk.Context, lbc *types.LastObserverCount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LastBlockObserverCountKey))
	b := k.cdc.MustMarshal(lbc)
	store.Set([]byte{0}, b)
}

func (k Keeper) GetLastObserverCount(ctx sdk.Context) (val types.LastObserverCount, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LastBlockObserverCountKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) SetObserverMapper(ctx sdk.Context, om *types.ObserverMapper) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ObserverMapperKey))
	om.Index = GetObserverMapperIndex(om.ObserverChain)
	b := k.cdc.MustMarshal(om)
	store.Set([]byte(om.Index), b)
}

func (k Keeper) GetObserverMapper(ctx sdk.Context, chain *common.Chain) (val types.ObserverMapper, found bool) {
	index := GetObserverMapperIndex(chain)
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

// AddObserver adds in a new observer to the store.It can be executed using an admin policy account
// Once added, the function also resets keygen and pauses inbound so that a new TSS can be generated.
func (k msgServer) AddObserver(goCtx context.Context, msg *types.MsgAddObserver) (*types.MsgAddObserverResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.GetParams(ctx).GetAdminPolicyAccount(types.Policy_Type_add_observer) {
		return &types.MsgAddObserverResponse{}, types.ErrNotAuthorizedPolicy
	}
	pubkey, err := common.NewPubKey(msg.ZetaclientGranteePubkey)
	if err != nil {
		return &types.MsgAddObserverResponse{}, cosmoserrors.Wrap(types.ErrInvalidPubKey, err.Error())
	}
	granteeAddress, err := common.GetAddressFromPubkeyString(msg.ZetaclientGranteePubkey)
	if err != nil {
		return &types.MsgAddObserverResponse{}, cosmoserrors.Wrap(types.ErrInvalidPubKey, err.Error())
	}
	k.DisableInboundOnly(ctx)
	// AddNodeAccountOnly flag usage
	// True: adds observer into the Node Account list but returns without adding to the observer list
	// False: adds observer to the observer list, and not the node account list
	// Inbound is disabled in both cases and needs to be enabled manually using an admin TX
	if msg.AddNodeAccountOnly {
		pubkeySet := common.PubKeySet{Secp256k1: pubkey, Ed25519: ""}
		k.SetNodeAccount(ctx, types.NodeAccount{
			Operator:       msg.ObserverAddress,
			GranteeAddress: granteeAddress.String(),
			GranteePubkey:  &pubkeySet,
			NodeStatus:     types.NodeStatus_Active,
		})
		k.SetKeygen(ctx, types.Keygen{BlockNumber: math.MaxInt64})
		return &types.MsgAddObserverResponse{}, nil
	}

	observerMappers := k.GetAllObserverMappers(ctx)
	totalObserverCountCurrentBlock := uint64(0)
	for _, mapper := range observerMappers {
		mapper.ObserverList = append(mapper.ObserverList, msg.ObserverAddress)
		totalObserverCountCurrentBlock += uint64(len(mapper.ObserverList))
		k.SetObserverMapper(ctx, mapper)
	}
	k.SetLastObserverCount(ctx, &types.LastObserverCount{Count: totalObserverCountCurrentBlock})
	EmitEventAddObserver(ctx, totalObserverCountCurrentBlock, msg.ObserverAddress, granteeAddress.String(), msg.ZetaclientGranteePubkey)
	return &types.MsgAddObserverResponse{}, nil
}

//Queries

func (k Keeper) ObserversByChain(goCtx context.Context, req *types.QueryObserversByChainRequest) (*types.QueryObserversByChainResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO move parsing to client
	// https://github.com/zeta-chain/node/issues/867

	chainName := common.ParseChainName(req.ObservationChain)
	chain := k.GetParams(ctx).GetChainFromChainName(chainName)
	if chain == nil {
		return &types.QueryObserversByChainResponse{}, types.ErrSupportedChains
	}
	mapper, found := k.GetObserverMapper(ctx, chain)
	if !found {
		return &types.QueryObserversByChainResponse{}, types.ErrObserverNotPresent
	}
	return &types.QueryObserversByChainResponse{Observers: mapper.ObserverList}, nil
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

func (k Keeper) AddObserverToMapper(ctx sdk.Context, chain *common.Chain, address string) {
	mapper, found := k.GetObserverMapper(ctx, chain)
	if !found {
		k.SetObserverMapper(ctx, &types.ObserverMapper{
			Index:         "",
			ObserverChain: chain,
			ObserverList:  []string{address},
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

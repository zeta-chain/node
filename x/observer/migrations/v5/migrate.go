package v5

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// observerKeeper prevents circular dependency
type observerKeeper interface {
	GetParams(ctx sdk.Context) types.Params
	SetParams(ctx sdk.Context, params types.Params)
	GetChainParamsList(ctx sdk.Context) (params types.ChainParamsList, found bool)
	SetChainParamsList(ctx sdk.Context, params types.ChainParamsList)
	StoreKey() storetypes.StoreKey
	Codec() codec.BinaryCodec
}

func MigrateStore(ctx sdk.Context, observerKeeper observerKeeper) error {
	if err := MigrateObserverMapper(ctx, observerKeeper.StoreKey(), observerKeeper.Codec()); err != nil {
		return err
	}
	return MigrateObserverParams(ctx, observerKeeper)
}

func MigrateObserverMapper(ctx sdk.Context, observerStoreKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	var legacyObserverMappers []*types.ObserverMapper
	legacyObserverMapperStore := prefix.NewStore(ctx.KVStore(observerStoreKey), types.KeyPrefix(types.ObserverMapperKey))
	iterator := sdk.KVStorePrefixIterator(legacyObserverMapperStore, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.ObserverMapper
		cdc.MustUnmarshal(iterator.Value(), &val)
		legacyObserverMappers = append(legacyObserverMappers, &val)
	}

	// We can safely assume that the observer list is the same for all the observer mappers
	observerList := legacyObserverMappers[0].ObserverList

	storelastBlockObserverCount := prefix.NewStore(ctx.KVStore(observerStoreKey), types.KeyPrefix(types.LastBlockObserverCountKey))
	b := cdc.MustMarshal(&types.LastObserverCount{Count: uint64(len(observerList)), LastChangeHeight: ctx.BlockHeight()})
	storelastBlockObserverCount.Set([]byte{0}, b)

	storeObserverSet := prefix.NewStore(ctx.KVStore(observerStoreKey), types.KeyPrefix(types.ObserverSetKey))
	b = cdc.MustMarshal(&types.ObserverSet{ObserverList: observerList})
	storeObserverSet.Set([]byte{0}, b)

	for _, legacyObserverMapper := range legacyObserverMappers {
		legacyObserverMapperStore.Delete(types.KeyPrefix(legacyObserverMapper.Index))
	}
	return nil
}

// MigrateObserverParams migrates the observer params to the chain params
// the function assumes that each observer params entry has a corresponding chain params entry
// if the chain is not found, the observer params entry is ignored because it is considered as not supported
func MigrateObserverParams(ctx sdk.Context, observerKeeper observerKeeper) error {
	chainParamsList, found := observerKeeper.GetChainParamsList(ctx)
	if !found {
		// no chain params found, nothing to migrate
		return nil
	}

	// search for the observer params with chain params entry
	observerParams := observerKeeper.GetParams(ctx).ObserverParams
	for _, observerParam := range observerParams {
		for i := range chainParamsList.ChainParams {
			// if the chain is found, update the chain params with the observer params
			if chainParamsList.ChainParams[i].ChainId == observerParam.Chain.ChainId {
				chainParamsList.ChainParams[i].MinObserverDelegation = observerParam.MinObserverDelegation
				chainParamsList.ChainParams[i].BallotThreshold = observerParam.BallotThreshold
				chainParamsList.ChainParams[i].IsSupported = observerParam.IsSupported
				break
			}
		}
	}

	observerKeeper.SetChainParamsList(ctx, chainParamsList)
	return nil
}

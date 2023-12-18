package v4

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
	GetCoreParamsList(ctx sdk.Context) (params types.CoreParamsList, found bool)
	SetCoreParamsList(ctx sdk.Context, params types.CoreParamsList)
	StoreKey() storetypes.StoreKey
	Codec() codec.BinaryCodec
}

func MigrateStore(ctx sdk.Context, observerKeeper observerKeeper) error {
	if err := MigrateCrosschainFlags(ctx, observerKeeper.StoreKey(), observerKeeper.Codec()); err != nil {
		return err
	}
	if err := MigrateObserverParams(ctx, observerKeeper); err != nil {
		return err
	}
	return nil
}

func MigrateCrosschainFlags(ctx sdk.Context, observerStoreKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	newCrossChainFlags := types.DefaultCrosschainFlags()
	var val types.LegacyCrosschainFlags
	store := prefix.NewStore(ctx.KVStore(observerStoreKey), types.KeyPrefix(types.CrosschainFlagsKey))
	b := store.Get([]byte{0})
	if b != nil {
		cdc.MustUnmarshal(b, &val)
		if val.GasPriceIncreaseFlags != nil {
			newCrossChainFlags.GasPriceIncreaseFlags = val.GasPriceIncreaseFlags
		}
		newCrossChainFlags.IsOutboundEnabled = val.IsOutboundEnabled
		newCrossChainFlags.IsInboundEnabled = val.IsInboundEnabled
	}
	b, err := cdc.Marshal(newCrossChainFlags)
	if err != nil {
		return err
	}
	store.Set([]byte{0}, b)
	return nil
}

// MigrateObserverParams migrates the observer params to the core params
// the function assumes that each oberver params entry has a corresponding core params entry
// if the chain is not found, the observer params entry is ignored because it is considered as not supported
func MigrateObserverParams(ctx sdk.Context, observerKeeper observerKeeper) error {
	coreParamsList, found := observerKeeper.GetCoreParamsList(ctx)
	if !found {
		// no core params found, nothing to migrate
		return nil
	}

	// search for the observer params with core params entry
	observerParams := observerKeeper.GetParams(ctx).ObserverParams
	for _, observerParam := range observerParams {
		for i := range coreParamsList.CoreParams {
			// if the chain is found, update the core params with the observer params
			if coreParamsList.CoreParams[i].ChainId == observerParam.Chain.ChainId {
				coreParamsList.CoreParams[i].MinObserverDelegation = observerParam.MinObserverDelegation
				coreParamsList.CoreParams[i].BallotThreshold = observerParam.BallotThreshold
				coreParamsList.CoreParams[i].IsSupported = observerParam.IsSupported
				break
			}
		}
	}

	observerKeeper.SetCoreParamsList(ctx, coreParamsList)
	return nil
}

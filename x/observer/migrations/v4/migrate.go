package v4

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func MigrateStore(ctx sdk.Context, observerStoreKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	if err := MigrateCrosschainFlags(ctx, observerStoreKey, cdc); err != nil {
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

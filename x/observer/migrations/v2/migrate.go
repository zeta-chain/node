package v2

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// MigrateStore migrates the x/observer module state from the consensus version 1 to 2
/* This migration adds a
- new permission flag to the observer module called IsOutboundEnabled
*/
func MigrateStore(
	ctx sdk.Context,
	observerStoreKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
) error {
	store := prefix.NewStore(ctx.KVStore(observerStoreKey), types.KeyPrefix(types.PermissionFlagsKey))
	b := cdc.MustMarshal(&types.PermissionFlags{
		IsInboundEnabled:  true,
		IsOutboundEnabled: true,
	})
	store.Set([]byte{0}, b)
	return nil
}

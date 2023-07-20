package v2

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// MigrateStore migrates the x/crosschain module state from the consensus version 1 to 2

func MigrateStore(
	ctx sdk.Context,
	observerKeeper types.ZetaObserverKeeper,
	crossChainStoreKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
) error {
	store := prefix.NewStore(ctx.KVStore(crossChainStoreKey), types.KeyPrefix(LegacyNodeAccountKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	var nodeAccounts []observerTypes.NodeAccount
	for ; iterator.Valid(); iterator.Next() {
		var val observerTypes.NodeAccount
		cdc.MustUnmarshal(iterator.Value(), &val)
		nodeAccounts = append(nodeAccounts, val)
	}
	for _, nodeAccount := range nodeAccounts {
		observerKeeper.SetNodeAccount(ctx, nodeAccount)
	}

	return nil
}

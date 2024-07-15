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
// This migration moves some data from the cross chain store to the observer store.
// The data moved is the node accounts, permission flags and keygen.
func MigrateStore(
	ctx sdk.Context,
	observerKeeper types.ObserverKeeper,
	crossChainStoreKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
) error {
	// Using New Types from observer module as the structure is the same
	var nodeAccounts []observerTypes.NodeAccount
	var crosschainFlags observerTypes.CrosschainFlags
	var keygen observerTypes.Keygen
	writePermissionFlags := false
	writeKeygen := false

	// Fetch data from cross chain store using the legacy keys directly
	store := prefix.NewStore(ctx.KVStore(crossChainStoreKey), types.KeyPrefix(LegacyNodeAccountKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val observerTypes.NodeAccount
		cdc.MustUnmarshal(iterator.Value(), &val)
		nodeAccounts = append(nodeAccounts, val)
	}

	store = prefix.NewStore(ctx.KVStore(crossChainStoreKey), types.KeyPrefix(LegacyKeygenKey))
	b := store.Get([]byte{0})
	if b != nil {
		cdc.MustUnmarshal(b, &keygen)
		writeKeygen = true
	}

	store = prefix.NewStore(ctx.KVStore(crossChainStoreKey), types.KeyPrefix(LegacyPermissionFlagsKey))
	b = store.Get([]byte{0})
	if b != nil {
		cdc.MustUnmarshal(b, &crosschainFlags)
		writePermissionFlags = true
	}

	// Write data to observer store using the new keys
	if nodeAccounts != nil {
		for _, nodeAccount := range nodeAccounts {
			observerKeeper.SetNodeAccount(ctx, nodeAccount)
		}
	}
	if writeKeygen {
		observerKeeper.SetKeygen(ctx, keygen)
	}
	if writePermissionFlags {
		observerKeeper.SetCrosschainFlags(ctx, crosschainFlags)
	}

	allObservers, found := observerKeeper.GetObserverSet(ctx)
	if !found {
		return observerTypes.ErrObserverSetNotFound
	}
	totalObserverCountCurrentBlock := allObservers.LenUint()

	observerKeeper.SetLastObserverCount(ctx, &observerTypes.LastObserverCount{
		// #nosec G115 always positive
		Count:            totalObserverCountCurrentBlock,
		LastChangeHeight: ctx.BlockHeight(),
	})

	return nil
}

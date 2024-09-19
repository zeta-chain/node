package v8

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/observer/types"
)

type observerKeeper interface {
	SetChainNonces(ctx sdk.Context, chainNonces types.ChainNonces)
	GetAllChainNonces(ctx sdk.Context) []types.ChainNonces
	StoreKey() storetypes.StoreKey
}

// MigrateStore migrates the x/observer module state from the consensus version 7 to 8
// It updates the indexing for chain nonces object to use chain ID instead of chain name
func MigrateStore(ctx sdk.Context, observerKeeper observerKeeper) error {
	updateChainNonceIndexing(ctx, observerKeeper)
	return nil
}

// updateChainNonceIndexing updates the chain nonces object to use chain ID instead of chain name
func updateChainNonceIndexing(
	ctx sdk.Context,
	observerKeeper observerKeeper,
) {
	// Iterate all chain nonces object in the store
	chainNonces := observerKeeper.GetAllChainNonces(ctx)
	for _, chainNonce := range chainNonces {
		// set again the object, since SetChainNonces uses chain ID as index, the indexing will be updated
		observerKeeper.SetChainNonces(ctx, chainNonce)

		// remove the old object
		removeChainNoncesLegacy(ctx, observerKeeper, chainNonce.Index)
	}
}

// removeChainNoncesLegacy removes a chainNonces from the store from index
func removeChainNoncesLegacy(ctx sdk.Context, observerKeeper observerKeeper, index string) {
	store := prefix.NewStore(ctx.KVStore(observerKeeper.StoreKey()), types.KeyPrefix(types.ChainNoncesKey))
	store.Delete(types.KeyPrefix(index))
}

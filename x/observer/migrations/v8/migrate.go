package v8

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// MigrateStore migrates the x/observer module state from the consensus version 7 to 8
// It updates the indexing for chain nonces object to use chain ID instead of chain name
func MigrateStore(ctx sdk.Context, observerKeeper types.ObserverKeeper) error {
	UpdateChainNonceIndexing(ctx, observerKeeper)

	return nil
}

// UpdateChainNonceIndexing updates the chain nonces object to use chain ID instead of chain name
func UpdateChainNonceIndexing(
	ctx sdk.Context,
	observerKeeper types.ObserverKeeper,
) {
	chainNonces := observerKeeper.GetAllChainNonces(ctx)
	for _, chainNonce := range chainNonces {
		observerKeeper.DeleteChainNonces(ctx, chainNonce.ChainName)
		chainNonce.ChainName = chainNonce.ChainID.String()
		observerKeeper.SetChainNonces(ctx, chainNonce)
	}
}

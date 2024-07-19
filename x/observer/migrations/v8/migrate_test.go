package v8_test

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	v8 "github.com/zeta-chain/zetacore/x/observer/migrations/v8"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"testing"
)

var chainNonces = []types.ChainNonces{
	{
		Index:   chains.BitcoinMainnet.ChainName.String(),
		ChainId: chains.BitcoinMainnet.ChainId,
		Nonce:   1,
	},
	{
		Index:   chains.BitcoinTestnet.ChainName.String(),
		ChainId: chains.BitcoinTestnet.ChainId,
		Nonce:   2,
	},
	{
		Index:   chains.Ethereum.ChainName.String(),
		ChainId: chains.Ethereum.ChainId,
		Nonce:   3,
	},
	{
		Index:   chains.Sepolia.ChainName.String(),
		ChainId: chains.Sepolia.ChainId,
		Nonce:   4,
	},
	{
		Index:   chains.BscMainnet.ChainName.String(),
		ChainId: chains.BscMainnet.ChainId,
		Nonce:   5,
	},
	{
		Index:   chains.BscTestnet.ChainName.String(),
		ChainId: chains.BscTestnet.ChainId,
		Nonce:   6,
	},
	{
		Index:   chains.Polygon.ChainName.String(),
		ChainId: chains.Polygon.ChainId,
		Nonce:   7,
	},
	{
		Index:   chains.Amoy.ChainName.String(),
		ChainId: chains.Amoy.ChainId,
		Nonce:   8,
	},
	{
		Index:   chains.BitcoinRegtest.String(),
		ChainId: chains.BitcoinRegtest.ChainId,
		Nonce:   9,
	},
	{
		Index:   chains.GoerliLocalnet.String(),
		ChainId: chains.GoerliLocalnet.ChainId,
		Nonce:   10,
	},
}

func TestMigrateStore(t *testing.T) {
	t.Run("MigrateStore", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// set chain nonces
		for _, chainNonce := range chainNonces {
			setChainNoncesLegacy(ctx, *k, chainNonce)
		}

		// there are 10 chain nonces in the store
		chainNonces := k.GetAllChainNonces(ctx)
		require.Len(t, chainNonces, 10)

		// no chain nonces can be found in the store
		for _, chainNonce := range chainNonces {
			_, found := k.GetChainNonces(ctx, chainNonce.ChainId)
			require.False(t, found)
		}

		// migrate the store
		err := v8.MigrateStore(ctx, *k)
		require.NoError(t, err)

		// there are 10 chain nonces in the store
		chainNonces = k.GetAllChainNonces(ctx)
		require.Len(t, chainNonces, 10)

		// all chain nonces can be found in the store
		for _, chainNonce := range chainNonces {
			_, found := k.GetChainNonces(ctx, chainNonce.ChainId)
			require.True(t, found)
		}
	})

}

// setChainNoncesLegacy set a specific chainNonces in the store from its index
func setChainNoncesLegacy(ctx sdk.Context, observerKeeper keeper.Keeper, chainNonces types.ChainNonces) {
	store := prefix.NewStore(ctx.KVStore(observerKeeper.StoreKey()), types.KeyPrefix(types.ChainNoncesKey))
	b := observerKeeper.Codec().MustMarshal(&chainNonces)
	store.Set(types.KeyPrefix(chainNonces.Index), b)
}

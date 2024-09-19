package v8_test

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/x/observer/keeper"
	v8 "github.com/zeta-chain/node/x/observer/migrations/v8"
	"github.com/zeta-chain/node/x/observer/types"
	"testing"
)

var chainNoncesArray = []types.ChainNonces{
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
	t.Run("can migrate chain nonces", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// set chain nonces
		for _, chainNonce := range chainNoncesArray {
			setChainNoncesLegacy(ctx, *k, chainNonce)
		}

		// there are 10 chain nonces in the store
		allChainNonces := k.GetAllChainNonces(ctx)
		require.Len(t, allChainNonces, 10)

		// no chain nonces can be found in the store using the new indexing
		for _, chainNonces := range allChainNonces {
			_, found := k.GetChainNonces(ctx, chainNonces.ChainId)
			require.False(t, found)
		}

		// migrate the store
		err := v8.MigrateStore(ctx, *k)
		require.NoError(t, err)

		// there are 10 chain nonces in the store
		allChainNonces = k.GetAllChainNonces(ctx)
		require.Len(t, allChainNonces, 10)

		chainIDMap := make(map[int64]struct{})

		// all chain nonces can be found in the store
		for _, chainNonces := range allChainNonces {
			// chain all chain IDs are different
			_, found := chainIDMap[chainNonces.ChainId]
			require.False(t, found)
			chainIDMap[chainNonces.ChainId] = struct{}{}

			// check value
			retrievedChainNonces, found := k.GetChainNonces(ctx, chainNonces.ChainId)
			require.True(t, found)
			require.Contains(t, chainNoncesArray, retrievedChainNonces)
		}
	})

	t.Run("migrate nothing with empty array", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		allChainNonces := k.GetAllChainNonces(ctx)
		require.Len(t, allChainNonces, 0)

		// migrate the store
		err := v8.MigrateStore(ctx, *k)
		require.NoError(t, err)

		allChainNonces = k.GetAllChainNonces(ctx)
		require.Len(t, allChainNonces, 0)
	})
}

// setChainNoncesLegacy set a specific chainNonces in the store from its index
func setChainNoncesLegacy(ctx sdk.Context, observerKeeper keeper.Keeper, chainNonces types.ChainNonces) {
	store := prefix.NewStore(ctx.KVStore(observerKeeper.StoreKey()), types.KeyPrefix(types.ChainNoncesKey))
	b := observerKeeper.Codec().MustMarshal(&chainNonces)
	store.Set(types.KeyPrefix(chainNonces.Index), b)
}

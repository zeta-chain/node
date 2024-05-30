package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/authority/types"
)

func TestKeeper_ChainInfo(t *testing.T) {
	t.Run("invalid request", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)

		_, err := k.ChainInfo(ctx, nil)
		require.ErrorContains(t, err, "invalid request")
	})

	t.Run("chain info not found", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)

		_, err := k.ChainInfo(ctx, &types.QueryGetChainInfoRequest{})
		require.ErrorContains(t, err, "chain info not found")
	})

	t.Run("can retrieve chain info", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)

		chainInfo := sample.ChainInfo(42)
		k.SetChainInfo(ctx, chainInfo)

		res, err := k.ChainInfo(ctx, &types.QueryGetChainInfoRequest{})
		require.NoError(t, err)
		require.Equal(t, chainInfo, res.ChainInfo)
	})
}

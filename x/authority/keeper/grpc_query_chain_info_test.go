package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/authority/types"
)

func TestKeeper_ChainInfo(t *testing.T) {
	t.Run("invalid request", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)

		_, err := k.ChainInfo(ctx, nil)
		require.ErrorContains(t, err, "invalid request")
	})

	t.Run("chain info not found", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)

		res, err := k.ChainInfo(ctx, &types.QueryGetChainInfoRequest{})
		require.NoError(t, err)
		require.Equal(t, res, &types.QueryGetChainInfoResponse{
			ChainInfo: types.ChainInfo{
				Chains: []chains.Chain{},
			},
		})
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

package keeper_test

import (
	"github.com/zeta-chain/node/pkg/chains"
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
)

func TestKeeper_GetAllForeignCoins(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	fc := sample.ForeignCoins(t, sample.EthAddress().Hex())
	fc.ForeignChainId = chains.LocalZetaChainID
	k.GetFungibleKeeper().SetForeignCoins(ctx, fc)

	res := k.GetAllForeignCoins(ctx)
	require.Equal(t, 1, len(res))
}

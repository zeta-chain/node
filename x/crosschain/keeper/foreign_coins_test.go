package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

func TestKeeper_GetAllForeignCoins(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	fc := sample.ForeignCoins(t, sample.EthAddress().Hex())
	fc.ForeignChainId = 101
	k.GetFungibleKeeper().SetForeignCoins(ctx, fc)

	res, err := k.GetAllForeignCoins(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
}

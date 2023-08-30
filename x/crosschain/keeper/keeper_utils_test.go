package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	testkeeper "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"testing"
)

func TestKeeper_IncreaseCctxGasPrice(t *testing.T) {
	k, ctx := testkeeper.CrosschainKeeper(t)

	t.Run("can increase gas", func(t *testing.T) {
		// sample cctx
		cctx := *sample.CrossChainTx(t, "foo")
		previousGasPrice, ok := math.NewIntFromString(cctx.GetCurrentOutTxParam().OutboundTxGasPrice)
		require.True(t, ok)

		// increase gas price
		err := k.IncreaseCctxGasPrice(ctx, cctx, math.NewInt(42))
		require.NoError(t, err)

		// can retrieve cctx
		cctx, found := k.GetCrossChainTx(ctx, "foo")
		require.True(t, found)

		// gas price increased
		currentGasPrice, ok := math.NewIntFromString(cctx.GetCurrentOutTxParam().OutboundTxGasPrice)
		require.True(t, ok)
		require.True(t, currentGasPrice.Equal(previousGasPrice.Add(math.NewInt(42))))
	})

	t.Run("fail if invalid cctx", func(t *testing.T) {
		cctx := *sample.CrossChainTx(t, "foo")
		cctx.GetCurrentOutTxParam().OutboundTxGasPrice = "invalid"
		err := k.IncreaseCctxGasPrice(ctx, cctx, math.NewInt(42))
		require.Error(t, err)
	})

}

package types_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// TODO: Complete tests for this file
// https://github.com/zeta-chain/node/issues/2669

func TestParseGatewayEvent(t *testing.T) {

}

func TestParseGatewayWithdrawalEvent(t *testing.T) {

}

func TestParseGatewayCallEvent(t *testing.T) {

}

func TestParseGatewayWithdrawAndCallEvent(t *testing.T) {

}

func TestNewWithdrawalInbound(t *testing.T) {
	t.Run("fail if sender chain ID is not valid", func(t *testing.T) {
		_, ctx, _, _ := keepertest.CrosschainKeeper(t)
		ctx = ctx.WithChainID("invalidChainID")

		fc := sample.ForeignCoins(t, sample.EthAddress().Hex())

		_, err := types.NewWithdrawalInbound(
			ctx,
			sample.EthAddress().Hex(),
			fc.CoinType,
			fc.Asset,
			nil,
			chains.GoerliLocalnet,
			big.NewInt(1000),
		)

		require.ErrorContains(t, err, " failed to convert chainID")
	})

}

func TestNewCallInbound(t *testing.T) {
	t.Run("fail if sender chain ID is not valid", func(t *testing.T) {
		_, ctx, _, _ := keepertest.CrosschainKeeper(t)
		ctx = ctx.WithChainID("invalidChainID")

		_, err := types.NewCallInbound(
			ctx,
			sample.EthAddress().Hex(),
			nil,
			chains.GoerliLocalnet,
			big.NewInt(1000),
		)

		require.ErrorContains(t, err, " failed to convert chainID")
	})

}

func TestNewWithdrawAndCallInbound(t *testing.T) {
	t.Run("fail if sender chain ID is not valid", func(t *testing.T) {
		_, ctx, _, _ := keepertest.CrosschainKeeper(t)
		ctx = ctx.WithChainID("invalidChainID")

		fc := sample.ForeignCoins(t, sample.EthAddress().Hex())

		_, err := types.NewWithdrawAndCallInbound(
			ctx,
			sample.EthAddress().Hex(),
			fc.CoinType,
			fc.Asset,
			nil,
			chains.GoerliLocalnet,
			big.NewInt(1000),
		)

		require.ErrorContains(t, err, " failed to convert chainID")
	})

	t.Run("fail if receiver address can't be decoded", func(t *testing.T) {

	})
}

package signer

import (
	"math/big"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestGasFromCCTX(t *testing.T) {
	logger := zerolog.New(zerolog.NewTestWriter(t))

	makeCCTX := func(gasLimit uint64, price, priorityFee string) *types.CrossChainTx {
		cctx := getCCTX(t)
		cctx.GetOutboundParams()[0].CallOptions.GasLimit = gasLimit
		cctx.GetOutboundParams()[0].GasPrice = price
		cctx.GetOutboundParams()[0].GasPriorityFee = priorityFee

		return cctx
	}

	erc20Withdraw := func(gasLimit uint64, price, priorityFee string) *types.CrossChainTx {
		cctx := makeCCTX(gasLimit, price, priorityFee)
		cctx.InboundParams.IsCrossChainCall = false
		cctx.InboundParams.CoinType = coin.CoinType_ERC20
		cctx.OutboundParams[0].CoinType = coin.CoinType_ERC20
		require.Len(t, cctx.GetOutboundParams(), 1)

		return cctx
	}

	gasWithdraw := func(gasLimit uint64, price, priorityFee string) *types.CrossChainTx {
		cctx := makeCCTX(gasLimit, price, priorityFee)
		cctx.InboundParams.IsCrossChainCall = false
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		cctx.OutboundParams[0].CoinType = coin.CoinType_Gas
		require.Len(t, cctx.GetOutboundParams(), 1)

		return cctx
	}

	gasWithdrawWithCall := func(gasLimit uint64, price, priorityFee string) *types.CrossChainTx {
		cctx := makeCCTX(gasLimit, price, priorityFee)
		cctx.InboundParams.IsCrossChainCall = true
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		cctx.OutboundParams[0].CoinType = coin.CoinType_Gas
		require.Len(t, cctx.GetOutboundParams(), 1)

		return cctx
	}

	gasWithdrawRevert := func(gasLimit uint64, price, priorityFee string) *types.CrossChainTx {
		cctx := makeCCTX(gasLimit, price, priorityFee)
		cctx.InboundParams.IsCrossChainCall = false
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		cctx.OutboundParams[0].CoinType = coin.CoinType_Gas
		cctx.OutboundParams = append(cctx.OutboundParams, cctx.OutboundParams[0])

		return cctx
	}

	gasWithdrawRevertWithCall := func(gasLimit uint64, price, priorityFee string) *types.CrossChainTx {
		cctx := makeCCTX(gasLimit, price, priorityFee)
		cctx.InboundParams.IsCrossChainCall = false
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		cctx.OutboundParams[0].CoinType = coin.CoinType_Gas
		cctx.OutboundParams = append(cctx.OutboundParams, cctx.OutboundParams[0])
		cctx.RevertOptions.CallOnRevert = true

		return cctx
	}

	for _, tt := range []struct {
		name          string
		cctx          *types.CrossChainTx
		errorContains string
		assert        func(t *testing.T, g Gas)
	}{
		{
			name: "gas limit is set to min if below min apply for erc20 withdraw",
			cctx: erc20Withdraw(contractCallMinGasLimit-1, gwei(3).String(), ""),
			assert: func(t *testing.T, g Gas) {
				assert.True(t, g.isLegacy())
				assertGasEquals(t, Gas{
					Limit:       contractCallMinGasLimit,
					Price:       gwei(3),
					PriorityFee: gwei(0),
				}, g)
			},
		},
		{
			name: "gas limit is set to min if below min doesn't apply for gas withdraw",
			cctx: gasWithdraw(contractCallMinGasLimit-1, gwei(3).String(), ""),
			assert: func(t *testing.T, g Gas) {
				assert.True(t, g.isLegacy())
				assertGasEquals(t, Gas{
					Limit:       contractCallMinGasLimit - 1,
					Price:       gwei(3),
					PriorityFee: gwei(0),
				}, g)
			},
		},
		{
			name: "gas limit is set to min if below min apply for gas withdraw with call",
			cctx: gasWithdrawWithCall(contractCallMinGasLimit-1, gwei(3).String(), ""),
			assert: func(t *testing.T, g Gas) {
				assert.True(t, g.isLegacy())
				assertGasEquals(t, Gas{
					Limit:       contractCallMinGasLimit,
					Price:       gwei(3),
					PriorityFee: gwei(0),
				}, g)
			},
		},
		{
			name: "gas limit is set to min if below min doesn't apply for gas withdraw revert",
			cctx: gasWithdrawRevert(contractCallMinGasLimit-1, gwei(3).String(), ""),
			assert: func(t *testing.T, g Gas) {
				assert.True(t, g.isLegacy())
				assertGasEquals(t, Gas{
					Limit:       contractCallMinGasLimit - 1,
					Price:       gwei(3),
					PriorityFee: gwei(0),
				}, g)
			},
		},
		{
			name: "gas limit is set to min if below min apply for gas withdraw revert with call",
			cctx: gasWithdrawRevertWithCall(contractCallMinGasLimit-1, gwei(3).String(), ""),
			assert: func(t *testing.T, g Gas) {
				assert.True(t, g.isLegacy())
				assertGasEquals(t, Gas{
					Limit:       contractCallMinGasLimit,
					Price:       gwei(3),
					PriorityFee: gwei(0),
				}, g)
			},
		},
		{
			name: "pre London gas logic",
			cctx: makeCCTX(21_000, gwei(3).String(), ""),
			assert: func(t *testing.T, g Gas) {
				assert.True(t, g.isLegacy())
				assertGasEquals(t, Gas{
					Limit:       21_000,
					Price:       gwei(3),
					PriorityFee: gwei(0),
				}, g)
			},
		},
		{
			name: "post London gas logic",
			cctx: makeCCTX(21_000, gwei(4).String(), gwei(1).String()),
			assert: func(t *testing.T, g Gas) {
				assert.False(t, g.isLegacy())
				assertGasEquals(t, Gas{
					Limit:       21_000,
					Price:       gwei(4),
					PriorityFee: gwei(1),
				}, g)
			},
		},
		{
			name: "gas is too high, force to the ceiling",
			cctx: makeCCTX(maxGasLimit+200, gwei(4).String(), gwei(1).String()),
			assert: func(t *testing.T, g Gas) {
				assert.False(t, g.isLegacy())
				assertGasEquals(t, Gas{
					Limit:       maxGasLimit,
					Price:       gwei(4),
					PriorityFee: gwei(1),
				}, g)
			},
		},
		{
			name:          "priority fee is invalid",
			cctx:          makeCCTX(123_000, gwei(4).String(), "oopsie"),
			errorContains: "unable to parse priorityFee",
		},
		{
			name:          "priority fee is negative",
			cctx:          makeCCTX(123_000, gwei(4).String(), "-1"),
			errorContains: "unable to parse priorityFee: big.Int is negative",
		},
		{
			name: "gasPrice is less than priorityFee",
			cctx: makeCCTX(123_000, gwei(4).String(), gwei(5).String()),
			assert: func(t *testing.T, g Gas) {
				assert.False(t, g.isLegacy())
				assertGasEquals(t, Gas{
					Limit:       123_000,
					Price:       gwei(4),
					PriorityFee: gwei(4),
				}, g)
			},
		},
		{
			name:          "gasPrice is invalid",
			cctx:          makeCCTX(123_000, "hello", gwei(5).String()),
			errorContains: "unable to parse gasPrice",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			g, err := gasFromCCTX(tt.cctx, logger)
			if tt.errorContains != "" {
				assert.ErrorContains(t, err, tt.errorContains)
				return
			}

			assert.NoError(t, err)
			assert.NoError(t, g.validate())
			tt.assert(t, g)
		})
	}

	t.Run("empty priority fee", func(t *testing.T) {
		gas := Gas{
			Limit:       123_000,
			Price:       gwei(4),
			PriorityFee: nil,
		}

		assert.Error(t, gas.validate())
	})
}

func assertGasEquals(t *testing.T, expected, actual Gas) {
	assert.Equal(t, int64(expected.Limit), int64(actual.Limit), "gas limit")
	assert.Equal(t, expected.Price.Int64(), actual.Price.Int64(), "max fee per unit")
	assert.Equal(t, expected.PriorityFee.Int64(), actual.PriorityFee.Int64(), "priority fee per unit")
}

func gwei(i int64) *big.Int {
	const g = 1_000_000_000
	return big.NewInt(i * g)
}

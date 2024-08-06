package signer

import (
	"math/big"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestGasFromCCTX(t *testing.T) {
	logger := zerolog.New(zerolog.NewTestWriter(t))

	makeCCTX := func(gasLimit uint64, price, priorityFee string) *types.CrossChainTx {
		cctx := getCCTX(t)
		cctx.GetOutboundParams()[0].GasLimit = gasLimit
		cctx.GetOutboundParams()[0].GasPrice = price
		cctx.GetOutboundParams()[0].GasPriorityFee = priorityFee

		return cctx
	}

	for _, tt := range []struct {
		name          string
		cctx          *types.CrossChainTx
		errorContains string
		assert        func(t *testing.T, g Gas)
	}{

		{
			name: "legacy: gas is too low",
			cctx: makeCCTX(minGasLimit-200, gwei(2).String(), ""),
			assert: func(t *testing.T, g Gas) {
				assert.True(t, g.IsLegacy())
				assertGasEquals(t, Gas{
					limit:       minGasLimit,
					price:       gwei(2),
					priorityFee: gwei(0),
				}, g)
			},
		},
		{
			name: "london: gas is too low",
			cctx: makeCCTX(minGasLimit-200, gwei(2).String(), gwei(1).String()),
			assert: func(t *testing.T, g Gas) {
				assert.False(t, g.IsLegacy())
				assertGasEquals(t, Gas{
					limit:       minGasLimit,
					price:       gwei(2),
					priorityFee: gwei(1),
				}, g)

				// gasPrice=2, priorityFee=1, so baseFee=1
				// gasFeeCap = 2*baseFee + priorityFee = 2 + 1 = 3
				assert.Equal(t, gwei(3).Int64(), g.GasFeeCap().Int64())
			},
		},
		{
			name: "pre London gas logic",
			cctx: makeCCTX(minGasLimit+100, gwei(3).String(), ""),
			assert: func(t *testing.T, g Gas) {
				assert.True(t, g.IsLegacy())
				assertGasEquals(t, Gas{
					limit:       100_100,
					price:       gwei(3),
					priorityFee: gwei(0),
				}, g)
			},
		},
		{
			name: "post London gas logic",
			cctx: makeCCTX(minGasLimit+200, gwei(4).String(), gwei(1).String()),
			assert: func(t *testing.T, g Gas) {
				assert.False(t, g.IsLegacy())
				assertGasEquals(t, Gas{
					limit:       100_200,
					price:       gwei(4),
					priorityFee: gwei(1),
				}, g)
			},
		},
		{
			name: "gas is too high, force to the ceiling",
			cctx: makeCCTX(maxGasLimit+200, gwei(4).String(), gwei(1).String()),
			assert: func(t *testing.T, g Gas) {
				assert.False(t, g.IsLegacy())
				assertGasEquals(t, Gas{
					limit:       maxGasLimit,
					price:       gwei(4),
					priorityFee: gwei(1),
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
			name:          "gasPrice is less than priorityFee",
			cctx:          makeCCTX(123_000, gwei(4).String(), gwei(5).String()),
			errorContains: "gasPrice (4000000000) is less than priorityFee (5000000000)",
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
			tt.assert(t, g)
		})
	}
}

func assertGasEquals(t *testing.T, expected, actual Gas) {
	assert.Equal(t, int64(expected.Limit()), int64(actual.Limit()), "gas limit")
	assert.Equal(t, expected.GasPrice().Int64(), actual.GasPrice().Int64(), "gas price")
	assert.Equal(t, expected.GasFeeCap().Int64(), actual.GasFeeCap().Int64(), "max fee cap")
	assert.Equal(t, expected.PriorityFee().Int64(), actual.PriorityFee().Int64(), "priority fee")
}

func gwei(i int64) *big.Int {
	const g = 1_000_000_000
	return big.NewInt(i * g)
}

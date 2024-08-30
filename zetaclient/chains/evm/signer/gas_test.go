package signer

import (
	"math/big"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/node/x/crosschain/types"
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
				assert.True(t, g.isLegacy())
				assertGasEquals(t, Gas{
					Limit:       minGasLimit,
					PriorityFee: gwei(0),
					Price:       gwei(2),
				}, g)
			},
		},
		{
			name: "london: gas is too low",
			cctx: makeCCTX(minGasLimit-200, gwei(2).String(), gwei(1).String()),
			assert: func(t *testing.T, g Gas) {
				assert.False(t, g.isLegacy())
				assertGasEquals(t, Gas{
					Limit:       minGasLimit,
					Price:       gwei(2),
					PriorityFee: gwei(1),
				}, g)
			},
		},
		{
			name: "pre London gas logic",
			cctx: makeCCTX(minGasLimit+100, gwei(3).String(), ""),
			assert: func(t *testing.T, g Gas) {
				assert.True(t, g.isLegacy())
				assertGasEquals(t, Gas{
					Limit:       100_100,
					Price:       gwei(3),
					PriorityFee: gwei(0),
				}, g)
			},
		},
		{
			name: "post London gas logic",
			cctx: makeCCTX(minGasLimit+200, gwei(4).String(), gwei(1).String()),
			assert: func(t *testing.T, g Gas) {
				assert.False(t, g.isLegacy())
				assertGasEquals(t, Gas{
					Limit:       100_200,
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

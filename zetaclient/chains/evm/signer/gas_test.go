package signer

import (
	"math/big"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func Test_makeGasFromCCTX(t *testing.T) {
	logger := zerolog.Nop()

	cctx1 := getCCTX(t)
	cctx1.GetOutboundParams()[0].GasLimit = MinGasLimit - 200
	cctx1.GetOutboundParams()[0].GasPrice = gwei(2).String()

	cctx2 := getCCTX(t)
	cctx2.GetOutboundParams()[0].GasLimit = 200_000
	cctx2.GetOutboundParams()[0].GasPrice = gwei(3).String()

	cctx3 := getCCTX(t)
	cctx3.GetOutboundParams()[0].GasLimit = 2_000_000
	cctx3.GetOutboundParams()[0].GasPrice = gwei(3).String()

	for _, tt := range []struct {
		name          string
		cctx          *types.CrossChainTx
		priorityFee   *big.Int
		errorContains string
		assert        func(t *testing.T, g Gas)
	}{
		{
			name:          "priority fee is nil",
			cctx:          cctx1,
			priorityFee:   nil,
			errorContains: "priorityFee is nil",
		},
		{
			name:        "gas is too low",
			cctx:        cctx1,
			priorityFee: gwei(1),
			assert: func(t *testing.T, g Gas) {
				assert.False(t, g.isLegacy())
				assertGasEquals(t, Gas{
					Limit:              MinGasLimit,
					PriorityFeePerUnit: gwei(1),
					MaxFeePerUnit:      gwei(2),
				}, g)
			},
		},
		{
			name:        "as is, no surprises",
			cctx:        cctx2,
			priorityFee: gwei(2),
			assert: func(t *testing.T, g Gas) {
				assert.False(t, g.isLegacy())
				assertGasEquals(t, Gas{
					Limit:              200_000,
					PriorityFeePerUnit: gwei(2),
					MaxFeePerUnit:      gwei(3),
				}, g)
			},
		},
		{
			name:        "pre London gas logic",
			cctx:        cctx2,
			priorityFee: gwei(0),
			assert: func(t *testing.T, g Gas) {
				assert.True(t, g.isLegacy())
				assertGasEquals(t, Gas{
					Limit:              200_000,
					PriorityFeePerUnit: gwei(0),
					MaxFeePerUnit:      gwei(3),
				}, g)
			},
		},
		{
			name:        "gas is too high, force to the ceiling",
			cctx:        cctx3,
			priorityFee: gwei(2),
			assert: func(t *testing.T, g Gas) {
				assert.False(t, g.isLegacy())
				assertGasEquals(t, Gas{
					Limit:              MaxGasLimit,
					PriorityFeePerUnit: gwei(2),
					MaxFeePerUnit:      gwei(3),
				}, g)
			},
		},
		{
			name:          "priority fee is too high",
			cctx:          cctx1,
			priorityFee:   gwei(200),
			errorContains: "less than priorityFee",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			g, err := makeGasFromCCTX(tt.cctx, tt.priorityFee, logger)
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
	assert.Equal(t, int64(expected.Limit), int64(actual.Limit), "gas limit")
	assert.Equal(t, expected.MaxFeePerUnit.Int64(), actual.MaxFeePerUnit.Int64(), "max fee per unit")
	assert.Equal(t, expected.PriorityFeePerUnit.Int64(), actual.PriorityFeePerUnit.Int64(), "priority fee per unit")
}

func gwei(i int64) *big.Int {
	const g = 1_000_000_000
	return big.NewInt(i * g)
}

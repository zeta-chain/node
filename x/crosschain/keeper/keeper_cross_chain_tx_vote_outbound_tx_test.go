package keeper_test

import (
	"errors"
	"math/big"
	"math/rand"
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	testkeeper "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

func TestKeeper_FundGasStabilityPoolFromRemainingFees(t *testing.T) {
	r := rand.New(rand.NewSource(42))

	tt := []struct {
		name                                  string
		gasLimit                              uint64
		gasUsed                               uint64
		effectiveGasPrice                     math.Int
		fundStabilityPoolReturn               error
		expectFundStabilityPoolCall           bool
		fundStabilityPoolExpectedRemainingFee *big.Int
		isError                               bool
	}{
		{
			name:                        "no call if gasLimit is 0",
			gasLimit:                    0,
			gasUsed:                     42,
			effectiveGasPrice:           math.NewInt(42),
			expectFundStabilityPoolCall: false,
		},
		{
			name:                        "no call if gasUsed is 0",
			gasLimit:                    42,
			gasUsed:                     0,
			effectiveGasPrice:           math.NewInt(42),
			expectFundStabilityPoolCall: false,
		},
		{
			name:                        "no call if effectiveGasPrice is 0",
			gasLimit:                    42,
			gasUsed:                     42,
			effectiveGasPrice:           math.NewInt(0),
			expectFundStabilityPoolCall: false,
		},
		{
			name:              "should return error if gas limit is less than gas used",
			gasLimit:          41,
			gasUsed:           42,
			effectiveGasPrice: math.NewInt(42),
			isError:           true,
		},
		{
			name:                                  "should call fund stability pool with correct remaining fees",
			gasLimit:                              100,
			gasUsed:                               90,
			effectiveGasPrice:                     math.NewInt(100),
			fundStabilityPoolReturn:               nil,
			expectFundStabilityPoolCall:           true,
			fundStabilityPoolExpectedRemainingFee: big.NewInt(500), // (100-90)*100 = 1000 * 50% = 500
		},
		{
			name:                                  "should return error if fund stability pool returns error",
			gasLimit:                              100,
			gasUsed:                               90,
			effectiveGasPrice:                     math.NewInt(100),
			fundStabilityPoolReturn:               errors.New("fund stability pool error"),
			expectFundStabilityPoolCall:           true,
			fundStabilityPoolExpectedRemainingFee: big.NewInt(500),
			isError:                               true,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			k, ctx := testkeeper.CrosschainKeeperAllMocks(t)
			fungibleMock := testkeeper.GetCrosschainFungibleMock(t, k)

			// OutboundTxParams
			outbound := sample.OutboundTxParams(r)
			outbound.OutboundTxGasLimit = tc.gasLimit
			outbound.OutboundTxGasUsed = tc.gasUsed
			outbound.OutboundTxEffectiveGasPrice = tc.effectiveGasPrice

			if tc.expectFundStabilityPoolCall {
				fungibleMock.On(
					"FundGasStabilityPool", ctx, int64(42), tc.fundStabilityPoolExpectedRemainingFee,
				).Return(tc.fundStabilityPoolReturn)
			}

			err := k.FundGasStabilityPoolFromRemainingFees(ctx, *outbound, 42)
			if tc.isError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			fungibleMock.AssertExpectations(t)
		})
	}
}

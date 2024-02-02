package keeper_test

import (
	"errors"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/assert"
	testkeeper "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestKeeper_CheckAndUpdateCctxGasPrice(t *testing.T) {
	sampleTimestamp := time.Now()
	retryIntervalReached := sampleTimestamp.Add(observertypes.DefaultGasPriceIncreaseFlags.RetryInterval + time.Second)
	retryIntervalNotReached := sampleTimestamp.Add(observertypes.DefaultGasPriceIncreaseFlags.RetryInterval - time.Second)

	tt := []struct {
		name                                   string
		cctx                                   types.CrossChainTx
		flags                                  observertypes.GasPriceIncreaseFlags
		blockTimestamp                         time.Time
		medianGasPrice                         uint64
		withdrawFromGasStabilityPoolReturn     error
		expectWithdrawFromGasStabilityPoolCall bool
		expectedGasPriceIncrease               math.Uint
		expectedAdditionalFees                 math.Uint
		isError                                bool
	}{
		{
			name: "can update gas price when retry interval is reached",
			cctx: types.CrossChainTx{
				Index: "a1",
				CctxStatus: &types.Status{
					LastUpdateTimestamp: sampleTimestamp.Unix(),
				},
				OutboundTxParams: []*types.OutboundTxParams{
					{
						ReceiverChainId:    42,
						OutboundTxGasLimit: 1000,
						OutboundTxGasPrice: "100",
					},
				},
			},
			flags:                                  observertypes.DefaultGasPriceIncreaseFlags,
			blockTimestamp:                         retryIntervalReached,
			medianGasPrice:                         50,
			withdrawFromGasStabilityPoolReturn:     nil,
			expectWithdrawFromGasStabilityPoolCall: true,
			expectedGasPriceIncrease:               math.NewUint(50),    // 100% medianGasPrice
			expectedAdditionalFees:                 math.NewUint(50000), // gasLimit * increase
		},
		{
			name: "can update gas price at max limit",
			cctx: types.CrossChainTx{
				Index: "a2",
				CctxStatus: &types.Status{
					LastUpdateTimestamp: sampleTimestamp.Unix(),
				},
				OutboundTxParams: []*types.OutboundTxParams{
					{
						ReceiverChainId:    42,
						OutboundTxGasLimit: 1000,
						OutboundTxGasPrice: "100",
					},
				},
			},
			flags: observertypes.GasPriceIncreaseFlags{
				EpochLength:             100,
				RetryInterval:           time.Minute * 10,
				GasPriceIncreasePercent: 200, // Increase gas price to 100+50*2 = 200
				GasPriceIncreaseMax:     400, // Max gas price is 50*4 = 200
			},
			blockTimestamp:                         retryIntervalReached,
			medianGasPrice:                         50,
			withdrawFromGasStabilityPoolReturn:     nil,
			expectWithdrawFromGasStabilityPoolCall: true,
			expectedGasPriceIncrease:               math.NewUint(100),    // 200% medianGasPrice
			expectedAdditionalFees:                 math.NewUint(100000), // gasLimit * increase
		},
		{
			name: "default gas price increase limit used if not defined",
			cctx: types.CrossChainTx{
				Index: "a3",
				CctxStatus: &types.Status{
					LastUpdateTimestamp: sampleTimestamp.Unix(),
				},
				OutboundTxParams: []*types.OutboundTxParams{
					{
						ReceiverChainId:    42,
						OutboundTxGasLimit: 1000,
						OutboundTxGasPrice: "100",
					},
				},
			},
			flags: observertypes.GasPriceIncreaseFlags{
				EpochLength:             100,
				RetryInterval:           time.Minute * 10,
				GasPriceIncreasePercent: 100,
				GasPriceIncreaseMax:     0, // Limit should not be reached
			},
			blockTimestamp:                         retryIntervalReached,
			medianGasPrice:                         50,
			withdrawFromGasStabilityPoolReturn:     nil,
			expectWithdrawFromGasStabilityPoolCall: true,
			expectedGasPriceIncrease:               math.NewUint(50),    // 100% medianGasPrice
			expectedAdditionalFees:                 math.NewUint(50000), // gasLimit * increase
		},
		{
			name: "skip if max limit reached",
			cctx: types.CrossChainTx{
				Index: "b0",
				CctxStatus: &types.Status{
					LastUpdateTimestamp: sampleTimestamp.Unix(),
				},
				OutboundTxParams: []*types.OutboundTxParams{
					{
						ReceiverChainId:    42,
						OutboundTxGasLimit: 1000,
						OutboundTxGasPrice: "100",
					},
				},
			},
			flags: observertypes.GasPriceIncreaseFlags{
				EpochLength:             100,
				RetryInterval:           time.Minute * 10,
				GasPriceIncreasePercent: 200, // Increase gas price to 100+50*2 = 200
				GasPriceIncreaseMax:     300, // Max gas price is 50*3 = 150
			},
			blockTimestamp:                         retryIntervalReached,
			medianGasPrice:                         50,
			expectWithdrawFromGasStabilityPoolCall: false,
			expectedGasPriceIncrease:               math.NewUint(0),
			expectedAdditionalFees:                 math.NewUint(0),
		},
		{
			name: "skip if gas price is not set",
			cctx: types.CrossChainTx{
				Index: "b1",
				CctxStatus: &types.Status{
					LastUpdateTimestamp: sampleTimestamp.Unix(),
				},
				OutboundTxParams: []*types.OutboundTxParams{
					{
						ReceiverChainId:    42,
						OutboundTxGasLimit: 100,
						OutboundTxGasPrice: "",
					},
				},
			},
			flags:                                  observertypes.DefaultGasPriceIncreaseFlags,
			blockTimestamp:                         retryIntervalReached,
			medianGasPrice:                         100,
			expectWithdrawFromGasStabilityPoolCall: false,
			expectedGasPriceIncrease:               math.NewUint(0),
			expectedAdditionalFees:                 math.NewUint(0),
		},
		{
			name: "skip if gas limit is not set",
			cctx: types.CrossChainTx{
				Index: "b2",
				CctxStatus: &types.Status{
					LastUpdateTimestamp: sampleTimestamp.Unix(),
				},
				OutboundTxParams: []*types.OutboundTxParams{
					{
						ReceiverChainId:    42,
						OutboundTxGasLimit: 0,
						OutboundTxGasPrice: "100",
					},
				},
			},
			flags:                                  observertypes.DefaultGasPriceIncreaseFlags,
			blockTimestamp:                         retryIntervalReached,
			medianGasPrice:                         100,
			expectWithdrawFromGasStabilityPoolCall: false,
			expectedGasPriceIncrease:               math.NewUint(0),
			expectedAdditionalFees:                 math.NewUint(0),
		},
		{
			name: "skip if retry interval is not reached",
			cctx: types.CrossChainTx{
				Index: "b3",
				CctxStatus: &types.Status{
					LastUpdateTimestamp: sampleTimestamp.Unix(),
				},
				OutboundTxParams: []*types.OutboundTxParams{
					{
						ReceiverChainId:    42,
						OutboundTxGasLimit: 0,
						OutboundTxGasPrice: "100",
					},
				},
			},
			flags:                                  observertypes.DefaultGasPriceIncreaseFlags,
			blockTimestamp:                         retryIntervalNotReached,
			medianGasPrice:                         100,
			expectWithdrawFromGasStabilityPoolCall: false,
			expectedGasPriceIncrease:               math.NewUint(0),
			expectedAdditionalFees:                 math.NewUint(0),
		},
		{
			name: "returns error if can't find median gas price",
			cctx: types.CrossChainTx{
				Index: "c1",
				CctxStatus: &types.Status{
					LastUpdateTimestamp: sampleTimestamp.Unix(),
				},
				OutboundTxParams: []*types.OutboundTxParams{
					{
						ReceiverChainId:    42,
						OutboundTxGasLimit: 1000,
						OutboundTxGasPrice: "100",
					},
				},
			},
			flags:                                  observertypes.DefaultGasPriceIncreaseFlags,
			expectWithdrawFromGasStabilityPoolCall: false,
			blockTimestamp:                         retryIntervalReached,
			medianGasPrice:                         0,
			isError:                                true,
		},
		{
			name: "returns error if can't withdraw from gas stability pool",
			cctx: types.CrossChainTx{
				Index: "c2",
				CctxStatus: &types.Status{
					LastUpdateTimestamp: sampleTimestamp.Unix(),
				},
				OutboundTxParams: []*types.OutboundTxParams{
					{
						ReceiverChainId:    42,
						OutboundTxGasLimit: 1000,
						OutboundTxGasPrice: "100",
					},
				},
			},
			flags:                                  observertypes.DefaultGasPriceIncreaseFlags,
			blockTimestamp:                         retryIntervalReached,
			medianGasPrice:                         50,
			expectWithdrawFromGasStabilityPoolCall: true,
			expectedGasPriceIncrease:               math.NewUint(50),    // 100% medianGasPrice
			expectedAdditionalFees:                 math.NewUint(50000), // gasLimit * increase
			withdrawFromGasStabilityPoolReturn:     errors.New("withdraw error"),
			isError:                                true,
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			k, ctx := testkeeper.CrosschainKeeperAllMocks(t)
			fungibleMock := testkeeper.GetCrosschainFungibleMock(t, k)
			chainID := tc.cctx.GetCurrentOutTxParam().ReceiverChainId
			previousGasPrice, err := tc.cctx.GetCurrentOutTxParam().GetGasPrice()
			if err != nil {
				previousGasPrice = 0
			}

			// set median gas price if not zero
			if tc.medianGasPrice != 0 {
				k.SetGasPrice(ctx, types.GasPrice{
					ChainId:     chainID,
					Prices:      []uint64{tc.medianGasPrice},
					MedianIndex: 0,
				})

				// ensure median gas price is set
				medianGasPrice, isFound := k.GetMedianGasPriceInUint(ctx, chainID)
				assert.True(t, isFound)
				assert.True(t, medianGasPrice.Equal(math.NewUint(tc.medianGasPrice)))
			}

			// set block timestamp
			ctx = ctx.WithBlockTime(tc.blockTimestamp)

			if tc.expectWithdrawFromGasStabilityPoolCall {
				fungibleMock.On(
					"WithdrawFromGasStabilityPool", ctx, chainID, tc.expectedAdditionalFees.BigInt(),
				).Return(tc.withdrawFromGasStabilityPoolReturn)
			}

			// check and update gas price
			gasPriceIncrease, feesPaid, err := k.CheckAndUpdateCctxGasPrice(ctx, tc.cctx, tc.flags)

			if tc.isError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// check values
			assert.True(t, gasPriceIncrease.Equal(tc.expectedGasPriceIncrease), "expected %s, got %s", tc.expectedGasPriceIncrease.String(), gasPriceIncrease.String())
			assert.True(t, feesPaid.Equal(tc.expectedAdditionalFees), "expected %s, got %s", tc.expectedAdditionalFees.String(), feesPaid.String())

			// check cctx
			if !tc.expectedGasPriceIncrease.IsZero() {
				cctx, found := k.GetCrossChainTx(ctx, tc.cctx.Index)
				assert.True(t, found)
				newGasPrice, err := cctx.GetCurrentOutTxParam().GetGasPrice()
				assert.NoError(t, err)
				assert.EqualValues(t, tc.expectedGasPriceIncrease.AddUint64(previousGasPrice).Uint64(), newGasPrice, "%d - %d", tc.expectedGasPriceIncrease.Uint64(), previousGasPrice)
				assert.EqualValues(t, tc.blockTimestamp.Unix(), cctx.CctxStatus.LastUpdateTimestamp)
			}
		})
	}
}

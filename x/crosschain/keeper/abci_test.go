package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	testkeeper "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestKeeper_CheckAndUpdateCctxGasPrice(t *testing.T) {
	k, ctx := testkeeper.CrosschainKeeper(t)
	sampleTimestamp := time.Now()
	retryIntervalReached := sampleTimestamp.Add(keeper.RetryInterval + time.Second)
	retryIntervalNotReached := sampleTimestamp.Add(keeper.RetryInterval - time.Second)

	tt := []struct {
		name                     string
		cctx                     types.CrossChainTx
		blockTimestamp           time.Time
		medianGasPrice           uint64
		expectedGasPriceIncrease math.Uint
		expectedFeesPaid         math.Uint
		isError                  bool
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
			blockTimestamp:           retryIntervalReached,
			medianGasPrice:           50,
			expectedGasPriceIncrease: math.NewUint(50),    // 100% medianGasPrice
			expectedFeesPaid:         math.NewUint(50000), // gasLimit * increase
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
			blockTimestamp:           retryIntervalReached,
			medianGasPrice:           100,
			expectedGasPriceIncrease: math.NewUint(0),
			expectedFeesPaid:         math.NewUint(0),
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
			blockTimestamp:           retryIntervalReached,
			medianGasPrice:           100,
			expectedGasPriceIncrease: math.NewUint(0),
			expectedFeesPaid:         math.NewUint(0),
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
			blockTimestamp:           retryIntervalNotReached,
			medianGasPrice:           100,
			expectedGasPriceIncrease: math.NewUint(0),
			expectedFeesPaid:         math.NewUint(0),
		},
		{
			name: "returns error if can't find median gas price",
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
			blockTimestamp: retryIntervalReached,
			medianGasPrice: 0,
			isError:        true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			chainID := tc.cctx.GetCurrentOutTxParam().ReceiverChainId

			// set median gas price if not zero
			if tc.medianGasPrice != 0 {
				k.SetGasPrice(ctx, types.GasPrice{
					ChainId:     chainID,
					Prices:      []uint64{tc.medianGasPrice},
					MedianIndex: 0,
				})

				// ensure median gas price is set
				medianGasPrice, isFound := k.GetMedianGasPriceInUint(ctx, chainID)
				require.True(t, isFound)
				require.True(t, medianGasPrice.Equal(math.NewUint(tc.medianGasPrice)))
			}

			// set block timestamp
			ctx = ctx.WithBlockTime(tc.blockTimestamp)

			// check and update gas price
			gasPriceIncrease, feesPaid, err := k.CheckAndUpdateCctxGasPrice(ctx, tc.cctx)

			if tc.isError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// check values
			require.True(t, gasPriceIncrease.Equal(tc.expectedGasPriceIncrease), "expected %s, got %s", tc.expectedGasPriceIncrease.String(), gasPriceIncrease.String())
			require.True(t, feesPaid.Equal(tc.expectedFeesPaid), "expected %s, got %s", tc.expectedFeesPaid.String(), feesPaid.String())
		})
	}
}

func TestKeeper_IncreaseCctxGasPrice(t *testing.T) {
	k, ctx := testkeeper.CrosschainKeeper(t)

	t.Run("can increase gas", func(t *testing.T) {
		// sample cctx
		cctx := *sample.CrossChainTx(t, "foo")
		previousGasPrice, ok := math.NewIntFromString(cctx.GetCurrentOutTxParam().OutboundTxGasPrice)
		require.True(t, ok)

		// increase gas price
		err := k.IncreaseCctxGasPrice(ctx, cctx, math.NewUint(42))
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
		err := k.IncreaseCctxGasPrice(ctx, cctx, math.NewUint(42))
		require.Error(t, err)
	})

}

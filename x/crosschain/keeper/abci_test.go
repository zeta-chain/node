package keeper_test

import (
	"errors"
	"testing"
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	testkeeper "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestKeeper_IterateAndUpdateCctxGasPrice(t *testing.T) {
	k, ctx, _, zk := testkeeper.CrosschainKeeper(t)

	// updateFuncMap tracks the calls done with cctx index
	updateFuncMap := make(map[string]struct{})

	// failMap gives the cctx index that should fail
	failMap := make(map[string]struct{})

	// updateFunc mocks the update function and keep track of the calls done with cctx index
	updateFunc := func(
		ctx sdk.Context,
		k keeper.Keeper,
		cctx types.CrossChainTx,
		flags observertypes.GasPriceIncreaseFlags,
	) (math.Uint, math.Uint, error) {
		if _, ok := failMap[cctx.Index]; ok {
			return math.NewUint(0), math.NewUint(0), errors.New("failed")
		}

		updateFuncMap[cctx.Index] = struct{}{}
		return math.NewUint(10), math.NewUint(10), nil
	}

	// add some evm and non-evm chains
	supportedChains := []*chains.Chain{
		{ChainId: chains.EthChain.ChainId},
		{ChainId: chains.BtcMainnetChain.ChainId},
		{ChainId: chains.BscMainnetChain.ChainId},
		{ChainId: chains.ZetaChainMainnet.ChainId},
	}

	// set pending cctx
	tss := sample.Tss()
	zk.ObserverKeeper.SetTSS(ctx, tss)
	createCctxWithNonceRange(t, ctx, *k, 10, 15, chains.EthChain.ChainId, tss, zk)
	createCctxWithNonceRange(t, ctx, *k, 20, 25, chains.BtcMainnetChain.ChainId, tss, zk)
	createCctxWithNonceRange(t, ctx, *k, 30, 35, chains.BscMainnetChain.ChainId, tss, zk)
	createCctxWithNonceRange(t, ctx, *k, 40, 45, chains.ZetaChainMainnet.ChainId, tss, zk)

	// set a cctx where the update function should fail to test that the next cctx are not updated but the next chains are
	failMap[sample.GetCctxIndexFromString("1-12")] = struct{}{}

	// test that the default crosschain flags are used when not set and the epoch length is not reached
	ctx = ctx.WithBlockHeight(observertypes.DefaultCrosschainFlags().GasPriceIncreaseFlags.EpochLength + 1)

	cctxCount, flags := k.IterateAndUpdateCctxGasPrice(ctx, supportedChains, updateFunc)
	require.Equal(t, 0, cctxCount)
	require.Equal(t, *observertypes.DefaultCrosschainFlags().GasPriceIncreaseFlags, flags)

	// test that custom crosschain flags are used when set and the epoch length is reached
	customFlags := observertypes.GasPriceIncreaseFlags{
		EpochLength:             100,
		RetryInterval:           time.Minute * 10,
		GasPriceIncreasePercent: 100,
		GasPriceIncreaseMax:     200,
		MaxPendingCctxs:         10,
	}
	crosschainFlags := sample.CrosschainFlags()
	crosschainFlags.GasPriceIncreaseFlags = &customFlags
	zk.ObserverKeeper.SetCrosschainFlags(ctx, *crosschainFlags)

	cctxCount, flags = k.IterateAndUpdateCctxGasPrice(ctx, supportedChains, updateFunc)
	require.Equal(t, 0, cctxCount)
	require.Equal(t, customFlags, flags)

	// test that cctx are iterated and updated when the epoch length is reached
	ctx = ctx.WithBlockHeight(observertypes.DefaultCrosschainFlags().GasPriceIncreaseFlags.EpochLength * 2)
	cctxCount, flags = k.IterateAndUpdateCctxGasPrice(ctx, supportedChains, updateFunc)

	// 2 eth + 5 bsc = 7
	require.Equal(t, 7, cctxCount)
	require.Equal(t, customFlags, flags)

	// check that the update function was called with the cctx index
	require.Equal(t, 7, len(updateFuncMap))
	require.Contains(t, updateFuncMap, sample.GetCctxIndexFromString("1-10"))
	require.Contains(t, updateFuncMap, sample.GetCctxIndexFromString("1-11"))

	require.Contains(t, updateFuncMap, sample.GetCctxIndexFromString("56-30"))
	require.Contains(t, updateFuncMap, sample.GetCctxIndexFromString("56-31"))
	require.Contains(t, updateFuncMap, sample.GetCctxIndexFromString("56-32"))
	require.Contains(t, updateFuncMap, sample.GetCctxIndexFromString("56-33"))
	require.Contains(t, updateFuncMap, sample.GetCctxIndexFromString("56-34"))
}

func TestCheckAndUpdateCctxGasPrice(t *testing.T) {
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
				OutboundParams: []*types.OutboundParams{
					{
						ReceiverChainId: 42,
						GasLimit:        1000,
						GasPrice:        "100",
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
				OutboundParams: []*types.OutboundParams{
					{
						ReceiverChainId: 42,
						GasLimit:        1000,
						GasPrice:        "100",
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
				OutboundParams: []*types.OutboundParams{
					{
						ReceiverChainId: 42,
						GasLimit:        1000,
						GasPrice:        "100",
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
				OutboundParams: []*types.OutboundParams{
					{
						ReceiverChainId: 42,
						GasLimit:        1000,
						GasPrice:        "100",
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
				OutboundParams: []*types.OutboundParams{
					{
						ReceiverChainId: 42,
						GasLimit:        100,
						GasPrice:        "",
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
				OutboundParams: []*types.OutboundParams{
					{
						ReceiverChainId: 42,
						GasLimit:        0,
						GasPrice:        "100",
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
				OutboundParams: []*types.OutboundParams{
					{
						ReceiverChainId: 42,
						GasLimit:        0,
						GasPrice:        "100",
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
				OutboundParams: []*types.OutboundParams{
					{
						ReceiverChainId: 42,
						GasLimit:        1000,
						GasPrice:        "100",
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
				OutboundParams: []*types.OutboundParams{
					{
						ReceiverChainId: 42,
						GasLimit:        1000,
						GasPrice:        "100",
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
			chainID := tc.cctx.GetCurrentOutboundParam().ReceiverChainId
			previousGasPrice, err := tc.cctx.GetCurrentOutboundParam().GetGasPriceUInt64()
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
				require.True(t, isFound)
				require.True(t, medianGasPrice.Equal(math.NewUint(tc.medianGasPrice)))
			}

			// set block timestamp
			ctx = ctx.WithBlockTime(tc.blockTimestamp)

			if tc.expectWithdrawFromGasStabilityPoolCall {
				fungibleMock.On(
					"WithdrawFromGasStabilityPool", ctx, chainID, tc.expectedAdditionalFees.BigInt(),
				).Return(tc.withdrawFromGasStabilityPoolReturn)
			}

			// check and update gas price
			gasPriceIncrease, feesPaid, err := keeper.CheckAndUpdateCctxGasPrice(ctx, *k, tc.cctx, tc.flags)

			if tc.isError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// check values
			require.True(t, gasPriceIncrease.Equal(tc.expectedGasPriceIncrease), "expected %s, got %s", tc.expectedGasPriceIncrease.String(), gasPriceIncrease.String())
			require.True(t, feesPaid.Equal(tc.expectedAdditionalFees), "expected %s, got %s", tc.expectedAdditionalFees.String(), feesPaid.String())

			// check cctx
			if !tc.expectedGasPriceIncrease.IsZero() {
				cctx, found := k.GetCrossChainTx(ctx, tc.cctx.Index)
				require.True(t, found)
				newGasPrice, err := cctx.GetCurrentOutboundParam().GetGasPriceUInt64()
				require.NoError(t, err)
				require.EqualValues(t, tc.expectedGasPriceIncrease.AddUint64(previousGasPrice).Uint64(), newGasPrice, "%d - %d", tc.expectedGasPriceIncrease.Uint64(), previousGasPrice)
				require.EqualValues(t, tc.blockTimestamp.Unix(), cctx.CctxStatus.LastUpdateTimestamp)
			}
		})
	}
}

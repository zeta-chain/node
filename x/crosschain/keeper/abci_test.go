package keeper_test

import (
	"errors"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	testkeeper "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
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
	supportedChains := []chains.Chain{
		{ChainId: chains.Ethereum.ChainId},
		{ChainId: chains.BitcoinMainnet.ChainId},
		{ChainId: chains.BscMainnet.ChainId},
		{ChainId: chains.ZetaChainMainnet.ChainId},
	}

	// set pending cctx
	tss := sample.Tss()
	zk.ObserverKeeper.SetTSS(ctx, tss)
	createCctxWithNonceRange(t, ctx, *k, 10, 15, chains.Ethereum.ChainId, tss, zk)
	createCctxWithNonceRange(t, ctx, *k, 20, 25, chains.BitcoinMainnet.ChainId, tss, zk)
	createCctxWithNonceRange(t, ctx, *k, 30, 35, chains.BscMainnet.ChainId, tss, zk)
	createCctxWithNonceRange(t, ctx, *k, 40, 45, chains.ZetaChainMainnet.ChainId, tss, zk)

	// set a cctx where the update function should fail to test that the next cctx are not updated but the next chains are
	failMap[sample.GetCctxIndexFromString("1-12")] = struct{}{}

	// test that the default crosschain flags are used when not set and the epoch length is not reached
	ctx = ctx.WithBlockHeight(observertypes.DefaultCrosschainFlags().GasPriceIncreaseFlags.EpochLength + 1)

	cctxCount, flags := k.IterateAndUpdateCCTXGasPrice(ctx, supportedChains, updateFunc)
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

	cctxCount, flags = k.IterateAndUpdateCCTXGasPrice(ctx, supportedChains, updateFunc)
	require.Equal(t, 0, cctxCount)
	require.Equal(t, customFlags, flags)

	// test that cctx are iterated and updated when the epoch length is reached
	ctx = ctx.WithBlockHeight(observertypes.DefaultCrosschainFlags().GasPriceIncreaseFlags.EpochLength * 2)
	cctxCount, flags = k.IterateAndUpdateCCTXGasPrice(ctx, supportedChains, updateFunc)

	// 2 eth + 5 btc + 5 bsc = 12
	require.Equal(t, 12, cctxCount)
	require.Equal(t, customFlags, flags)

	// check that the update function was called with the cctx index
	require.Equal(t, 12, len(updateFuncMap))
	require.Contains(t, updateFuncMap, sample.GetCctxIndexFromString("1-10"))
	require.Contains(t, updateFuncMap, sample.GetCctxIndexFromString("1-11"))

	require.Contains(t, updateFuncMap, sample.GetCctxIndexFromString("8332-20"))
	require.Contains(t, updateFuncMap, sample.GetCctxIndexFromString("8332-21"))
	require.Contains(t, updateFuncMap, sample.GetCctxIndexFromString("8332-22"))
	require.Contains(t, updateFuncMap, sample.GetCctxIndexFromString("8332-23"))
	require.Contains(t, updateFuncMap, sample.GetCctxIndexFromString("8332-24"))

	require.Contains(t, updateFuncMap, sample.GetCctxIndexFromString("56-30"))
	require.Contains(t, updateFuncMap, sample.GetCctxIndexFromString("56-31"))
	require.Contains(t, updateFuncMap, sample.GetCctxIndexFromString("56-32"))
	require.Contains(t, updateFuncMap, sample.GetCctxIndexFromString("56-33"))
	require.Contains(t, updateFuncMap, sample.GetCctxIndexFromString("56-34"))
}

func Test_CheckAndUpdateCCTXGasPrice(t *testing.T) {
	sampleTimestamp := time.Now()
	retryIntervalReached := sampleTimestamp.Add(observertypes.DefaultGasPriceIncreaseFlags.RetryInterval + time.Second)

	tt := []struct {
		name                                   string
		cctx                                   types.CrossChainTx
		flags                                  observertypes.GasPriceIncreaseFlags
		blockTimestamp                         time.Time
		medianGasPrice                         uint64
		medianPriorityFee                      uint64
		withdrawFromGasStabilityPoolReturn     error
		expectWithdrawFromGasStabilityPoolCall bool
		expectedGasPriceIncrease               math.Uint
		expectedAdditionalFees                 math.Uint
		isError                                bool
	}{
		{
			name:                                   "can update EVM chain gas price when retry interval is reached",
			cctx:                                   mkCustomCCTX(t, sampleTimestamp, 1, "100", 1000),
			flags:                                  observertypes.DefaultGasPriceIncreaseFlags,
			blockTimestamp:                         retryIntervalReached,
			medianGasPrice:                         50,
			medianPriorityFee:                      20,
			withdrawFromGasStabilityPoolReturn:     nil,
			expectWithdrawFromGasStabilityPoolCall: true,
			expectedGasPriceIncrease:               math.NewUint(50),    // 100% medianGasPrice
			expectedAdditionalFees:                 math.NewUint(50000), // gasLimit * increase
		},
		{
			name:                                   "can update BTC chain gas price when retry interval is reached",
			cctx:                                   mkCustomCCTX(t, sampleTimestamp, 8332, "10", 254),
			flags:                                  observertypes.DefaultGasPriceIncreaseFlags,
			blockTimestamp:                         sampleTimestamp.Add(time.Hour),
			medianGasPrice:                         20,
			medianPriorityFee:                      0,
			withdrawFromGasStabilityPoolReturn:     nil,
			expectWithdrawFromGasStabilityPoolCall: true,
			expectedGasPriceIncrease:               math.NewUint(10),   // medianGasPrice - gasPrice
			expectedAdditionalFees:                 math.NewUint(2540), // gasLimit * gasPriceIncrease
		},
		{
			name:                                   "skip if gas price is not set",
			cctx:                                   mkCustomCCTX(t, sampleTimestamp, 42, "", 100),
			flags:                                  observertypes.DefaultGasPriceIncreaseFlags,
			blockTimestamp:                         retryIntervalReached,
			medianGasPrice:                         100,
			expectWithdrawFromGasStabilityPoolCall: false,
			expectedGasPriceIncrease:               math.NewUint(0),
			expectedAdditionalFees:                 math.NewUint(0),
		},
		{
			name:                                   "skip if gas limit is not set",
			cctx:                                   mkCustomCCTX(t, sampleTimestamp, 42, "100", 0),
			flags:                                  observertypes.DefaultGasPriceIncreaseFlags,
			blockTimestamp:                         retryIntervalReached,
			medianGasPrice:                         100,
			expectWithdrawFromGasStabilityPoolCall: false,
			expectedGasPriceIncrease:               math.NewUint(0),
			expectedAdditionalFees:                 math.NewUint(0),
		},
		{
			name:                                   "returns error if can't find median gas price",
			cctx:                                   mkCustomCCTX(t, sampleTimestamp, 42, "100", 1000),
			flags:                                  observertypes.DefaultGasPriceIncreaseFlags,
			expectWithdrawFromGasStabilityPoolCall: false,
			blockTimestamp:                         retryIntervalReached,
			medianGasPrice:                         0,
			isError:                                true,
		},
		{
			name:                               "do nothing for non-EVM, non-BTC chain",
			cctx:                               mkCustomCCTX(t, sampleTimestamp, 100, "100", 1000),
			flags:                              observertypes.DefaultGasPriceIncreaseFlags,
			blockTimestamp:                     retryIntervalReached,
			medianGasPrice:                     50,
			medianPriorityFee:                  20,
			withdrawFromGasStabilityPoolReturn: nil,
			expectedGasPriceIncrease:           math.NewUint(0),
			expectedAdditionalFees:             math.NewUint(0),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// ARRANGE
			k, ctx := testkeeper.CrosschainKeeperAllMocks(t)
			fungibleMock := testkeeper.GetCrosschainFungibleMock(t, k)
			authorityMock := testkeeper.GetCrosschainAuthorityMock(t, k)
			chainID := tc.cctx.GetCurrentOutboundParam().ReceiverChainId
			previousGasPrice, err := tc.cctx.GetCurrentOutboundParam().GetGasPriceUInt64()
			if err != nil {
				previousGasPrice = 0
			}

			// set median gas price if not zero
			if tc.medianGasPrice != 0 {
				k.SetGasPrice(ctx, types.GasPrice{
					ChainId:      chainID,
					Prices:       []uint64{tc.medianGasPrice},
					PriorityFees: []uint64{tc.medianPriorityFee},
					MedianIndex:  0,
				})

				// ensure median gas price is set
				medianGasPrice, medianPriorityFee, isFound := k.GetMedianGasValues(ctx, chainID)
				require.True(t, isFound)
				require.True(t, medianGasPrice.Equal(math.NewUint(tc.medianGasPrice)))
				require.True(t, medianPriorityFee.Equal(math.NewUint(tc.medianPriorityFee)))
			}

			// set block timestamp
			ctx = ctx.WithBlockTime(tc.blockTimestamp)

			authorityMock.On("GetAdditionalChainList", ctx).Maybe().Return([]chains.Chain{})

			if tc.expectWithdrawFromGasStabilityPoolCall {
				fungibleMock.On(
					"WithdrawFromGasStabilityPool", ctx, chainID, tc.expectedAdditionalFees.BigInt(),
				).Return(tc.withdrawFromGasStabilityPoolReturn)
			}

			// ACT
			// check and update gas price
			gasPriceIncrease, feesPaid, err := keeper.CheckAndUpdateCCTXGasPrice(ctx, *k, tc.cctx, tc.flags)

			// ASSERT
			if tc.isError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// check values
			require.True(
				t,
				gasPriceIncrease.Equal(tc.expectedGasPriceIncrease),
				"expected %s, got %s",
				tc.expectedGasPriceIncrease.String(),
				gasPriceIncrease.String(),
			)
			require.True(
				t,
				feesPaid.Equal(tc.expectedAdditionalFees),
				"expected %s, got %s",
				tc.expectedAdditionalFees.String(),
				feesPaid.String(),
			)

			// check cctx
			if !tc.expectedGasPriceIncrease.IsZero() {
				cctx, found := k.GetCrossChainTx(ctx, tc.cctx.Index)
				require.True(t, found)
				newGasPrice, err := cctx.GetCurrentOutboundParam().GetGasPriceUInt64()
				require.NoError(t, err)
				require.EqualValues(
					t,
					tc.expectedGasPriceIncrease.AddUint64(previousGasPrice).Uint64(),
					newGasPrice,
					"%d - %d",
					tc.expectedGasPriceIncrease.Uint64(),
					previousGasPrice,
				)
				require.EqualValues(t, tc.blockTimestamp.Unix(), cctx.CctxStatus.LastUpdateTimestamp)
			}
		})
	}
}

func Test_CheckAndUpdateCCTXGasPriceEVM(t *testing.T) {
	sampleTimestamp := time.Now()
	chainID := chains.Ethereum.ChainId
	retryIntervalReached := sampleTimestamp.Add(observertypes.DefaultGasPriceIncreaseFlags.RetryInterval + time.Second)
	retryIntervalNotReached := sampleTimestamp.Add(
		observertypes.DefaultGasPriceIncreaseFlags.RetryInterval - time.Second,
	)

	tt := []struct {
		name                                   string
		cctx                                   types.CrossChainTx
		flags                                  observertypes.GasPriceIncreaseFlags
		blockTimestamp                         time.Time
		medianGasPrice                         uint64
		medianPriorityFee                      uint64
		withdrawFromGasStabilityPoolReturn     error
		expectWithdrawFromGasStabilityPoolCall bool
		expectedGasPriceIncrease               math.Uint
		expectedAdditionalFees                 math.Uint
		isError                                bool
	}{
		{
			name:                                   "can update gas price",
			cctx:                                   mkCustomCCTX(t, sampleTimestamp, chainID, "100", 1000),
			flags:                                  observertypes.DefaultGasPriceIncreaseFlags,
			blockTimestamp:                         retryIntervalReached,
			medianGasPrice:                         50,
			medianPriorityFee:                      20,
			withdrawFromGasStabilityPoolReturn:     nil,
			expectWithdrawFromGasStabilityPoolCall: true,
			expectedGasPriceIncrease:               math.NewUint(50),    // 100% medianGasPrice
			expectedAdditionalFees:                 math.NewUint(50000), // gasLimit * increase
		},
		{
			name: "can update gas price at max limit",
			cctx: mkCustomCCTX(t, sampleTimestamp, chainID, "100", 1000),
			flags: observertypes.GasPriceIncreaseFlags{
				EpochLength:             100,
				RetryInterval:           time.Minute * 10,
				GasPriceIncreasePercent: 200, // Increase gas price to 100+50*2 = 200
				GasPriceIncreaseMax:     400, // Max gas price is 50*4 = 200
			},
			blockTimestamp:                         retryIntervalReached,
			medianGasPrice:                         50,
			medianPriorityFee:                      20,
			withdrawFromGasStabilityPoolReturn:     nil,
			expectWithdrawFromGasStabilityPoolCall: true,
			expectedGasPriceIncrease:               math.NewUint(100),    // 200% medianGasPrice
			expectedAdditionalFees:                 math.NewUint(100000), // gasLimit * increase
		},
		{
			name: "default gas price increase limit used if not defined",
			cctx: mkCustomCCTX(t, sampleTimestamp, chainID, "100", 1000),
			flags: observertypes.GasPriceIncreaseFlags{
				EpochLength:             100,
				RetryInterval:           time.Minute * 10,
				GasPriceIncreasePercent: 100,
				GasPriceIncreaseMax:     0, // Limit should not be reached
			},
			blockTimestamp:                         retryIntervalReached,
			medianGasPrice:                         50,
			medianPriorityFee:                      20,
			withdrawFromGasStabilityPoolReturn:     nil,
			expectWithdrawFromGasStabilityPoolCall: true,
			expectedGasPriceIncrease:               math.NewUint(50),    // 100% medianGasPrice
			expectedAdditionalFees:                 math.NewUint(50000), // gasLimit * increase
		},
		{
			name:                                   "skip if retry interval is not reached",
			cctx:                                   mkCustomCCTX(t, sampleTimestamp, chainID, "100", 1000),
			flags:                                  observertypes.DefaultGasPriceIncreaseFlags,
			blockTimestamp:                         retryIntervalNotReached,
			medianGasPrice:                         50,
			medianPriorityFee:                      20,
			expectWithdrawFromGasStabilityPoolCall: false,
			expectedGasPriceIncrease:               math.NewUint(0),
			expectedAdditionalFees:                 math.NewUint(0),
		},
		{
			name: "returns error if can't get CCTX gas price",
			cctx: mkCustomCCTX(
				t,
				sampleTimestamp,
				chainID,
				"invalid_gas_price",
				1000,
			),
			flags:                                  observertypes.DefaultGasPriceIncreaseFlags,
			blockTimestamp:                         retryIntervalReached,
			medianGasPrice:                         50,
			medianPriorityFee:                      20,
			expectWithdrawFromGasStabilityPoolCall: false,
			expectedGasPriceIncrease:               math.NewUint(0),
			expectedAdditionalFees:                 math.NewUint(0),
			isError:                                true,
		},
		{
			name: "skip if max limit reached",
			cctx: mkCustomCCTX(t, sampleTimestamp, chainID, "100", 1000),
			flags: observertypes.GasPriceIncreaseFlags{
				EpochLength:             100,
				RetryInterval:           time.Minute * 10,
				GasPriceIncreasePercent: 200, // Increase gas price to 100+50*2 = 200
				GasPriceIncreaseMax:     300, // Max gas price is 50*3 = 150
			},
			blockTimestamp:                         retryIntervalReached,
			medianGasPrice:                         50,
			medianPriorityFee:                      20,
			expectWithdrawFromGasStabilityPoolCall: false,
			expectedGasPriceIncrease:               math.NewUint(0),
			expectedAdditionalFees:                 math.NewUint(0),
		},
		{
			name:                                   "returns error if can't withdraw from gas stability pool",
			cctx:                                   mkCustomCCTX(t, sampleTimestamp, chainID, "100", 1000),
			flags:                                  observertypes.DefaultGasPriceIncreaseFlags,
			blockTimestamp:                         retryIntervalReached,
			medianGasPrice:                         50,
			medianPriorityFee:                      20,
			expectWithdrawFromGasStabilityPoolCall: true,
			expectedGasPriceIncrease:               math.NewUint(0),     // expect 0 on error
			expectedAdditionalFees:                 math.NewUint(50000), // gasLimit * increase
			withdrawFromGasStabilityPoolReturn:     errors.New("withdraw error"),
			isError:                                true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// ARRANGE
			k, ctx := testkeeper.CrosschainKeeperAllMocks(t)
			fungibleMock := testkeeper.GetCrosschainFungibleMock(t, k)
			chainID := tc.cctx.GetCurrentOutboundParam().ReceiverChainId
			previousGasPrice, err := tc.cctx.GetCurrentOutboundParam().GetGasPriceUInt64()
			if err != nil {
				previousGasPrice = 0
			}

			// set block timestamp
			ctx = ctx.WithBlockTime(tc.blockTimestamp)

			if tc.expectWithdrawFromGasStabilityPoolCall {
				fungibleMock.On(
					"WithdrawFromGasStabilityPool", ctx, chainID, tc.expectedAdditionalFees.BigInt(),
				).Return(tc.withdrawFromGasStabilityPoolReturn)
			}

			// ACT
			// check and update gas price
			gasPriceIncrease, feesPaid, err := keeper.CheckAndUpdateCCTXGasPriceEVM(
				ctx,
				*k,
				math.NewUint(tc.medianGasPrice),
				math.NewUint(tc.medianPriorityFee),
				tc.cctx,
				tc.flags,
			)

			// ASSERT
			if tc.isError {
				require.Error(t, err)
				require.True(t, gasPriceIncrease.IsZero())
				require.True(t, feesPaid.IsZero())
				return
			}
			require.NoError(t, err)

			// check values
			require.True(
				t,
				gasPriceIncrease.Equal(tc.expectedGasPriceIncrease),
				"expected %s, got %s",
				tc.expectedGasPriceIncrease.String(),
				gasPriceIncrease.String(),
			)
			require.True(
				t,
				feesPaid.Equal(tc.expectedAdditionalFees),
				"expected %s, got %s",
				tc.expectedAdditionalFees.String(),
				feesPaid.String(),
			)

			// check cctx
			if !tc.expectedGasPriceIncrease.IsZero() {
				cctx, found := k.GetCrossChainTx(ctx, tc.cctx.Index)
				require.True(t, found)
				newGasPrice, err := cctx.GetCurrentOutboundParam().GetGasPriceUInt64()
				require.NoError(t, err)
				require.EqualValues(
					t,
					tc.expectedGasPriceIncrease.AddUint64(previousGasPrice).Uint64(),
					newGasPrice,
					"%d - %d",
					tc.expectedGasPriceIncrease.Uint64(),
					previousGasPrice,
				)
				require.EqualValues(t, tc.blockTimestamp.Unix(), cctx.CctxStatus.LastUpdateTimestamp)
			}
		})
	}
}

func Test_CheckAndUpdateCCTXGasPriceBTC(t *testing.T) {
	sampleTimestamp := time.Now()
	chainID := chains.BitcoinMainnet.ChainId
	flags := observertypes.DefaultGasPriceIncreaseFlags
	gasRateUpdateInterval := flags.RetryIntervalBTC
	retryIntervalReached := sampleTimestamp.Add(gasRateUpdateInterval + time.Second)
	retryIntervalNotReached := sampleTimestamp.Add(gasRateUpdateInterval - time.Second)

	tt := []struct {
		name                                   string
		cctx                                   types.CrossChainTx
		flags                                  observertypes.GasPriceIncreaseFlags
		blockTimestamp                         time.Time
		medianGasPrice                         uint64
		withdrawFromGasStabilityPoolReturn     error
		expectWithdrawFromGasStabilityPoolCall bool
		expectGasPriorityFeeUpdate             bool
		expectedGasPriceIncrease               math.Uint
		expectedAdditionalFees                 math.Uint
		isError                                bool
	}{
		{
			name:                                   "can update fee rate",
			cctx:                                   mkCustomCCTX(t, sampleTimestamp, chainID, "10", 200),
			flags:                                  flags,
			blockTimestamp:                         retryIntervalReached,
			medianGasPrice:                         12,
			withdrawFromGasStabilityPoolReturn:     nil,
			expectWithdrawFromGasStabilityPoolCall: true,
			expectGasPriorityFeeUpdate:             true,
			expectedGasPriceIncrease:               math.NewUint(2),   // medianGasPrice - gasPrice
			expectedAdditionalFees:                 math.NewUint(400), // gasLimit * increase
		},
		{
			name:                                   "skip if retry interval is not reached",
			cctx:                                   mkCustomCCTX(t, sampleTimestamp, chainID, "10", 200),
			flags:                                  flags,
			blockTimestamp:                         retryIntervalNotReached,
			medianGasPrice:                         12,
			expectWithdrawFromGasStabilityPoolCall: false,
			expectedGasPriceIncrease:               math.NewUint(0),
			expectedAdditionalFees:                 math.NewUint(0),
		},
		{
			name:                                   "default interval used if not defined",
			cctx:                                   mkCustomCCTX(t, sampleTimestamp, chainID, "10", 200),
			flags:                                  observertypes.GasPriceIncreaseFlags{}, // interval not set
			blockTimestamp:                         retryIntervalNotReached,
			medianGasPrice:                         12,
			expectWithdrawFromGasStabilityPoolCall: false,
			expectedGasPriceIncrease:               math.NewUint(0),
			expectedAdditionalFees:                 math.NewUint(0),
		},
		{
			name:                                   "returns error if can't get CCTX gas price",
			cctx:                                   mkCustomCCTX(t, sampleTimestamp, chainID, "invalid_gas_price", 200),
			flags:                                  flags,
			blockTimestamp:                         retryIntervalReached,
			medianGasPrice:                         12,
			expectWithdrawFromGasStabilityPoolCall: false,
			expectGasPriorityFeeUpdate:             true,
			expectedGasPriceIncrease:               math.NewUint(0),
			expectedAdditionalFees:                 math.NewUint(0),
			isError:                                true,
		},
		{
			name:                                   "skip if current gas price is equal or higher than median gas price",
			cctx:                                   mkCustomCCTX(t, sampleTimestamp, chainID, "10", 200),
			flags:                                  flags,
			blockTimestamp:                         retryIntervalReached,
			medianGasPrice:                         10,
			expectWithdrawFromGasStabilityPoolCall: false,
			expectGasPriorityFeeUpdate:             true,
			expectedGasPriceIncrease:               math.NewUint(0),
			expectedAdditionalFees:                 math.NewUint(0),
		},
		{
			name:                                   "returns error if can't withdraw from gas stability pool",
			cctx:                                   mkCustomCCTX(t, sampleTimestamp, chainID, "10", 200),
			flags:                                  flags,
			blockTimestamp:                         retryIntervalReached,
			medianGasPrice:                         12,
			expectWithdrawFromGasStabilityPoolCall: true,
			expectGasPriorityFeeUpdate:             true,
			expectedGasPriceIncrease:               math.NewUint(0),   // expect 0 on error
			expectedAdditionalFees:                 math.NewUint(400), // gasLimit * increase
			withdrawFromGasStabilityPoolReturn:     errors.New("withdraw error"),
			isError:                                true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// ARRANGE
			k, ctx := testkeeper.CrosschainKeeperAllMocks(t)
			fungibleMock := testkeeper.GetCrosschainFungibleMock(t, k)
			chainID := tc.cctx.GetCurrentOutboundParam().ReceiverChainId
			previousGasPrice, err := tc.cctx.GetCurrentOutboundParam().GetGasPriceUInt64()
			if err != nil {
				previousGasPrice = 0
			}

			// set block timestamp
			ctx = ctx.WithBlockTime(tc.blockTimestamp)

			// mock up gas stability pool withdraw
			if tc.expectWithdrawFromGasStabilityPoolCall {
				fungibleMock.On(
					"WithdrawFromGasStabilityPool", ctx, chainID, tc.expectedAdditionalFees.BigInt(),
				).Return(tc.withdrawFromGasStabilityPoolReturn)
			}

			// ACT
			// check and update gas rate
			gasPriceIncrease, feesPaid, err := keeper.CheckAndUpdateCCTXGasPriceBTC(
				ctx,
				*k,
				math.NewUint(tc.medianGasPrice),
				tc.cctx,
				tc.flags,
			)

			// ASSERT
			if tc.isError {
				require.Error(t, err)
				require.True(t, gasPriceIncrease.IsZero())
				require.True(t, feesPaid.IsZero())
				return
			}
			require.NoError(t, err)

			// check values
			require.True(
				t,
				gasPriceIncrease.Equal(tc.expectedGasPriceIncrease),
				"expected %s, got %s",
				tc.expectedGasPriceIncrease.String(),
				gasPriceIncrease.String(),
			)
			require.True(
				t,
				feesPaid.Equal(tc.expectedAdditionalFees),
				"expected %s, got %s",
				tc.expectedAdditionalFees.String(),
				feesPaid.String(),
			)

			// check gas priority fee update and last update timestamp
			if tc.expectGasPriorityFeeUpdate {
				cctx, found := k.GetCrossChainTx(ctx, tc.cctx.Index)
				require.True(t, found)
				newPriorityFee, err := cctx.GetCurrentOutboundParam().GetGasPriorityFeeUInt64()
				require.NoError(t, err)
				require.Equal(t, tc.medianGasPrice, newPriorityFee)
				require.EqualValues(t, tc.blockTimestamp.Unix(), cctx.CctxStatus.LastUpdateTimestamp)
			}

			// check gas price update and last update timestamp
			if !tc.expectedGasPriceIncrease.IsZero() {
				cctx, found := k.GetCrossChainTx(ctx, tc.cctx.Index)
				require.True(t, found)
				newGasPrice, err := cctx.GetCurrentOutboundParam().GetGasPriceUInt64()
				require.NoError(t, err)
				require.EqualValues(
					t,
					tc.expectedGasPriceIncrease.AddUint64(previousGasPrice).Uint64(),
					newGasPrice,
					"%d - %d",
					tc.expectedGasPriceIncrease.Uint64(),
					previousGasPrice,
				)
				require.EqualValues(t, tc.blockTimestamp.Unix(), cctx.CctxStatus.LastUpdateTimestamp)
			}
		})
	}
}

func Test_IsCCTXGasPriceUpdateSupported(t *testing.T) {
	tt := []struct {
		name      string
		chainID   int64
		isSupport bool
	}{
		{
			name:      "Zetachain is unsupported for gas price update",
			chainID:   chains.ZetaChainMainnet.ChainId,
			isSupport: false,
		},
		{
			name:      "Zetachain testnet is unsupported for gas price update",
			chainID:   chains.ZetaChainTestnet.ChainId,
			isSupport: false,
		},
		{
			name:      "Ethereum is supported for gas price update",
			chainID:   chains.Ethereum.ChainId,
			isSupport: true,
		},
		{
			name:      "Ethereum Sepolia is supported for gas price update",
			chainID:   chains.Sepolia.ChainId,
			isSupport: true,
		},
		{
			name:      "BSC is supported for gas price update",
			chainID:   chains.BscMainnet.ChainId,
			isSupport: true,
		},
		{
			name:      "BSC testnet is supported for gas price update",
			chainID:   chains.BscTestnet.ChainId,
			isSupport: true,
		},
		{
			name:      "Polygon is supported for gas price update",
			chainID:   chains.Polygon.ChainId,
			isSupport: true,
		},
		{
			name:      "Polygon Amoy is supported for gas price update",
			chainID:   chains.Amoy.ChainId,
			isSupport: true,
		},
		{
			name:      "Base is supported for gas price update",
			chainID:   chains.BaseMainnet.ChainId,
			isSupport: true,
		},
		{
			name:      "Base Sepolia is supported for gas price update",
			chainID:   chains.BaseSepolia.ChainId,
			isSupport: true,
		},
		{
			name:      "Bitcoin is supported for gas price update",
			chainID:   chains.BitcoinMainnet.ChainId,
			isSupport: true,
		},
		{
			name:      "Bitcoin testnet is supported for gas price update",
			chainID:   chains.BitcoinTestnet4.ChainId,
			isSupport: true,
		},
		{
			name:      "Solana is unsupported for gas price update",
			chainID:   chains.SolanaMainnet.ChainId,
			isSupport: false,
		},
		{
			name:      "Solana devnet is unsupported for gas price update",
			chainID:   chains.SolanaDevnet.ChainId,
			isSupport: false,
		},
		{
			name:      "TON is unsupported for gas price update",
			chainID:   chains.TONMainnet.ChainId,
			isSupport: false,
		},
		{
			name:      "TON testnet is unsupported for gas price update",
			chainID:   chains.TONTestnet.ChainId,
			isSupport: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			isSupported := keeper.IsCCTXGasPriceUpdateSupported(tc.chainID, []chains.Chain{})
			require.Equal(t, tc.isSupport, isSupported)
		})
	}
}

// mkCustomCCTX is a helper function to create a CCTX with given time, chainID, gasPrice, gasLimit
func mkCustomCCTX(t *testing.T, time time.Time, chainID int64, gasPrice string, gasLimit uint64) types.CrossChainTx {
	return types.CrossChainTx{
		Index: sample.ZetaIndex(t),
		CctxStatus: &types.Status{
			CreatedTimestamp:    time.Unix(),
			LastUpdateTimestamp: time.Unix(),
		},
		OutboundParams: []*types.OutboundParams{
			{
				ReceiverChainId: chainID,
				CallOptions: &types.CallOptions{
					GasLimit: gasLimit,
				},
				GasPrice: gasPrice,
			},
		},
	}
}

package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

var (
	// local eth chain ID
	ethChainID = getValidEthChainID()

	// local btc chain ID
	btcChainID = getValidBtcChainID()
)

// createTestRateLimiterFlags creates a custom rate limiter flags
func createTestRateLimiterFlags(
	window int64,
	rate math.Uint,
	zrc20ETH string,
	zrc20BTC string,
	zrc20USDT string,
	ethRate string,
	btcRate string,
	usdtRate string,
) *types.RateLimiterFlags {
	return &types.RateLimiterFlags{
		Enabled: true,
		Window:  window, // for instance: 500 zeta blocks, 50 minutes
		Rate:    rate,
		Conversions: []types.Conversion{
			// ETH
			{
				Zrc20: zrc20ETH,
				Rate:  sdk.MustNewDecFromStr(ethRate),
			},
			// BTC
			{
				Zrc20: zrc20BTC,
				Rate:  sdk.MustNewDecFromStr(btcRate),
			},
			// USDT
			{
				Zrc20: zrc20USDT,
				Rate:  sdk.MustNewDecFromStr(usdtRate),
			},
		},
	}
}

// setCctxsInKeeper sets the given cctxs to the keeper
func setCctxsInKeeper(
	ctx sdk.Context,
	k keeper.Keeper,
	zk keepertest.ZetaKeepers,
	tss observertypes.TSS,
	cctxs []*types.CrossChainTx,
) {
	for _, cctx := range cctxs {
		k.SetCrossChainTx(ctx, *cctx)
		zk.ObserverKeeper.SetNonceToCctx(ctx, observertypes.NonceToCctx{
			ChainId: cctx.GetCurrentOutboundParam().ReceiverChainId,
			// #nosec G115 always in range for tests
			Nonce:     int64(cctx.GetCurrentOutboundParam().TssNonce),
			CctxIndex: cctx.Index,
			Tss:       tss.TssPubkey,
		})
	}
}

// setupForeignCoins adds ETH, BTC, USDT to the foreign coins store
func setupForeignCoins(
	t *testing.T,
	ctx sdk.Context,
	zk keepertest.ZetaKeepers,
	zrc20ETH, zrc20BTC, zrc20USDT, assetUSDT string,
) {
	// set foreign coins
	fCoins := sample.ForeignCoinList(t, zrc20ETH, zrc20BTC, zrc20USDT, assetUSDT)
	for _, fc := range fCoins {
		zk.FungibleKeeper.SetForeignCoins(ctx, fc)
	}
}

func TestKeeper_RateLimiterInput(t *testing.T) {
	// create sample TSS
	tss := sample.Tss()
	zetaChainID := chains.ZetaChainMainnet.ChainId

	// create sample zrc20 addresses for ETH, BTC, USDT
	zrc20ETH := sample.EthAddress().Hex()
	zrc20BTC := sample.EthAddress().Hex()
	zrc20USDT := sample.EthAddress().Hex()

	// create Eth chain 999 mined and 200 pending cctxs for rate limiter test
	// the number 999 is to make it less than `MaxLookbackNonce` so the LoopBackwards gets the chance to hit nonce 0
	ethMinedCctxs := sample.CustomCctxsInBlockRange(
		t,
		1,
		999,
		zetaChainID,
		ethChainID,
		coin.CoinType_Gas,
		"",
		uint64(1e15),
		types.CctxStatus_OutboundMined,
	)
	ethPendingCctxs := sample.CustomCctxsInBlockRange(
		t,
		1000,
		1199,
		zetaChainID,
		ethChainID,
		coin.CoinType_Gas,
		"",
		uint64(1e15),
		types.CctxStatus_PendingOutbound,
	)

	// create Eth chain 999 reverted and 200 pending revert cctxs for rate limiter test
	// the number 999 is to make it less than `MaxLookbackNonce` so the LoopBackwards gets the chance to hit nonce 0
	// these cctxs should be ignored by the rate limiter as it can't compare `ObservedExternalHeight` against the window boundary
	ethRevertedCctxs := sample.CustomCctxsInBlockRange(
		t,
		1,
		999,
		ethChainID,
		ethChainID,
		coin.CoinType_Gas,
		"",
		uint64(1e15),
		types.CctxStatus_Reverted,
	)
	ethPendingRevertCctxs := sample.CustomCctxsInBlockRange(
		t,
		1000,
		1199,
		ethChainID,
		ethChainID,
		coin.CoinType_Gas,
		"",
		uint64(1e15),
		types.CctxStatus_PendingRevert,
	)

	// create Btc chain 999 mined and 200 pending cctxs for rate limiter test
	// the number 999 is to make it less than `MaxLookbackNonce` so the LoopBackwards gets the chance to hit nonce 0
	btcMinedCctxs := sample.CustomCctxsInBlockRange(
		t,
		1,
		999,
		zetaChainID,
		btcChainID,
		coin.CoinType_Gas,
		"",
		1000,
		types.CctxStatus_OutboundMined,
	)
	btcPendingCctxs := sample.CustomCctxsInBlockRange(
		t,
		1000,
		1199,
		zetaChainID,
		btcChainID,
		coin.CoinType_Gas,
		"",
		1000,
		types.CctxStatus_PendingOutbound,
	)

	// define test cases
	tests := []struct {
		name           string
		rateLimitFlags *types.RateLimiterFlags

		// Eth chain cctxs setup
		ethMinedCctxs    []*types.CrossChainTx
		ethPendingCctxs  []*types.CrossChainTx
		ethPendingNonces observertypes.PendingNonces

		// Btc chain cctxs setup
		btcMinedCctxs    []*types.CrossChainTx
		btcPendingCctxs  []*types.CrossChainTx
		btcPendingNonces observertypes.PendingNonces

		// block height and limit of cctxs to retrieve
		currentHeight int64
		queryLimit    uint32

		// expected results
		expectedHeight                  int64
		expectedCctxsMissed             []*types.CrossChainTx
		expectedCctxsPending            []*types.CrossChainTx
		expectedTotalPending            uint64
		expectedPastCctxsValue          string
		expectedPendingCctxsValue       string
		expectedLowestPendingCctxHeight int64
	}{
		{
			name: "can retrieve all pending cctxs",
			rateLimitFlags: createTestRateLimiterFlags(
				500,
				math.NewUint(10*1e18),
				zrc20ETH,
				zrc20BTC,
				zrc20USDT,
				"2500",
				"50000",
				"0.8",
			),
			ethMinedCctxs:   ethMinedCctxs,
			ethPendingCctxs: ethPendingCctxs,
			ethPendingNonces: observertypes.PendingNonces{
				ChainId:   ethChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			btcMinedCctxs:   btcMinedCctxs,
			btcPendingCctxs: btcPendingCctxs,
			btcPendingNonces: observertypes.PendingNonces{
				ChainId:   btcChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			currentHeight: 1199,
			queryLimit:    0, // use default MaxPendingCctxs

			// expected results
			expectedHeight: 1199,
			expectedCctxsMissed: keeper.SortCctxsByHeightAndChainID(
				append(append([]*types.CrossChainTx{}, ethPendingCctxs[0:100]...), btcPendingCctxs[0:100]...),
			),
			expectedCctxsPending: keeper.SortCctxsByHeightAndChainID(
				append(append([]*types.CrossChainTx{}, ethPendingCctxs[100:200]...), btcPendingCctxs[100:200]...),
			),
			expectedTotalPending:            400,
			expectedPastCctxsValue:          sdk.NewInt(1200).Mul(sdk.NewInt(1e18)).String(), // 400 * (2.5 + 0.5) ZETA
			expectedPendingCctxsValue:       sdk.NewInt(300).Mul(sdk.NewInt(1e18)).String(),  // 100 * (2.5 + 0.5) ZETA
			expectedLowestPendingCctxHeight: 1100,
		},
		{
			name: "scan retrieve all pending cctxs and ignore revert cctxs",
			rateLimitFlags: createTestRateLimiterFlags(
				500,
				math.NewUint(10*1e18),
				zrc20ETH,
				zrc20BTC,
				zrc20USDT,
				"2500",
				"50000",
				"0.8",
			),
			ethMinedCctxs:   ethRevertedCctxs,
			ethPendingCctxs: ethPendingRevertCctxs,
			ethPendingNonces: observertypes.PendingNonces{
				ChainId:   ethChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			btcMinedCctxs:   btcMinedCctxs,
			btcPendingCctxs: btcPendingCctxs,
			btcPendingNonces: observertypes.PendingNonces{
				ChainId:   btcChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			currentHeight: 1199,
			queryLimit:    0, // use default MaxPendingCctxs

			// expected results
			expectedHeight: 1199,
			expectedCctxsMissed: keeper.SortCctxsByHeightAndChainID(
				append(append([]*types.CrossChainTx{}, ethPendingRevertCctxs[0:100]...), btcPendingCctxs[0:100]...),
			),
			expectedCctxsPending: keeper.SortCctxsByHeightAndChainID(
				append(append([]*types.CrossChainTx{}, ethPendingRevertCctxs[100:200]...), btcPendingCctxs[100:200]...),
			),
			expectedTotalPending: 400,
			expectedPastCctxsValue: sdk.NewInt(200).
				Mul(sdk.NewInt(1e18)).
				String(),
			// 400 * 0.5 ZETA, ignore Eth chain reverted cctxs
			expectedPendingCctxsValue:       sdk.NewInt(300).Mul(sdk.NewInt(1e18)).String(), // 100 * (2.5 + 0.5) ZETA
			expectedLowestPendingCctxHeight: 1100,
		},
		{
			name: "should use left window boundary 1 if window > currentHeight",
			rateLimitFlags: createTestRateLimiterFlags(
				1200,
				math.NewUint(10*1e18),
				zrc20ETH,
				zrc20BTC,
				zrc20USDT,
				"2500",
				"50000",
				"0.8",
			),
			ethMinedCctxs:   ethMinedCctxs,
			ethPendingCctxs: ethPendingCctxs,
			ethPendingNonces: observertypes.PendingNonces{
				ChainId:   ethChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			btcMinedCctxs:   btcMinedCctxs,
			btcPendingCctxs: btcPendingCctxs,
			btcPendingNonces: observertypes.PendingNonces{
				ChainId:   btcChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			currentHeight: 1199,                       // window 1200 > 1199
			queryLimit:    keeper.MaxPendingCctxs + 1, // should use default MaxPendingCctxs

			// expected results
			expectedHeight: 1199,
			expectedCctxsMissed: keeper.SortCctxsByHeightAndChainID(
				append(append([]*types.CrossChainTx{}, ethPendingCctxs[0:100]...), btcPendingCctxs[0:100]...),
			),
			expectedCctxsPending: keeper.SortCctxsByHeightAndChainID(
				append(append([]*types.CrossChainTx{}, ethPendingCctxs[100:200]...), btcPendingCctxs[100:200]...),
			),
			expectedTotalPending:            400,
			expectedPastCctxsValue:          sdk.NewInt(3297).Mul(sdk.NewInt(1e18)).String(), // 1099 * (2.5 + 0.5) ZETA
			expectedPendingCctxsValue:       sdk.NewInt(300).Mul(sdk.NewInt(1e18)).String(),  // 100 * (2.5 + 0.5) ZETA
			expectedLowestPendingCctxHeight: 1100,
		},
		{
			name: "should loop from nonce 0",
			rateLimitFlags: createTestRateLimiterFlags(
				500,
				math.NewUint(10*1e18),
				zrc20ETH,
				zrc20BTC,
				zrc20USDT,
				"2500",
				"50000",
				"0.8",
			),
			ethMinedCctxs:   ethMinedCctxs,
			ethPendingCctxs: ethPendingCctxs,
			ethPendingNonces: observertypes.PendingNonces{
				ChainId:   ethChainID,
				NonceLow:  999, // startNonce will be set to 0 (NonceLow - 1000 < 0)
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			btcMinedCctxs:   btcMinedCctxs,
			btcPendingCctxs: btcPendingCctxs,
			btcPendingNonces: observertypes.PendingNonces{
				ChainId:   btcChainID,
				NonceLow:  999, // startNonce will be set to 0 (NonceLow - 1000 < 0)
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			currentHeight: 1199,
			queryLimit:    keeper.MaxPendingCctxs,

			// expected results
			expectedHeight:      1199,
			expectedCctxsMissed: []*types.CrossChainTx{}, // no missed cctxs
			expectedCctxsPending: keeper.SortCctxsByHeightAndChainID(
				append(append([]*types.CrossChainTx{}, ethPendingCctxs...), btcPendingCctxs...),
			),
			expectedTotalPending:            400,
			expectedPastCctxsValue:          sdk.NewInt(900).Mul(sdk.NewInt(1e18)).String(), // 300 * (2.5 + 0.5) ZETA
			expectedPendingCctxsValue:       sdk.NewInt(600).Mul(sdk.NewInt(1e18)).String(), // 200 * (2.5 + 0.5) ZETA
			expectedLowestPendingCctxHeight: 1000,
		},
		{
			name: "set a lower gRPC request limit < len(pending_cctxs)",
			rateLimitFlags: createTestRateLimiterFlags(
				500,
				math.NewUint(10*1e18),
				zrc20ETH,
				zrc20BTC,
				zrc20USDT,
				"2500",
				"50000",
				"0.8",
			),
			ethMinedCctxs:   ethMinedCctxs,
			ethPendingCctxs: ethPendingCctxs,
			ethPendingNonces: observertypes.PendingNonces{
				ChainId:   ethChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			btcMinedCctxs:   btcMinedCctxs,
			btcPendingCctxs: btcPendingCctxs,
			btcPendingNonces: observertypes.PendingNonces{
				ChainId:   btcChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			currentHeight: 1199,
			queryLimit:    100, // 100 < keeper.MaxPendingCctxs

			// expected results
			expectedHeight: 1199,
			// should include all missed and 50 pending cctxs for each chain
			expectedCctxsMissed: keeper.SortCctxsByHeightAndChainID(
				append(append([]*types.CrossChainTx{}, ethPendingCctxs[0:100]...), btcPendingCctxs[0:100]...),
			),
			expectedCctxsPending: keeper.SortCctxsByHeightAndChainID(
				append(append([]*types.CrossChainTx{}, ethPendingCctxs[100:150]...), btcPendingCctxs[100:150]...),
			),
			expectedTotalPending:            400,
			expectedPastCctxsValue:          sdk.NewInt(1200).Mul(sdk.NewInt(1e18)).String(), // 400 * (2.5 + 0.5) ZETA
			expectedPendingCctxsValue:       sdk.NewInt(300).Mul(sdk.NewInt(1e18)).String(),  // 100 * (2.5 + 0.5) ZETA
			expectedLowestPendingCctxHeight: 1100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create test keepers
			k, ctx, _, zk := keepertest.CrosschainKeeper(t)

			// Set TSS
			zk.ObserverKeeper.SetTSS(ctx, tss)

			// Set foreign coins
			assetUSDT := sample.EthAddress().Hex()
			setupForeignCoins(t, ctx, zk, zrc20ETH, zrc20BTC, zrc20USDT, assetUSDT)

			// Set rate limiter flags
			if tt.rateLimitFlags != nil {
				k.SetRateLimiterFlags(ctx, *tt.rateLimitFlags)
			}

			// Set Eth chain mined cctxs, pending ccxts and pending nonces
			setCctxsInKeeper(ctx, *k, zk, tss, tt.ethMinedCctxs)
			setCctxsInKeeper(ctx, *k, zk, tss, tt.ethPendingCctxs)
			zk.ObserverKeeper.SetPendingNonces(ctx, tt.ethPendingNonces)

			// Set Btc chain mined cctxs, pending ccxts and pending nonces
			setCctxsInKeeper(ctx, *k, zk, tss, tt.btcMinedCctxs)
			setCctxsInKeeper(ctx, *k, zk, tss, tt.btcPendingCctxs)
			zk.ObserverKeeper.SetPendingNonces(ctx, tt.btcPendingNonces)

			// Set current block height
			ctx = ctx.WithBlockHeight(tt.currentHeight)

			// Query input data for the rate limiter
			request := &types.QueryRateLimiterInputRequest{
				Limit:  tt.queryLimit,
				Window: tt.rateLimitFlags.Window,
			}
			res, err := k.RateLimiterInput(ctx, request)

			// check results
			require.NoError(t, err)
			require.Equal(t, tt.expectedHeight, res.Height)
			require.Equal(t, tt.expectedCctxsMissed, res.CctxsMissed)
			require.Equal(t, tt.expectedCctxsPending, res.CctxsPending)
			require.Equal(t, tt.expectedTotalPending, res.TotalPending)
			require.Equal(t, tt.expectedPastCctxsValue, res.PastCctxsValue)
			require.Equal(t, tt.expectedPendingCctxsValue, res.PendingCctxsValue)
			require.Equal(t, tt.expectedLowestPendingCctxHeight, res.LowestPendingCctxHeight)
		})
	}
}

func TestKeeper_RateLimiterInput_Errors(t *testing.T) {
	t.Run("should fail for empty req", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		_, err := k.RateLimiterInput(ctx, nil)
		require.ErrorContains(t, err, "invalid request")
	})

	t.Run("window must be positive", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		_, err := k.RateLimiterInput(ctx, &types.QueryRateLimiterInputRequest{
			Window: 0, // 0 window
		})
		require.ErrorContains(t, err, "window must be positive")
	})
	t.Run("height out of range", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		// set current height to 0
		ctx = ctx.WithBlockHeight(0)
		_, err := k.RateLimiterInput(ctx, &types.QueryRateLimiterInputRequest{
			Window: 100,
		})
		require.ErrorContains(t, err, "height out of range")
	})

	t.Run("tss not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		// no TSS set
		_, err := k.RateLimiterInput(ctx, &types.QueryRateLimiterInputRequest{
			Window: 100,
		})
		require.ErrorContains(t, err, "tss not found")
	})

	t.Run("asset rates not found", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)

		// Set TSS but no rate limiter flags
		tss := sample.Tss()
		zk.ObserverKeeper.SetTSS(ctx, tss)

		_, err := k.RateLimiterInput(ctx, &types.QueryRateLimiterInputRequest{
			Window: 100,
		})
		require.ErrorContains(t, err, "asset rates not found")
	})

	t.Run("pending nonces not found", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)

		// Set TSS
		tss := sample.Tss()
		zk.ObserverKeeper.SetTSS(ctx, tss)

		// Set rate limiter flags as disabled
		rFlags := sample.RateLimiterFlags()
		k.SetRateLimiterFlags(ctx, rFlags)

		_, err := k.RateLimiterInput(ctx, &types.QueryRateLimiterInputRequest{
			Window: 100,
		})
		require.ErrorContains(t, err, "pending nonces not found")
	})
}

func TestKeeper_ListPendingCctxWithinRateLimit(t *testing.T) {
	// create sample TSS
	tss := sample.Tss()
	zetaChainID := chains.ZetaChainMainnet.ChainId

	// create sample zrc20 addresses for ETH, BTC, USDT
	zrc20ETH := sample.EthAddress().Hex()
	zrc20BTC := sample.EthAddress().Hex()
	zrc20USDT := sample.EthAddress().Hex()

	// create Eth chain 999 mined and 200 pending cctxs for rate limiter test
	// the number 999 is to make it less than `MaxLookbackNonce` so the LoopBackwards gets the chance to hit nonce 0
	ethMinedCctxs := sample.CustomCctxsInBlockRange(
		t,
		1,
		999,
		zetaChainID,
		ethChainID,
		coin.CoinType_Gas,
		"",
		uint64(1e15),
		types.CctxStatus_OutboundMined,
	)
	ethPendingCctxs := sample.CustomCctxsInBlockRange(
		t,
		1000,
		1199,
		zetaChainID,
		ethChainID,
		coin.CoinType_Gas,
		"",
		uint64(1e15),
		types.CctxStatus_PendingOutbound,
	)

	// create Eth chain 999 reverted and 200 pending revert cctxs for rate limiter test
	// these cctxs should be just ignored by the rate limiter as we can't compare their `ObservedExternalHeight` with window boundary
	ethRevertedCctxs := sample.CustomCctxsInBlockRange(
		t,
		1,
		999,
		ethChainID,
		ethChainID,
		coin.CoinType_Gas,
		"",
		uint64(1e15),
		types.CctxStatus_Reverted,
	)
	ethPendingRevertCctxs := sample.CustomCctxsInBlockRange(
		t,
		1000,
		1199,
		ethChainID,
		ethChainID,
		coin.CoinType_Gas,
		"",
		uint64(1e15),
		types.CctxStatus_PendingRevert,
	)

	// create Btc chain 999 mined and 200 pending cctxs for rate limiter test
	// the number 999 is to make it less than `MaxLookbackNonce` so the LoopBackwards gets the chance to hit nonce 0
	btcMinedCctxs := sample.CustomCctxsInBlockRange(
		t,
		1,
		999,
		zetaChainID,
		btcChainID,
		coin.CoinType_Gas,
		"",
		1000,
		types.CctxStatus_OutboundMined,
	)
	btcPendingCctxs := sample.CustomCctxsInBlockRange(
		t,
		1000,
		1199,
		zetaChainID,
		btcChainID,
		coin.CoinType_Gas,
		"",
		1000,
		types.CctxStatus_PendingOutbound,
	)

	// define test cases
	tests := []struct {
		name           string
		fallback       bool
		rateLimitFlags *types.RateLimiterFlags

		// Eth chain cctxs setup
		ethMinedCctxs    []*types.CrossChainTx
		ethPendingCctxs  []*types.CrossChainTx
		ethPendingNonces observertypes.PendingNonces

		// Btc chain cctxs setup
		btcMinedCctxs    []*types.CrossChainTx
		btcPendingCctxs  []*types.CrossChainTx
		btcPendingNonces observertypes.PendingNonces

		// current block height and limit
		currentHeight int64
		queryLimit    uint32

		// expected results
		expectedCctxs          []*types.CrossChainTx
		expectedTotalPending   uint64
		expectedWithdrawWindow int64
		expectedWithdrawRate   string
		rateLimitExceeded      bool
	}{
		{
			name:            "should use fallback query if rate limiter is disabled",
			fallback:        true,
			rateLimitFlags:  nil, // no rate limiter flags set in the keeper
			ethMinedCctxs:   ethMinedCctxs,
			ethPendingCctxs: ethPendingCctxs,
			ethPendingNonces: observertypes.PendingNonces{
				ChainId:   ethChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			btcMinedCctxs:   btcMinedCctxs,
			btcPendingCctxs: btcPendingCctxs,
			btcPendingNonces: observertypes.PendingNonces{
				ChainId:   btcChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			currentHeight:        1199,
			queryLimit:           keeper.MaxPendingCctxs,
			expectedCctxs:        append(append([]*types.CrossChainTx{}, btcPendingCctxs...), ethPendingCctxs...),
			expectedTotalPending: 400,
		},
		{
			name:     "should use fallback query if rate is 0",
			fallback: true,
			rateLimitFlags: createTestRateLimiterFlags(
				500,
				math.NewUint(0),
				zrc20ETH,
				zrc20BTC,
				zrc20USDT,
				"2500",
				"50000",
				"0.8",
			),
			ethMinedCctxs:   ethMinedCctxs,
			ethPendingCctxs: ethPendingCctxs,
			ethPendingNonces: observertypes.PendingNonces{
				ChainId:   ethChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			btcMinedCctxs:   btcMinedCctxs,
			btcPendingCctxs: btcPendingCctxs,
			btcPendingNonces: observertypes.PendingNonces{
				ChainId:   btcChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			currentHeight:        1199,
			queryLimit:           keeper.MaxPendingCctxs,
			expectedCctxs:        append(append([]*types.CrossChainTx{}, btcPendingCctxs...), ethPendingCctxs...),
			expectedTotalPending: 400,
		},
		{
			name: "can retrieve all pending cctx without exceeding rate limit",
			rateLimitFlags: createTestRateLimiterFlags(
				500,
				math.NewUint(10*1e18),
				zrc20ETH,
				zrc20BTC,
				zrc20USDT,
				"2500",
				"50000",
				"0.8",
			),
			ethMinedCctxs:   ethMinedCctxs,
			ethPendingCctxs: ethPendingCctxs,
			ethPendingNonces: observertypes.PendingNonces{
				ChainId:   ethChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			btcMinedCctxs:   btcMinedCctxs,
			btcPendingCctxs: btcPendingCctxs,
			btcPendingNonces: observertypes.PendingNonces{
				ChainId:   btcChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			currentHeight:          1199,
			queryLimit:             keeper.MaxPendingCctxs,
			expectedCctxs:          append(append([]*types.CrossChainTx{}, ethPendingCctxs...), btcPendingCctxs...),
			expectedTotalPending:   400,
			expectedWithdrawWindow: 500,                       // the sliding window
			expectedWithdrawRate:   sdk.NewInt(3e18).String(), // 3 ZETA, (2.5 + 0.5) per block
			rateLimitExceeded:      false,
		},
		{
			name: "can ignore reverted or pending revert cctxs and retrieve all pending cctx without exceeding rate limit",
			rateLimitFlags: createTestRateLimiterFlags(
				500,
				math.NewUint(10*1e18),
				zrc20ETH,
				zrc20BTC,
				zrc20USDT,
				"2500",
				"50000",
				"0.8",
			),
			ethMinedCctxs:   ethRevertedCctxs,      // replace mined cctxs with reverted cctxs, should be ignored
			ethPendingCctxs: ethPendingRevertCctxs, // replace pending cctxs with pending revert cctxs, should be ignored
			ethPendingNonces: observertypes.PendingNonces{
				ChainId:   ethChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			btcMinedCctxs:   btcMinedCctxs,
			btcPendingCctxs: btcPendingCctxs,
			btcPendingNonces: observertypes.PendingNonces{
				ChainId:   btcChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			currentHeight: 1199,
			queryLimit:    keeper.MaxPendingCctxs,
			expectedCctxs: append(
				append([]*types.CrossChainTx{}, ethPendingRevertCctxs...),
				btcPendingCctxs...),
			expectedTotalPending:   400,
			expectedWithdrawWindow: 500,                       // the sliding window
			expectedWithdrawRate:   sdk.NewInt(5e17).String(), // 0.5 ZETA per block, only btc cctxs should be counted
			rateLimitExceeded:      false,
		},
		{
			name: "can loop backwards all the way to endNonce 0",
			rateLimitFlags: createTestRateLimiterFlags(
				500,
				math.NewUint(10*1e18),
				zrc20ETH,
				zrc20BTC,
				zrc20USDT,
				"2500",
				"50000",
				"0.8",
			),
			ethMinedCctxs:   ethMinedCctxs,
			ethPendingCctxs: ethPendingCctxs,
			ethPendingNonces: observertypes.PendingNonces{
				ChainId:   ethChainID,
				NonceLow:  999, // endNonce will be set to 0 (NonceLow - 1000 < 0)
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			btcMinedCctxs:   btcMinedCctxs,
			btcPendingCctxs: btcPendingCctxs,
			btcPendingNonces: observertypes.PendingNonces{
				ChainId:   btcChainID,
				NonceLow:  999, // endNonce will be set to 0 (NonceLow - 1000 < 0)
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			currentHeight:          1199,
			queryLimit:             keeper.MaxPendingCctxs,
			expectedCctxs:          append(append([]*types.CrossChainTx{}, ethPendingCctxs...), btcPendingCctxs...),
			expectedTotalPending:   400,
			expectedWithdrawWindow: 500,                       // the sliding window
			expectedWithdrawRate:   sdk.NewInt(3e18).String(), // 3 ZETA, (2.5 + 0.5) per block
			rateLimitExceeded:      false,
		},
		{
			name: "set a low rate (rate < 2.4 ZETA) to exceed rate limit in backward loop",
			rateLimitFlags: createTestRateLimiterFlags(
				500,
				math.NewUint(2*1e18),
				zrc20ETH,
				zrc20BTC,
				zrc20USDT,
				"2500",
				"50000",
				"0.8",
			),
			ethMinedCctxs:   ethMinedCctxs,
			ethPendingCctxs: ethPendingCctxs,
			ethPendingNonces: observertypes.PendingNonces{
				ChainId:   ethChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			btcMinedCctxs:   btcMinedCctxs,
			btcPendingCctxs: btcPendingCctxs,
			btcPendingNonces: observertypes.PendingNonces{
				ChainId:   btcChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			currentHeight: 1199,
			queryLimit:    keeper.MaxPendingCctxs,
			// return missed cctxs only if rate limit is exceeded
			expectedCctxs: append(
				append([]*types.CrossChainTx{}, ethPendingCctxs[0:100]...),
				btcPendingCctxs[0:100]...),
			expectedTotalPending:   400,
			expectedWithdrawWindow: 500,                       // the sliding window
			expectedWithdrawRate:   sdk.NewInt(3e18).String(), // 3 ZETA, (2.5 + 0.5) per block
			rateLimitExceeded:      true,
		},
		{
			name: "set a lower gRPC request limit and reach the limit of the query in forward loop",
			rateLimitFlags: createTestRateLimiterFlags(
				500,
				math.NewUint(10*1e18),
				zrc20ETH,
				zrc20BTC,
				zrc20USDT,
				"2500",
				"50000",
				"0.8",
			),
			ethMinedCctxs:   ethMinedCctxs,
			ethPendingCctxs: ethPendingCctxs,
			ethPendingNonces: observertypes.PendingNonces{
				ChainId:   ethChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			btcMinedCctxs:   btcMinedCctxs,
			btcPendingCctxs: btcPendingCctxs,
			btcPendingNonces: observertypes.PendingNonces{
				ChainId:   btcChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			currentHeight: 1199,
			queryLimit:    300, // 300 < keeper.MaxPendingCctxs
			expectedCctxs: append(
				append([]*types.CrossChainTx{}, ethPendingCctxs[0:100]...),
				btcPendingCctxs...),
			expectedTotalPending:   400,
			expectedWithdrawWindow: 500,                       // the sliding window
			expectedWithdrawRate:   sdk.NewInt(3e18).String(), // 3 ZETA, (2.5 + 0.5) per block
			rateLimitExceeded:      false,
		},
		{
			name: "set a median rate (2.4 ZETA < rate < 3 ZETA) to exceed rate limit in forward loop",
			rateLimitFlags: createTestRateLimiterFlags(
				500,
				math.NewUint(26*1e17),
				zrc20ETH,
				zrc20BTC,
				zrc20USDT,
				"2500",
				"50000",
				"0.8",
			),
			ethMinedCctxs:   ethMinedCctxs,
			ethPendingCctxs: ethPendingCctxs,
			ethPendingNonces: observertypes.PendingNonces{
				ChainId:   ethChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			btcMinedCctxs:   btcMinedCctxs,
			btcPendingCctxs: btcPendingCctxs,
			btcPendingNonces: observertypes.PendingNonces{
				ChainId:   btcChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			currentHeight: 1199,
			queryLimit:    keeper.MaxPendingCctxs,
			// return missed cctxs only if rate limit is exceeded
			expectedCctxs: append(
				append([]*types.CrossChainTx{}, ethPendingCctxs[0:100]...),
				btcPendingCctxs[0:100]...),
			expectedTotalPending:   400,
			expectedWithdrawWindow: 500,                       // the sliding window
			expectedWithdrawRate:   sdk.NewInt(3e18).String(), // 3 ZETA, (2.5 + 0.5) per block
			rateLimitExceeded:      true,
		},
		{
			// the pending cctxs window is wider than the rate limiter sliding window in this test case.
			name: "set low rate and narrow window to early exceed rate limit in forward loop",
			// the left boundary will be 1149 (currentHeight-50), the pending nonces [1099, 1148] fall beyond the left boundary.
			// `pendingCctxWindow` is 100 which is wider than rate limiter window 50.
			//  give a block rate of 2 ZETA/block, the max value allowed should be 100 * 2 = 200 ZETA
			rateLimitFlags: createTestRateLimiterFlags(
				50,
				math.NewUint(2*1e18),
				zrc20ETH,
				zrc20BTC,
				zrc20USDT,
				"2500",
				"50000",
				"0.8",
			),
			ethMinedCctxs:   ethMinedCctxs,
			ethPendingCctxs: ethPendingCctxs,
			ethPendingNonces: observertypes.PendingNonces{
				ChainId:   ethChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			btcMinedCctxs:   btcMinedCctxs,
			btcPendingCctxs: btcPendingCctxs,
			btcPendingNonces: observertypes.PendingNonces{
				ChainId:   btcChainID,
				NonceLow:  1099,
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			currentHeight: 1199,
			queryLimit:    keeper.MaxPendingCctxs,
			// return missed cctxs only if rate limit is exceeded
			expectedCctxs: append(
				append([]*types.CrossChainTx{}, ethPendingCctxs[0:100]...),
				btcPendingCctxs[0:100]...),
			expectedTotalPending:   400,
			expectedWithdrawWindow: 100,                       // 100 > sliding window 50
			expectedWithdrawRate:   sdk.NewInt(3e18).String(), // 3 ZETA, (2.5 + 0.5) per block
			rateLimitExceeded:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create test keepers
			k, ctx, _, zk := keepertest.CrosschainKeeper(t)

			// Set TSS
			zk.ObserverKeeper.SetTSS(ctx, tss)

			// Set foreign coins
			assetUSDT := sample.EthAddress().Hex()
			setupForeignCoins(t, ctx, zk, zrc20ETH, zrc20BTC, zrc20USDT, assetUSDT)

			// Set rate limiter flags
			if tt.rateLimitFlags != nil {
				k.SetRateLimiterFlags(ctx, *tt.rateLimitFlags)
			}

			// Set Eth chain mined cctxs, pending ccxts and pending nonces
			setCctxsInKeeper(ctx, *k, zk, tss, tt.ethMinedCctxs)
			setCctxsInKeeper(ctx, *k, zk, tss, tt.ethPendingCctxs)
			zk.ObserverKeeper.SetPendingNonces(ctx, tt.ethPendingNonces)

			// Set Btc chain mined cctxs, pending ccxts and pending nonces
			setCctxsInKeeper(ctx, *k, zk, tss, tt.btcMinedCctxs)
			setCctxsInKeeper(ctx, *k, zk, tss, tt.btcPendingCctxs)
			zk.ObserverKeeper.SetPendingNonces(ctx, tt.btcPendingNonces)

			// Set current block height
			ctx = ctx.WithBlockHeight(tt.currentHeight)

			// Query pending cctxs
			res, err := k.ListPendingCctxWithinRateLimit(
				ctx,
				&types.QueryListPendingCctxWithinRateLimitRequest{Limit: tt.queryLimit},
			)
			require.NoError(t, err)
			require.EqualValues(t, tt.expectedCctxs, res.CrossChainTx)
			require.Equal(t, tt.expectedTotalPending, res.TotalPending)

			// check rate limiter related fields only if it's not a fallback query
			if !tt.fallback {
				require.Equal(t, tt.expectedWithdrawWindow, res.CurrentWithdrawWindow)
				require.Equal(t, tt.expectedWithdrawRate, res.CurrentWithdrawRate)
				require.Equal(t, tt.rateLimitExceeded, res.RateLimitExceeded)
			}
		})
	}
}

func TestKeeper_ListPendingCctxWithinRateLimit_Errors(t *testing.T) {
	t.Run("should fail for empty req", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		_, err := k.ListPendingCctxWithinRateLimit(ctx, nil)
		require.ErrorContains(t, err, "invalid request")
	})
	t.Run("height out of range", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		// Set rate limiter flags as disabled
		rFlags := sample.RateLimiterFlags()
		k.SetRateLimiterFlags(ctx, rFlags)

		ctx = ctx.WithBlockHeight(0)
		_, err := k.ListPendingCctxWithinRateLimit(ctx, &types.QueryListPendingCctxWithinRateLimitRequest{})
		require.ErrorContains(t, err, "height out of range")
	})
	t.Run("tss not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		// Set rate limiter flags as disabled
		rFlags := sample.RateLimiterFlags()
		k.SetRateLimiterFlags(ctx, rFlags)

		_, err := k.ListPendingCctxWithinRateLimit(ctx, &types.QueryListPendingCctxWithinRateLimitRequest{})
		require.ErrorContains(t, err, "tss not found")
	})
	t.Run("pending nonces not found", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)

		// Set rate limiter flags as disabled
		rFlags := sample.RateLimiterFlags()
		k.SetRateLimiterFlags(ctx, rFlags)

		// Set TSS
		tss := sample.Tss()
		zk.ObserverKeeper.SetTSS(ctx, tss)

		_, err := k.ListPendingCctxWithinRateLimit(ctx, &types.QueryListPendingCctxWithinRateLimitRequest{})
		require.ErrorContains(t, err, "pending nonces not found")
	})
}

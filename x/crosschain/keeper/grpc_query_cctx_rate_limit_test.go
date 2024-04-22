package keeper_test

import (
	"fmt"
	"strings"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
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

// createCctxsWithCoinTypeAndHeightRange
//   - create 1 cctx per block from lowBlock to highBlock (inclusive)
//
// return created cctxs
func createCctxsWithCoinTypeAndHeightRange(
	t *testing.T,
	lowBlock uint64,
	highBlock uint64,
	chainID int64,
	coinType coin.CoinType,
	asset string,
	amount uint64,
	status types.CctxStatus,
) (cctxs []*types.CrossChainTx) {
	// create 1 pending cctxs per block
	for i := lowBlock; i <= highBlock; i++ {
		nonce := i - 1
		cctx := sample.CrossChainTx(t, fmt.Sprintf("%d-%d", chainID, nonce))
		cctx.CctxStatus.Status = status
		cctx.InboundTxParams.SenderChainId = chainID
		cctx.InboundTxParams.CoinType = coinType
		cctx.InboundTxParams.Asset = asset
		cctx.InboundTxParams.InboundTxObservedExternalHeight = i
		cctx.GetCurrentOutTxParam().ReceiverChainId = chainID
		cctx.GetCurrentOutTxParam().Amount = sdk.NewUint(amount)
		cctx.GetCurrentOutTxParam().OutboundTxTssNonce = nonce
		cctxs = append(cctxs, cctx)
	}
	return cctxs
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
			ChainId: cctx.InboundTxParams.SenderChainId,
			// #nosec G701 always in range for tests
			Nonce:     int64(cctx.GetCurrentOutTxParam().OutboundTxTssNonce),
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

func Test_ConvertCctxValue(t *testing.T) {
	// chain IDs
	ethChainID := getValidEthChainID()
	btcChainID := getValidBtcChainID()

	// zrc20 addresses for ETH, BTC, USDT and asset for USDT
	zrc20ETH := sample.EthAddress().Hex()
	zrc20BTC := sample.EthAddress().Hex()
	zrc20USDT := sample.EthAddress().Hex()
	assetUSDT := sample.EthAddress().Hex()

	k, ctx, _, zk := keepertest.CrosschainKeeper(t)

	// Set TSS
	tss := sample.Tss()
	zk.ObserverKeeper.SetTSS(ctx, tss)

	// Set foreign coins
	setupForeignCoins(t, ctx, zk, zrc20ETH, zrc20BTC, zrc20USDT, assetUSDT)

	// Set rate limiter flags
	rateLimiterFlags := createTestRateLimiterFlags(500, math.NewUint(5000), zrc20ETH, zrc20BTC, zrc20USDT, "2500", "50000", "0.8")
	k.SetRateLimiterFlags(ctx, *rateLimiterFlags)

	// get rate limiter rates
	gasCoinRates, erc20CoinRates := k.GetRateLimiterRates(ctx)
	foreignCoinMap := zk.FungibleKeeper.GetAllForeignCoinMap(ctx)

	t.Run("should convert cctx ZETA value correctly", func(t *testing.T) {
		// create cctx with 0.3 ZETA
		cctx := sample.CrossChainTx(t, fmt.Sprintf("%d-%d", ethChainID, 1))
		cctx.InboundTxParams.CoinType = coin.CoinType_Zeta
		cctx.InboundTxParams.Asset = ""
		cctx.GetCurrentOutTxParam().Amount = sdk.NewUint(3e17) // 0.3 ZETA

		// convert cctx value
		value := keeper.ConvertCctxValue(ethChainID, cctx, gasCoinRates, erc20CoinRates, foreignCoinMap)
		require.Equal(t, sdk.MustNewDecFromStr("0.3"), value)
	})
	t.Run("should convert cctx ETH value correctly", func(t *testing.T) {
		// create cctx with 0.003 ETH
		cctx := sample.CrossChainTx(t, fmt.Sprintf("%d-%d", ethChainID, 1))
		cctx.InboundTxParams.CoinType = coin.CoinType_Gas
		cctx.InboundTxParams.Asset = ""
		cctx.GetCurrentOutTxParam().Amount = sdk.NewUint(3e15) // 0.003 ETH

		// convert cctx value
		value := keeper.ConvertCctxValue(ethChainID, cctx, gasCoinRates, erc20CoinRates, foreignCoinMap)
		require.Equal(t, sdk.MustNewDecFromStr("7.5"), value)
	})
	t.Run("should convert cctx BTC value correctly", func(t *testing.T) {
		// create cctx with 0.0007 BTC
		cctx := sample.CrossChainTx(t, fmt.Sprintf("%d-%d", btcChainID, 1))
		cctx.InboundTxParams.CoinType = coin.CoinType_Gas
		cctx.InboundTxParams.Asset = ""
		cctx.GetCurrentOutTxParam().Amount = sdk.NewUint(70000) // 0.0007 BTC

		// convert cctx value
		value := keeper.ConvertCctxValue(btcChainID, cctx, gasCoinRates, erc20CoinRates, foreignCoinMap)
		require.Equal(t, sdk.MustNewDecFromStr("35.0"), value)
	})
	t.Run("should convert cctx USDT value correctly", func(t *testing.T) {
		// create cctx with 3 USDT
		cctx := sample.CrossChainTx(t, fmt.Sprintf("%d-%d", ethChainID, 1))
		cctx.InboundTxParams.CoinType = coin.CoinType_ERC20
		cctx.InboundTxParams.Asset = assetUSDT
		cctx.GetCurrentOutTxParam().Amount = sdk.NewUint(3e6) // 3 USDT

		// convert cctx value
		value := keeper.ConvertCctxValue(ethChainID, cctx, gasCoinRates, erc20CoinRates, foreignCoinMap)
		require.Equal(t, sdk.MustNewDecFromStr("2.4"), value)
	})
	t.Run("should return 0 if no rate found for chainID", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, fmt.Sprintf("%d-%d", ethChainID, 1))
		cctx.InboundTxParams.CoinType = coin.CoinType_ERC20
		cctx.GetCurrentOutTxParam().Amount = sdk.NewUint(100)

		// use nil erc20CoinRates map to convert cctx value
		value := keeper.ConvertCctxValue(ethChainID, cctx, gasCoinRates, nil, foreignCoinMap)
		require.Equal(t, sdk.NewDec(0), value)
	})
	t.Run("should return 0 if coinType is CoinType_Cmd", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, fmt.Sprintf("%d-%d", ethChainID, 1))
		cctx.InboundTxParams.CoinType = coin.CoinType_Cmd
		cctx.GetCurrentOutTxParam().Amount = sdk.NewUint(100)

		// convert cctx value
		value := keeper.ConvertCctxValue(ethChainID, cctx, gasCoinRates, erc20CoinRates, foreignCoinMap)
		require.Equal(t, sdk.NewDec(0), value)
	})
	t.Run("should return 0 on nil rate or rate <= 0", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, fmt.Sprintf("%d-%d", ethChainID, 1))
		cctx.InboundTxParams.CoinType = coin.CoinType_Gas
		cctx.GetCurrentOutTxParam().Amount = sdk.NewUint(100)

		// use nil gasCoinRates map to convert cctx value
		value := keeper.ConvertCctxValue(ethChainID, cctx, nil, erc20CoinRates, foreignCoinMap)
		require.Equal(t, sdk.NewDec(0), value)

		// set rate to 0
		zeroCoinRates, _ := k.GetRateLimiterRates(ctx)
		zeroCoinRates[ethChainID] = sdk.NewDec(0)

		// convert cctx value
		value = keeper.ConvertCctxValue(ethChainID, cctx, zeroCoinRates, erc20CoinRates, foreignCoinMap)
		require.Equal(t, sdk.NewDec(0), value)

		// set rate to -1
		negativeCoinRates, _ := k.GetRateLimiterRates(ctx)
		negativeCoinRates[ethChainID] = sdk.NewDec(-1)

		// convert cctx value
		value = keeper.ConvertCctxValue(ethChainID, cctx, negativeCoinRates, erc20CoinRates, foreignCoinMap)
		require.Equal(t, sdk.NewDec(0), value)
	})
	t.Run("should return 0 if no coin found for chainID", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, fmt.Sprintf("%d-%d", ethChainID, 1))
		cctx.InboundTxParams.CoinType = coin.CoinType_Gas
		cctx.GetCurrentOutTxParam().Amount = sdk.NewUint(100)

		// use empty foreignCoinMap to convert cctx value
		value := keeper.ConvertCctxValue(ethChainID, cctx, gasCoinRates, erc20CoinRates, nil)
		require.Equal(t, sdk.NewDec(0), value)
	})
	t.Run("should return 0 if no coin found for asset", func(t *testing.T) {
		cctx := sample.CrossChainTx(t, fmt.Sprintf("%d-%d", ethChainID, 1))
		cctx.InboundTxParams.CoinType = coin.CoinType_ERC20
		cctx.InboundTxParams.Asset = assetUSDT
		cctx.GetCurrentOutTxParam().Amount = sdk.NewUint(100)

		// delete assetUSDT from foreignCoinMap for ethChainID
		tempCoinMap := zk.FungibleKeeper.GetAllForeignCoinMap(ctx)
		delete(tempCoinMap[ethChainID], strings.ToLower(assetUSDT))

		// convert cctx value
		value := keeper.ConvertCctxValue(ethChainID, cctx, gasCoinRates, erc20CoinRates, tempCoinMap)
		require.Equal(t, sdk.NewDec(0), value)
	})
}

func TestKeeper_ListPendingCctxWithinRateLimit(t *testing.T) {
	// create sample TSS
	tss := sample.Tss()

	// create sample zrc20 addresses for ETH, BTC, USDT
	zrc20ETH := sample.EthAddress().Hex()
	zrc20BTC := sample.EthAddress().Hex()
	zrc20USDT := sample.EthAddress().Hex()

	// create Eth chain mined and pending cctxs for rate limiter test
	ethMindedCctxs := createCctxsWithCoinTypeAndHeightRange(t, 1, 999, ethChainID, coin.CoinType_Gas, "", uint64(1e15), types.CctxStatus_OutboundMined)
	ethPendingCctxs := createCctxsWithCoinTypeAndHeightRange(t, 1000, 1199, ethChainID, coin.CoinType_Gas, "", uint64(1e15), types.CctxStatus_PendingOutbound)

	// create Btc chain mined and pending cctxs for rate limiter test
	btcMinedCctxs := createCctxsWithCoinTypeAndHeightRange(t, 1, 999, btcChainID, coin.CoinType_Gas, "", 1000, types.CctxStatus_OutboundMined)
	btcPendingCctxs := createCctxsWithCoinTypeAndHeightRange(t, 1000, 1199, btcChainID, coin.CoinType_Gas, "", 1000, types.CctxStatus_PendingOutbound)

	// define test cases
	tests := []struct {
		name           string
		rateLimitFlags *types.RateLimiterFlags

		// Eth chain cctxs setup
		ethMindedCctxs   []*types.CrossChainTx
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
		expectedCctxs        []*types.CrossChainTx
		expectedTotalPending uint64
		expectedTotalValue   uint64
		rateLimitExceeded    bool
	}{
		{
			name:            "should use fallback query if rate limiter is disabled",
			rateLimitFlags:  nil, // no rate limiter flags set in the keeper
			ethMindedCctxs:  ethMindedCctxs,
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
			name:            "can retrieve pending cctx in range without exceeding rate limit",
			rateLimitFlags:  createTestRateLimiterFlags(500, math.NewUint(5000), zrc20ETH, zrc20BTC, zrc20USDT, "2500", "50000", "0.8"),
			ethMindedCctxs:  ethMindedCctxs,
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
			expectedCctxs:        append(append([]*types.CrossChainTx{}, ethPendingCctxs...), btcPendingCctxs...),
			expectedTotalPending: 400,
			expectedTotalValue:   1500, // 500 (window) * (2.5 + 0.5)
			rateLimitExceeded:    false,
		},
		{
			name:            "can loop backwards all the way to endNonce 0",
			rateLimitFlags:  createTestRateLimiterFlags(500, math.NewUint(5000), zrc20ETH, zrc20BTC, zrc20USDT, "2500", "50000", "0.8"),
			ethMindedCctxs:  ethMindedCctxs,
			ethPendingCctxs: ethPendingCctxs,
			ethPendingNonces: observertypes.PendingNonces{
				ChainId:   ethChainID,
				NonceLow:  999, // endNonce will set to 0 as NonceLow - 1000 < 0
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			btcMinedCctxs:   btcMinedCctxs,
			btcPendingCctxs: btcPendingCctxs,
			btcPendingNonces: observertypes.PendingNonces{
				ChainId:   btcChainID,
				NonceLow:  999, // endNonce will set to 0 as NonceLow - 1000 < 0
				NonceHigh: 1199,
				Tss:       tss.TssPubkey,
			},
			currentHeight:        1199,
			queryLimit:           keeper.MaxPendingCctxs,
			expectedCctxs:        append(append([]*types.CrossChainTx{}, ethPendingCctxs...), btcPendingCctxs...),
			expectedTotalPending: 400,
			expectedTotalValue:   1500, // 500 (window) * (2.5 + 0.5)
			rateLimitExceeded:    false,
		},
		{
			name:            "set a low rate (< 1200) to early break the LoopBackwards with criteria #2",
			rateLimitFlags:  createTestRateLimiterFlags(500, math.NewUint(1000), zrc20ETH, zrc20BTC, zrc20USDT, "2500", "50000", "0.8"), // 1000 < 1200
			ethMindedCctxs:  ethMindedCctxs,
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
			expectedCctxs:        append(append([]*types.CrossChainTx{}, ethPendingCctxs[0:100]...), btcPendingCctxs[0:100]...),
			expectedTotalPending: 400,
			rateLimitExceeded:    true,
		},
		{
			name: "set high rate and big window to early to break inner loop with the criteria #1",
			// The left boundary will be 49 (currentHeight-1150), which will be less than the endNonce 99 (1099 - 1000)
			rateLimitFlags:  createTestRateLimiterFlags(1150, math.NewUint(10000), zrc20ETH, zrc20BTC, zrc20USDT, "2500", "50000", "0.8"),
			ethMindedCctxs:  ethMindedCctxs,
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
			expectedCctxs:        append(append([]*types.CrossChainTx{}, ethPendingCctxs...), btcPendingCctxs...),
			expectedTotalPending: 400,
			expectedTotalValue:   3450, // 1150 (window) * (2.5 + 0.5)
			rateLimitExceeded:    false,
		},
		{
			name:            "set lower request limit to early break the LoopForwards loop",
			rateLimitFlags:  createTestRateLimiterFlags(500, math.NewUint(5000), zrc20ETH, zrc20BTC, zrc20USDT, "2500", "50000", "0.8"),
			ethMindedCctxs:  ethMindedCctxs,
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
			queryLimit:           300, // 300 < keeper.MaxPendingCctxs
			expectedCctxs:        append(append([]*types.CrossChainTx{}, ethPendingCctxs[0:100]...), btcPendingCctxs...),
			expectedTotalPending: 400,
			expectedTotalValue:   1250, // 500 * 0.5 + 400 * 2.5
			rateLimitExceeded:    false,
		},
		{
			name:            "set rate to middle value (1200 < rate < 1500) to early break the LoopForwards loop with criteria #2",
			rateLimitFlags:  createTestRateLimiterFlags(500, math.NewUint(1300), zrc20ETH, zrc20BTC, zrc20USDT, "2500", "50000", "0.8"), // 1200 < 1300 < 1500
			ethMindedCctxs:  ethMindedCctxs,
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
			// 120 ETH cctx + 200 BTC cctx
			expectedCctxs:        append(append([]*types.CrossChainTx{}, ethPendingCctxs[0:120]...), btcPendingCctxs...),
			expectedTotalPending: 400,
			rateLimitExceeded:    true,
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
			setCctxsInKeeper(ctx, *k, zk, tss, tt.ethMindedCctxs)
			setCctxsInKeeper(ctx, *k, zk, tss, tt.ethPendingCctxs)
			zk.ObserverKeeper.SetPendingNonces(ctx, tt.ethPendingNonces)

			// Set Btc chain mined cctxs, pending ccxts and pending nonces
			setCctxsInKeeper(ctx, *k, zk, tss, tt.btcMinedCctxs)
			setCctxsInKeeper(ctx, *k, zk, tss, tt.btcPendingCctxs)
			zk.ObserverKeeper.SetPendingNonces(ctx, tt.btcPendingNonces)

			// Set current block height
			ctx = ctx.WithBlockHeight(tt.currentHeight)

			// Query pending cctxs
			res, err := k.ListPendingCctxWithinRateLimit(ctx, &types.QueryListPendingCctxWithinRateLimitRequest{Limit: tt.queryLimit})
			require.NoError(t, err)
			require.EqualValues(t, tt.expectedCctxs, res.CrossChainTx)
			require.EqualValues(t, tt.expectedTotalPending, res.TotalPending)

			// check rate limiter related fields only if rate limiter flags is enabled
			if tt.rateLimitFlags != nil {
				if !tt.rateLimitExceeded {
					require.EqualValues(t, tt.expectedTotalValue, res.ValueWithinWindow)
				} else {
					require.True(t, res.ValueWithinWindow > tt.rateLimitFlags.Rate.Uint64())
				}
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

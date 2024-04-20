package keeper_test

import (
	"fmt"
	"sort"
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

// createTestRateLimiterFlags creates a custom rate limiter flags
func createTestRateLimiterFlags(
	zrc20ETH string,
	zrc20BTC string,
	zrc20USDT string,
	ethRate string,
	btcRate string,
	usdtRate string,
) types.RateLimiterFlags {
	var rateLimiterFlags = types.RateLimiterFlags{
		Enabled: true,
		Window:  100,                // 100 zeta blocks, 10 minutes
		Rate:    math.NewUint(1000), // 1000 ZETA
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
	return rateLimiterFlags
}

// createCctxWithCopyTypeAndBlockRange
//   - create 1 cctx per block from lowBlock to highBlock (inclusive)
//
// return created cctxs
func createCctxWithCopyTypeAndHeightRange(
	t *testing.T,
	ctx sdk.Context,
	k keeper.Keeper,
	zk keepertest.ZetaKeepers,
	tss observertypes.TSS,
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
		cctx.GetCurrentOutTxParam().Amount = sdk.NewUint(amount)
		cctx.GetCurrentOutTxParam().OutboundTxTssNonce = nonce
		k.SetCrossChainTx(ctx, *cctx)
		zk.ObserverKeeper.SetNonceToCctx(ctx, observertypes.NonceToCctx{
			ChainId: chainID,
			// #nosec G701 always in range for tests
			Nonce:     int64(nonce),
			CctxIndex: cctx.Index,
			Tss:       tss.TssPubkey,
		})
		cctxs = append(cctxs, cctx)
	}
	return cctxs
}

// setPendingNonces sets the pending nonces for the given chainID
func setPendingNonces(
	ctx sdk.Context,
	zk keepertest.ZetaKeepers,
	chainID int64,
	nonceLow int64,
	nonceHigh int64,
	tssPubKey string,
) {
	zk.ObserverKeeper.SetPendingNonces(ctx, observertypes.PendingNonces{
		ChainId:   chainID,
		NonceLow:  nonceLow,
		NonceHigh: nonceHigh,
		Tss:       tssPubKey,
	})
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
	rateLimiterFlags := createTestRateLimiterFlags(zrc20ETH, zrc20BTC, zrc20USDT, "2500", "50000", "0.8")
	k.SetRateLimiterFlags(ctx, rateLimiterFlags)

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
	// chain IDs
	ethChainID := getValidEthChainID()

	// define cctx status
	statusPending := types.CctxStatus_PendingOutbound
	statusMined := types.CctxStatus_OutboundMined

	t.Run("should fail for empty req", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		_, err := k.ListPendingCctxWithinRateLimit(ctx, nil)
		require.ErrorContains(t, err, "invalid request")
	})
	t.Run("should use fallback query", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)

		// Set TSS
		tss := sample.Tss()
		zk.ObserverKeeper.SetTSS(ctx, tss)

		// Set rate limiter flags as disabled
		rateLimiterFlags := sample.RateLimiterFlags()
		rateLimiterFlags.Enabled = false
		k.SetRateLimiterFlags(ctx, rateLimiterFlags)

		// Create cctxs [0~999] and [1000~1199] for Eth chain, 0.001 ETH per cctx
		_ = createCctxWithCopyTypeAndHeightRange(t, ctx, *k, zk, tss, 1, 1000, ethChainID, coin.CoinType_Gas, "", uint64(1e15), statusMined)
		cctxETH := createCctxWithCopyTypeAndHeightRange(t, ctx, *k, zk, tss, 1001, 1200, ethChainID, coin.CoinType_Gas, "", uint64(1e15), statusPending)

		// Set Eth chain pending nonces which contains 100 missing cctxs
		setPendingNonces(ctx, zk, ethChainID, 1100, 1200, tss.TssPubkey)

		// Query pending cctxs use small limit
		res, err := k.ListPendingCctxWithinRateLimit(ctx, &types.QueryListPendingCctxWithinRateLimitRequest{Limit: 100})
		require.NoError(t, err)
		require.Equal(t, 100, len(res.CrossChainTx))

		// sort res.CrossChainTx by outbound nonce ascending so that we can compare with cctxETH
		sort.Slice(res.CrossChainTx, func(i, j int) bool {
			return res.CrossChainTx[i].GetCurrentOutTxParam().OutboundTxTssNonce < res.CrossChainTx[j].GetCurrentOutTxParam().OutboundTxTssNonce
		})
		require.EqualValues(t, cctxETH[:100], res.CrossChainTx)
		require.EqualValues(t, uint64(200), res.TotalPending)

		// Query pending cctxs use max limit
		res, err = k.ListPendingCctxWithinRateLimit(ctx, &types.QueryListPendingCctxWithinRateLimitRequest{Limit: keeper.MaxPendingCctxs})
		require.NoError(t, err)
		require.Equal(t, 200, len(res.CrossChainTx))

		// sort res.CrossChainTx by outbound nonce ascending so that we can compare with cctxETH
		sort.Slice(res.CrossChainTx, func(i, j int) bool {
			return res.CrossChainTx[i].GetCurrentOutTxParam().OutboundTxTssNonce < res.CrossChainTx[j].GetCurrentOutTxParam().OutboundTxTssNonce
		})
		require.EqualValues(t, cctxETH, res.CrossChainTx)
		require.EqualValues(t, uint64(200), res.TotalPending)
	})
}

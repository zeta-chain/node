package types_test

import (
	"fmt"
	"strings"
	"testing"

	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestRateLimiterFlags_Validate(t *testing.T) {
	dec, err := sdk.NewDecFromStr("0.00042")
	require.NoError(t, err)
	duplicatedAddress := sample.EthAddress().String()

	tt := []struct {
		name  string
		flags types.RateLimiterFlags
		isErr bool
	}{
		{
			name: "valid flags",
			flags: types.RateLimiterFlags{
				Enabled: true,
				Window:  42,
				Rate:    sdk.NewUint(42),
				Conversions: []types.Conversion{
					{
						Zrc20: sample.EthAddress().String(),
						Rate:  sdk.NewDec(42),
					},
					{
						Zrc20: sample.EthAddress().String(),
						Rate:  dec,
					},
				},
			},
		},
		{
			name:  "empty is valid",
			flags: types.RateLimiterFlags{},
		},
		{
			name: "invalid zrc20 address",
			flags: types.RateLimiterFlags{
				Enabled: true,
				Window:  42,
				Rate:    sdk.NewUint(42),
				Conversions: []types.Conversion{
					{
						Zrc20: "invalid",
						Rate:  sdk.NewDec(42),
					},
					{
						Zrc20: sample.EthAddress().String(),
						Rate:  dec,
					},
				},
			},
			isErr: true,
		},
		{
			name: "duplicated conversion",
			flags: types.RateLimiterFlags{
				Enabled: true,
				Window:  42,
				Rate:    sdk.NewUint(42),
				Conversions: []types.Conversion{
					{
						Zrc20: duplicatedAddress,
						Rate:  sdk.NewDec(42),
					},
					{
						Zrc20: duplicatedAddress,
						Rate:  dec,
					},
				},
			},
			isErr: true,
		},
		{
			name: "invalid conversion rate",
			flags: types.RateLimiterFlags{
				Enabled: true,
				Window:  42,
				Rate:    sdk.NewUint(42),
				Conversions: []types.Conversion{
					{
						Zrc20: sample.EthAddress().String(),
						Rate:  sdk.NewDec(42),
					},
					{
						Zrc20: sample.EthAddress().String(),
					},
				},
			},
			isErr: true,
		},
		{
			name: "negative window",
			flags: types.RateLimiterFlags{
				Enabled: true,
				Window:  -1,
			},
			isErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.flags.Validate()
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}

}

func TestRateLimiterFlags_GetConversionRate(t *testing.T) {
	dec, err := sdk.NewDecFromStr("0.00042")
	require.NoError(t, err)
	address := sample.EthAddress().String()

	tt := []struct {
		name       string
		flags      types.RateLimiterFlags
		zrc20      string
		expected   sdk.Dec
		shouldFind bool
	}{
		{
			name: "valid conversion",
			flags: types.RateLimiterFlags{
				Enabled: true,
				Window:  42,
				Rate:    sdk.NewUint(42),
				Conversions: []types.Conversion{
					{
						Zrc20: address,
						Rate:  sdk.NewDec(42),
					},
					{
						Zrc20: sample.EthAddress().String(),
						Rate:  dec,
					},
				},
			},
			zrc20:      address,
			expected:   sdk.NewDec(42),
			shouldFind: true,
		},
		{
			name: "not found",
			flags: types.RateLimiterFlags{
				Enabled: true,
				Window:  42,
				Rate:    sdk.NewUint(42),
				Conversions: []types.Conversion{
					{
						Zrc20: sample.EthAddress().String(),
						Rate:  sdk.NewDec(42),
					},
					{
						Zrc20: sample.EthAddress().String(),
						Rate:  dec,
					},
				},
			},
			zrc20:      address,
			expected:   sdk.NewDec(0),
			shouldFind: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual, found := tc.flags.GetConversionRate(tc.zrc20)
			require.Equal(t, tc.expected, actual)
			require.Equal(t, tc.shouldFind, found)
		})
	}
}

func TestBuildAssetRateMapFromList(t *testing.T) {
	// define asset rate list
	assetRates := []types.AssetRate{
		{
			ChainId:  1,
			Asset:    "eth",
			Decimals: 18,
			CoinType: coin.CoinType_Gas,
			Rate:     sdk.NewDec(1),
		},
		{
			ChainId:  1,
			Asset:    "usdt",
			Decimals: 6,
			CoinType: coin.CoinType_ERC20,
			Rate:     sdk.NewDec(2),
		},
		{
			ChainId:  2,
			Asset:    "btc",
			Decimals: 8,
			CoinType: coin.CoinType_Gas,
			Rate:     sdk.NewDec(3),
		},
	}

	// build asset rate map
	gasAssetRateMap, erc20AssetRateMap := types.BuildAssetRateMapFromList(assetRates)

	// check length
	require.Equal(t, 2, len(gasAssetRateMap))
	require.Equal(t, 1, len(erc20AssetRateMap))
	require.Equal(t, 1, len(erc20AssetRateMap[1]))

	// check gas asset rate map
	require.EqualValues(t, assetRates[0], gasAssetRateMap[1])
	require.EqualValues(t, assetRates[2], gasAssetRateMap[2])

	// check erc20 asset rate map
	require.EqualValues(t, assetRates[1], erc20AssetRateMap[1]["usdt"])
}

func TestConvertCctxValue(t *testing.T) {
	// chain IDs
	ethChainID := chains.GoerliLocalnetChain.ChainId
	btcChainID := chains.BtcRegtestChain.ChainId

	// setup test asset rates
	assetETH := sample.EthAddress().Hex()
	assetBTC := sample.EthAddress().Hex()
	assetUSDT := sample.EthAddress().Hex()
	assetRateList := []types.AssetRate{
		sample.CustomAssetRate(ethChainID, assetETH, 18, coin.CoinType_Gas, sdk.NewDec(2500)),
		sample.CustomAssetRate(btcChainID, assetBTC, 8, coin.CoinType_Gas, sdk.NewDec(50000)),
		sample.CustomAssetRate(ethChainID, assetUSDT, 6, coin.CoinType_ERC20, sdk.MustNewDecFromStr("0.8")),
	}
	gasAssetRateMap, erc20AssetRateMap := types.BuildAssetRateMapFromList(assetRateList)

	// test cases
	tests := []struct {
		name string

		// input
		chainID         int64
		coinType        coin.CoinType
		asset           string
		amount          math.Uint
		gasAssetRates   map[int64]types.AssetRate
		erc20AssetRates map[int64]map[string]types.AssetRate

		// output
		expectedValue sdkmath.Int
	}{
		{
			name:            "should convert cctx ZETA value correctly",
			chainID:         ethChainID,
			coinType:        coin.CoinType_Zeta,
			asset:           "",
			amount:          sdk.NewUint(3e17), // 0.3 ZETA
			gasAssetRates:   gasAssetRateMap,
			erc20AssetRates: erc20AssetRateMap,
			expectedValue:   sdk.NewInt(3e17),
		},
		{
			name:            "should convert cctx ETH value correctly",
			chainID:         ethChainID,
			coinType:        coin.CoinType_Gas,
			asset:           "",
			amount:          sdk.NewUint(3e15), // 0.003 ETH
			gasAssetRates:   gasAssetRateMap,
			erc20AssetRates: erc20AssetRateMap,
			expectedValue:   sdk.NewInt(75e17), // 0.003 ETH * 2500 = 7.5 ZETA
		},
		{
			name:            "should convert cctx BTC value correctly",
			chainID:         btcChainID,
			coinType:        coin.CoinType_Gas,
			asset:           "",
			amount:          sdk.NewUint(70000), // 0.0007 BTC
			gasAssetRates:   gasAssetRateMap,
			erc20AssetRates: erc20AssetRateMap,
			expectedValue:   sdk.NewInt(35).Mul(sdk.NewInt(1e18)), // 0.0007 BTC * 50000 = 35.0 ZETA
		},
		{
			name:            "should convert cctx USDT value correctly",
			chainID:         ethChainID,
			coinType:        coin.CoinType_ERC20,
			asset:           assetUSDT,
			amount:          sdk.NewUint(3e6), // 3 USDT
			gasAssetRates:   gasAssetRateMap,
			erc20AssetRates: erc20AssetRateMap,
			expectedValue:   sdk.NewInt(24e17), // 3 USDT * 0.8 = 2.4 ZETA
		},
		{
			name:            "should return 0 if no gas asset rate found for chainID",
			chainID:         ethChainID,
			coinType:        coin.CoinType_Gas,
			asset:           "",
			amount:          sdk.NewUint(100),
			gasAssetRates:   nil,
			erc20AssetRates: erc20AssetRateMap,
			expectedValue:   sdk.NewInt(0),
		},
		{
			name:            "should return 0 if no erc20 asset rate found for chainID",
			chainID:         ethChainID,
			coinType:        coin.CoinType_ERC20,
			asset:           assetUSDT,
			amount:          sdk.NewUint(100),
			gasAssetRates:   gasAssetRateMap,
			erc20AssetRates: nil,
			expectedValue:   sdk.NewInt(0),
		},
		{
			name:            "should return 0 if coinType is CoinType_Cmd",
			chainID:         ethChainID,
			coinType:        coin.CoinType_Cmd,
			asset:           "",
			amount:          sdk.NewUint(100),
			gasAssetRates:   gasAssetRateMap,
			erc20AssetRates: erc20AssetRateMap,
			expectedValue:   sdk.NewInt(0),
		},
		{
			name:     "should return 0 on nil rate",
			chainID:  ethChainID,
			coinType: coin.CoinType_Gas,
			asset:    "",
			amount:   sdk.NewUint(100),
			gasAssetRates: func() map[int64]types.AssetRate {
				// set rate to nil
				nilAssetRateMap, _ := types.BuildAssetRateMapFromList(assetRateList)
				nilRate := nilAssetRateMap[ethChainID]
				nilRate.Rate = sdk.Dec{}
				nilAssetRateMap[ethChainID] = nilRate
				return nilAssetRateMap
			}(),
			erc20AssetRates: erc20AssetRateMap,
			expectedValue:   sdk.NewInt(0),
		},
		{
			name:          "should return 0 on rate <= 0",
			chainID:       ethChainID,
			coinType:      coin.CoinType_ERC20,
			asset:         assetUSDT,
			amount:        sdk.NewUint(100),
			gasAssetRates: gasAssetRateMap,
			erc20AssetRates: func() map[int64]map[string]types.AssetRate {
				// set rate to 0
				_, zeroAssetRateMap := types.BuildAssetRateMapFromList(assetRateList)
				zeroRate := zeroAssetRateMap[ethChainID][strings.ToLower(assetUSDT)]
				zeroRate.Rate = sdk.NewDec(0)
				zeroAssetRateMap[ethChainID][strings.ToLower(assetUSDT)] = zeroRate
				return zeroAssetRateMap
			}(),
			expectedValue: sdk.NewInt(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create cctx with given input
			cctx := sample.CrossChainTx(t, fmt.Sprintf("%d-%d", tt.chainID, 1))
			cctx.InboundParams.CoinType = tt.coinType
			cctx.InboundParams.Asset = tt.asset
			cctx.GetCurrentOutboundParam().Amount = tt.amount

			// convert cctx value
			value := types.ConvertCctxValueToAzeta(tt.chainID, cctx, tt.gasAssetRates, tt.erc20AssetRates)
			require.Equal(t, tt.expectedValue, value)
		})
	}
}

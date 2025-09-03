package sample

import (
	"testing"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/x/fungible/types"
)

func ForeignCoins(t *testing.T, address string) types.ForeignCoins {
	r := newRandFromStringSeed(t, address)

	return types.ForeignCoins{
		Zrc20ContractAddress: address,
		Asset:                EthAddress().String(),
		ForeignChainId:       r.Int63(),
		Decimals:             uint32(r.Uint64()),
		Name:                 StringRandom(r, 32),
		Symbol:               StringRandom(r, 32),
		CoinType:             coin.CoinType_ERC20,
		GasLimit:             r.Uint64(),
		LiquidityCap:         UintInRange(0, 10000000000),
	}
}

func ForeignCoinList(t *testing.T, zrc20ETH, zrc20BTC, zrc20ERC20, erc20Asset string) []types.ForeignCoins {
	// eth and btc chain id
	ethChainID := chains.GoerliLocalnet.ChainId
	btcChainID := chains.BitcoinRegtest.ChainId

	// add zrc20 ETH
	fcGas := ForeignCoins(t, zrc20ETH)
	fcGas.Asset = ""
	fcGas.ForeignChainId = ethChainID
	fcGas.Decimals = 18
	fcGas.CoinType = coin.CoinType_Gas

	// add zrc20 BTC
	fcBTC := ForeignCoins(t, zrc20BTC)
	fcBTC.Asset = ""
	fcBTC.ForeignChainId = btcChainID
	fcBTC.Decimals = 8
	fcBTC.CoinType = coin.CoinType_Gas

	// add zrc20 ERC20
	fcERC20 := ForeignCoins(t, zrc20ERC20)
	fcERC20.Asset = erc20Asset
	fcERC20.ForeignChainId = ethChainID
	fcERC20.Decimals = 6
	fcERC20.CoinType = coin.CoinType_ERC20

	return []types.ForeignCoins{fcGas, fcBTC, fcERC20}
}

func SystemContract() *types.SystemContract {
	return &types.SystemContract{
		SystemContract:  EthAddress().String(),
		ConnectorZevm:   EthAddress().String(),
		Gateway:         EthAddress().String(),
		GatewayGasLimit: types.DefaultGatewayGasLimit,
	}
}

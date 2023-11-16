package common

import (
	"github.com/btcsuite/btcd/chaincfg"
)

var (
	BitcoinMainnetParams = &chaincfg.MainNetParams
	BitcoinRegnetParams  = &chaincfg.RegressionNetParams
	BitcoinTestnetParams = &chaincfg.TestNet3Params
)

// BitcoinNetParamsFromChainID returns the bitcoin net params to be used from the chain id
func BitcoinNetParamsFromChainID(chainID int64) *chaincfg.Params {
	switch chainID {
	case BtcRegtestChain().ChainId:
		return BitcoinRegnetParams
	case BtcMainnetChain().ChainId:
		return BitcoinMainnetParams
	case BtcTestNetChain().ChainId:
		return BitcoinTestnetParams
	default:
		return BitcoinRegnetParams
	}
}

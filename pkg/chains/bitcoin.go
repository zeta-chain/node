package chains

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
)

var (
	BitcoinMainnetParams = &chaincfg.MainNetParams
	BitcoinRegnetParams  = &chaincfg.RegressionNetParams
	BitcoinTestnetParams = &chaincfg.TestNet3Params
)

// BitcoinNetParamsFromChainID returns the bitcoin net params to be used from the chain id
func BitcoinNetParamsFromChainID(chainID int64) (*chaincfg.Params, error) {
	switch chainID {
	case BtcRegtestChain.ChainId:
		return BitcoinRegnetParams, nil
	case BtcMainnetChain.ChainId:
		return BitcoinMainnetParams, nil
	case BtcTestNetChain.ChainId:
		return BitcoinTestnetParams, nil
	default:
		return nil, fmt.Errorf("no Bitcoin net params for chain ID: %d", chainID)
	}
}

// BitcoinChainIDFromNetworkName returns the chain id for the given bitcoin network name
func BitcoinChainIDFromNetworkName(name string) (int64, error) {
	switch name {
	case BitcoinRegnetParams.Name:
		return BtcRegtestChain.ChainId, nil
	case BitcoinMainnetParams.Name:
		return BtcMainnetChain.ChainId, nil
	case BitcoinTestnetParams.Name:
		return BtcTestNetChain.ChainId, nil
	default:
		return 0, fmt.Errorf("invalid Bitcoin network name: %s", name)
	}
}

// IsBitcoinRegnet returns true if the chain id is for the regnet
func IsBitcoinRegnet(chainID int64) bool {
	return chainID == BtcRegtestChain.ChainId
}

// IsBitcoinMainnet returns true if the chain id is for the mainnet
func IsBitcoinMainnet(chainID int64) bool {
	return chainID == BtcMainnetChain.ChainId
}

package chains

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
)

var (
	// chainIDToNetworkParams maps the Bitcoin chain ID to the network parameters
	chainIDToNetworkParams = map[int64]*chaincfg.Params{
		BitcoinRegtest.ChainId:       &chaincfg.RegressionNetParams,
		BitcoinMainnet.ChainId:       &chaincfg.MainNetParams,
		BitcoinTestnet.ChainId:       &chaincfg.TestNet3Params,
		BitcoinSignetTestnet.ChainId: &chaincfg.SigNetParams,
		BitcoinTestnet4.ChainId:      &TestNet4Params,
	}

	// networkNameToChainID maps the Bitcoin network name to the chain ID
	networkNameToChainID = map[string]int64{
		chaincfg.RegressionNetParams.Name: BitcoinRegtest.ChainId,
		chaincfg.MainNetParams.Name:       BitcoinMainnet.ChainId,
		chaincfg.TestNet3Params.Name:      BitcoinTestnet.ChainId,
		chaincfg.SigNetParams.Name:        BitcoinSignetTestnet.ChainId,
		TestNet4Params.Name:               BitcoinTestnet4.ChainId,
	}
)

// BitcoinNetParamsFromChainID returns the bitcoin net params to be used from the chain id
func BitcoinNetParamsFromChainID(chainID int64) (*chaincfg.Params, error) {
	if params, found := chainIDToNetworkParams[chainID]; found {
		return params, nil
	}
	return nil, fmt.Errorf("no Bitcoin network params for chain ID: %d", chainID)
}

// BitcoinChainIDFromNetworkName returns the chain id for the given bitcoin network name
func BitcoinChainIDFromNetworkName(name string) (int64, error) {
	if chainID, found := networkNameToChainID[name]; found {
		return chainID, nil
	}
	return 0, fmt.Errorf("invalid Bitcoin network name: %s", name)
}

// IsBitcoinRegnet returns true if the chain id is for the regnet
func IsBitcoinRegnet(chainID int64) bool {
	return chainID == BitcoinRegtest.ChainId
}

// IsBitcoinMainnet returns true if the chain id is for the mainnet
func IsBitcoinMainnet(chainID int64) bool {
	return chainID == BitcoinMainnet.ChainId
}

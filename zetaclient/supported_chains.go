package zetaclient

import (
	"github.com/zeta-chain/zetacore/common"
)

// Modify to update this from the core later
func GetSupportedChains() []*common.Chain {
	return common.DefaultChainsList()
}

func GetChainIDFromChainName(chainName common.ChainName) int64 {
	chains := GetSupportedChains()
	for _, chain := range chains {
		if chainName == chain.ChainName {
			return chain.ChainId
		}
	}
	return -1
}
func GetChainFromChainName(chainName common.ChainName) *common.Chain {
	chains := GetSupportedChains()
	for _, chain := range chains {
		if chainName == chain.ChainName {
			return chain
		}
	}
	return nil
}

func GetChainFromChainID(chainID int64) *common.Chain {
	chains := GetSupportedChains()
	for _, chain := range chains {
		if chainID == chain.ChainId {
			return chain
		}
	}
	return nil
}

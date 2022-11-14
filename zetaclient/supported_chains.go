package zetaclient

import zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"

// Modify to update this from the core later
func GetSupportedChains() []*zetaObserverTypes.Chain {
	return zetaObserverTypes.DefaultChainsList()
}

func GetChainIdFromChainName(chainName zetaObserverTypes.ChainName) int64 {
	chains := GetSupportedChains()
	for _, chain := range chains {
		if chainName == chain.ChainName {
			return chain.ChainId
		}
	}
	return -1
}
func GetChainFromChainName(chainName zetaObserverTypes.ChainName) *zetaObserverTypes.Chain {
	chains := GetSupportedChains()
	for _, chain := range chains {
		if chainName == chain.ChainName {
			return chain
		}
	}
	return nil
}

func GetChainNameFromChainId(chainId int64) zetaObserverTypes.ChainName {
	chains := GetSupportedChains()
	for _, chain := range chains {
		if chainId == chain.ChainId {
			return chain.ChainName
		}
	}
	return zetaObserverTypes.ChainName_Empty
}

func GetChainFromChainId(chainId int64) *zetaObserverTypes.Chain {
	chains := GetSupportedChains()
	for _, chain := range chains {
		if chainId == chain.ChainId {
			return chain
		}
	}
	return nil
}

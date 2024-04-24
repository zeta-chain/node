package chains

import "fmt"

// Mainnet chains
func ZetaChainMainnet() Chain {
	return Chain{
		ChainName:         ChainName_zeta_mainnet,
		ChainId:           7000,
		Network:           Network_ZETA,
		NetworkType:       NetworkType_MAINNET,
		Vm:                Vm_EVM,
		Consensus:         Consensus_Tendermint,
		IsExternal:        false,
		IsHeaderSupported: false,
	}
}
func EthChain() Chain {
	return Chain{
		ChainName:         ChainName_eth_mainnet,
		ChainId:           1,
		Network:           Network_ETH,
		NetworkType:       NetworkType_MAINNET,
		Vm:                Vm_EVM,
		Consensus:         Consensus_Ethereum,
		IsExternal:        true,
		IsHeaderSupported: true,
	}
}

func BscMainnetChain() Chain {
	return Chain{
		ChainName:         ChainName_bsc_mainnet,
		ChainId:           56,
		Network:           Network_BSC,
		NetworkType:       NetworkType_MAINNET,
		Vm:                Vm_EVM,
		Consensus:         Consensus_Ethereum,
		IsExternal:        true,
		IsHeaderSupported: true,
	}
}

func BtcMainnetChain() Chain {
	return Chain{
		ChainName:         ChainName_btc_mainnet,
		ChainId:           8332,
		Network:           Network_BTC,
		NetworkType:       NetworkType_MAINNET,
		Vm:                Vm_NO_VM,
		Consensus:         Consensus_Bitcoin,
		IsExternal:        true,
		IsHeaderSupported: false,
	}
}

func PolygonChain() Chain {
	return Chain{
		ChainName:         ChainName_polygon_mainnet,
		ChainId:           137,
		Network:           Network_POLYGON,
		NetworkType:       NetworkType_MAINNET,
		Vm:                Vm_EVM,
		Consensus:         Consensus_Ethereum,
		IsExternal:        true,
		IsHeaderSupported: false,
	}
}

// Testnet chains

func ZetaTestnetChain() Chain {
	return Chain{
		ChainName:         ChainName_zeta_testnet,
		ChainId:           7001,
		Network:           Network_ZETA,
		NetworkType:       NetworkType_TESTNET,
		Vm:                Vm_EVM,
		Consensus:         Consensus_Tendermint,
		IsExternal:        false,
		IsHeaderSupported: false,
	}
}

func SepoliaChain() Chain {
	return Chain{
		ChainName:         ChainName_sepolia_testnet,
		ChainId:           11155111,
		Network:           Network_ETH,
		NetworkType:       NetworkType_TESTNET,
		Vm:                Vm_EVM,
		Consensus:         Consensus_Ethereum,
		IsExternal:        true,
		IsHeaderSupported: true,
	}
}

// GoerliChain Deprecated
func GoerliChain() Chain {
	return Chain{
		ChainName:         ChainName_goerli_testnet,
		ChainId:           5,
		Network:           Network_ETH,
		NetworkType:       NetworkType_TESTNET,
		Vm:                Vm_EVM,
		Consensus:         Consensus_Ethereum,
		IsExternal:        true,
		IsHeaderSupported: true,
	}
}

func BscTestnetChain() Chain {
	return Chain{
		ChainName:         ChainName_bsc_testnet,
		ChainId:           97,
		Network:           Network_BSC,
		NetworkType:       NetworkType_TESTNET,
		Vm:                Vm_EVM,
		Consensus:         Consensus_Ethereum,
		IsExternal:        true,
		IsHeaderSupported: true,
	}
}

func BtcTestNetChain() Chain {
	return Chain{
		ChainName:         ChainName_btc_testnet,
		ChainId:           18332,
		Network:           Network_BTC,
		NetworkType:       NetworkType_TESTNET,
		Vm:                Vm_NO_VM,
		Consensus:         Consensus_Bitcoin,
		IsExternal:        true,
		IsHeaderSupported: false,
	}
}

// MumbaiChain Deprecated
func MumbaiChain() Chain {
	return Chain{
		ChainName:         ChainName_mumbai_testnet,
		ChainId:           80001,
		Network:           Network_POLYGON,
		NetworkType:       NetworkType_TESTNET,
		Vm:                Vm_EVM,
		Consensus:         Consensus_Ethereum,
		IsExternal:        true,
		IsHeaderSupported: false,
	}
}

func AmoyChain() Chain {
	return Chain{
		ChainName:         ChainName_amoy_testnet,
		ChainId:           80002,
		Network:           Network_POLYGON,
		NetworkType:       NetworkType_TESTNET,
		Vm:                Vm_EVM,
		Consensus:         Consensus_Ethereum,
		IsExternal:        true,
		IsHeaderSupported: false,
	}
}

// Devnet chains
func ZetaMocknetChain() Chain {
	return Chain{
		ChainName:         ChainName_zeta_mainnet,
		ChainId:           70000,
		Network:           Network_ZETA,
		NetworkType:       NetworkType_DEVNET,
		Vm:                Vm_EVM,
		Consensus:         Consensus_Tendermint,
		IsExternal:        false,
		IsHeaderSupported: false,
	}
}

// Privnet chains

func ZetaPrivnetChain() Chain {
	return Chain{
		ChainName:         ChainName_zeta_mainnet,
		ChainId:           101,
		Network:           Network_ZETA,
		NetworkType:       NetworkType_PRIVNET,
		Vm:                Vm_EVM,
		Consensus:         Consensus_Tendermint,
		IsExternal:        false,
		IsHeaderSupported: false,
	}
}
func BtcRegtestChain() Chain {
	return Chain{
		ChainName:         ChainName_btc_regtest,
		ChainId:           18444,
		Network:           Network_BTC,
		NetworkType:       NetworkType_PRIVNET,
		Vm:                Vm_NO_VM,
		Consensus:         Consensus_Bitcoin,
		IsExternal:        true,
		IsHeaderSupported: false,
	}
}

func GoerliLocalnetChain() Chain {
	return Chain{
		ChainName:         ChainName_goerli_localnet,
		ChainId:           1337,
		Network:           Network_ETH,
		NetworkType:       NetworkType_PRIVNET,
		Vm:                Vm_EVM,
		Consensus:         Consensus_Ethereum,
		IsExternal:        true,
		IsHeaderSupported: true,
	}
}

func BtcDustOffset() int64 {
	return 2000
}

// DefaultChainsList returns a list of default chains
func DefaultChainsList() []*Chain {
	return chainListPointers([]Chain{
		BtcMainnetChain(),
		BscMainnetChain(),
		EthChain(),
		BtcTestNetChain(),
		MumbaiChain(),
		AmoyChain(),
		BscTestnetChain(),
		GoerliChain(),
		SepoliaChain(),
		BtcRegtestChain(),
		GoerliLocalnetChain(),
		ZetaChainMainnet(),
		ZetaTestnetChain(),
		ZetaMocknetChain(),
		ZetaPrivnetChain(),
		PolygonChain(),
	})
}

// ChainListByNetworkType returns a list of chains by network type
func ChainListByNetworkType(networkType NetworkType) []*Chain {
	var mainNetList []*Chain
	for _, chain := range DefaultChainsList() {
		if chain.NetworkType == networkType {
			mainNetList = append(mainNetList, chain)
		}
	}
	return mainNetList
}

// ChainListByNetwork returns a list of chains by network
func ChainListByNetwork(network Network) []*Chain {
	var chainList []*Chain
	for _, chain := range DefaultChainsList() {
		if chain.Network == network {
			chainList = append(chainList, chain)
		}
	}
	return chainList
}

// ExternalChainList returns a list chains that are not Zeta
func ExternalChainList() []*Chain {
	var chainList []*Chain
	for _, chain := range DefaultChainsList() {
		if chain.IsExternal {
			chainList = append(chainList, chain)
		}
	}
	return chainList
}

// ChainListByConsensus returns a list of chains by consensus
func ChainListByConsensus(consensus Consensus) []*Chain {
	var chainList []*Chain
	for _, chain := range DefaultChainsList() {
		if chain.Consensus == consensus {
			chainList = append(chainList, chain)
		}
	}
	return chainList

}

// ChainListForHeaderSupport returns a list of chains that support headers
func ChainListForHeaderSupport() []*Chain {
	var chainList []*Chain
	for _, chain := range DefaultChainsList() {
		if chain.IsHeaderSupported {
			chainList = append(chainList, chain)
		}
	}
	return chainList
}

// ZetaChainFromChainID returns a ZetaChain chain object  from a Cosmos chain ID
func ZetaChainFromChainID(chainID string) (Chain, error) {
	ethChainID, err := CosmosToEthChainID(chainID)
	if err != nil {
		return Chain{}, err
	}

	switch ethChainID {
	case ZetaPrivnetChain().ChainId:
		return ZetaPrivnetChain(), nil
	case ZetaChainMainnet().ChainId:
		return ZetaChainMainnet(), nil
	case ZetaTestnetChain().ChainId:
		return ZetaTestnetChain(), nil
	case ZetaMocknetChain().ChainId:
		return ZetaMocknetChain(), nil
	default:
		return Chain{}, fmt.Errorf("chain %d not found", ethChainID)
	}
}

// chainListPointers returns a list of chain pointers
func chainListPointers(chains []Chain) []*Chain {
	var c []*Chain
	for i := 0; i < len(chains); i++ {
		c = append(c, &chains[i])
	}
	return c
}

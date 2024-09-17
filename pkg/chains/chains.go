package chains

import "fmt"

var (
	/**
	* Mainnet chains
	 */

	// ZetaChainMainnet is the mainnet chain for Zeta
	ZetaChainMainnet = Chain{
		ChainName:   ChainName_zeta_mainnet,
		ChainId:     7000,
		Network:     Network_zeta,
		NetworkType: NetworkType_mainnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_tendermint,
		IsExternal:  false,
		CctxGateway: CCTXGateway_zevm,
		Name:        "zeta_mainnet",
	}

	// Ethereum is Ethereum mainnet
	Ethereum = Chain{
		ChainName:   ChainName_eth_mainnet,
		ChainId:     1,
		Network:     Network_eth,
		NetworkType: NetworkType_mainnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_ethereum,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "eth_mainnet",
	}

	// BscMainnet is Binance Smart Chain mainnet
	BscMainnet = Chain{
		ChainName:   ChainName_bsc_mainnet,
		ChainId:     56,
		Network:     Network_bsc,
		NetworkType: NetworkType_mainnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_ethereum,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "bsc_mainnet",
	}

	// BitcoinMainnet is Bitcoin mainnet
	BitcoinMainnet = Chain{
		ChainName:   ChainName_btc_mainnet,
		ChainId:     8332,
		Network:     Network_btc,
		NetworkType: NetworkType_mainnet,
		Vm:          Vm_no_vm,
		Consensus:   Consensus_bitcoin,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "btc_mainnet",
	}

	// Polygon is Polygon mainnet
	Polygon = Chain{
		ChainName:   ChainName_polygon_mainnet,
		ChainId:     137,
		Network:     Network_polygon,
		NetworkType: NetworkType_mainnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_ethereum,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "polygon_mainnet",
	}

	// OptimismMainnet is Optimism mainnet
	OptimismMainnet = Chain{
		ChainName:   ChainName_optimism_mainnet,
		ChainId:     10,
		Network:     Network_optimism,
		NetworkType: NetworkType_mainnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_op_stack,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "optimism_mainnet",
	}

	// BaseMainnet is Base mainnet
	BaseMainnet = Chain{
		ChainName:   ChainName_base_mainnet,
		ChainId:     8453,
		Network:     Network_base,
		NetworkType: NetworkType_mainnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_op_stack,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "base_mainnet",
	}

	// SolanaMainnet is Solana mainnet
	// TODO: define final chain ID
	// https://github.com/zeta-chain/node/issues/2421
	SolanaMainnet = Chain{
		ChainName:   ChainName_solana_mainnet,
		ChainId:     900,
		Network:     Network_solana,
		NetworkType: NetworkType_mainnet,
		Vm:          Vm_svm,
		Consensus:   Consensus_solana_consensus,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "solana_mainnet",
	}

	/**
	* Testnet chains
	 */

	// ZetaChainTestnet is the testnet chain for Zeta
	ZetaChainTestnet = Chain{
		ChainName:   ChainName_zeta_testnet,
		ChainId:     7001,
		Network:     Network_zeta,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_tendermint,
		IsExternal:  false,
		CctxGateway: CCTXGateway_zevm,
		Name:        "zeta_testnet",
	}

	// Sepolia is Ethereum sepolia testnet
	Sepolia = Chain{
		ChainName:   ChainName_sepolia_testnet,
		ChainId:     11155111,
		Network:     Network_eth,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_ethereum,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "sepolia_testnet",
	}

	// BscTestnet is Binance Smart Chain testnet
	BscTestnet = Chain{
		ChainName:   ChainName_bsc_testnet,
		ChainId:     97,
		Network:     Network_bsc,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_ethereum,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "bsc_testnet",
	}

	// BitcoinTestnet is Bitcoin testnet3
	BitcoinTestnet = Chain{
		ChainName:   ChainName_btc_testnet,
		ChainId:     18332,
		Network:     Network_btc,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_no_vm,
		Consensus:   Consensus_bitcoin,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "btc_testnet",
	}

	BitcoinSignetTestnet = Chain{
		ChainName:   ChainName_btc_signet_testnet,
		ChainId:     18334,
		Network:     Network_btc,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_no_vm,
		Consensus:   Consensus_bitcoin,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "btc_signet_testnet",
	}

	// Amoy is Polygon amoy testnet
	Amoy = Chain{
		ChainName:   ChainName_amoy_testnet,
		ChainId:     80002,
		Network:     Network_polygon,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_ethereum,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "amoy_testnet",
	}

	// OptimismSepolia is Optimism sepolia testnet
	OptimismSepolia = Chain{
		ChainName:   ChainName_optimism_sepolia,
		ChainId:     11155420,
		Network:     Network_optimism,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_op_stack,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "optimism_sepolia",
	}

	// BaseSepolia is Base sepolia testnet
	BaseSepolia = Chain{
		ChainName:   ChainName_base_sepolia,
		ChainId:     84532,
		Network:     Network_base,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_op_stack,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "base_sepolia",
	}

	// SolanaDevnet is Solana devnet
	// NOTE: Solana devnet refers to Solana testnet in our terminology
	// Solana uses devnet denomitation for network for development
	// TODO: define final chain ID
	// https://github.com/zeta-chain/node/issues/2421
	SolanaDevnet = Chain{
		ChainName:   ChainName_solana_devnet,
		ChainId:     901,
		Network:     Network_solana,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_svm,
		Consensus:   Consensus_solana_consensus,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "solana_devnet",
	}

	/**
	* Devnet chains
	 */

	// ZetaChainDevnet is the devnet chain for Zeta
	// used as live testing environment
	ZetaChainDevnet = Chain{
		ChainName:   ChainName_zeta_mainnet,
		ChainId:     70000,
		Network:     Network_zeta,
		NetworkType: NetworkType_devnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_tendermint,
		IsExternal:  false,
		CctxGateway: CCTXGateway_zevm,
		Name:        "zeta_mainnet",
	}

	/**
	* Privnet chains
	 */

	// ZetaChainPrivnet is the privnet chain for Zeta (localnet)
	ZetaChainPrivnet = Chain{
		ChainName:   ChainName_zeta_mainnet,
		ChainId:     101,
		Network:     Network_zeta,
		NetworkType: NetworkType_privnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_tendermint,
		IsExternal:  false,
		CctxGateway: CCTXGateway_zevm,
		Name:        "zeta_mainnet",
	}

	// BitcoinRegtest is Bitcoin regtest (localnet)
	BitcoinRegtest = Chain{
		ChainName:   ChainName_btc_regtest,
		ChainId:     18444,
		Network:     Network_btc,
		NetworkType: NetworkType_privnet,
		Vm:          Vm_no_vm,
		Consensus:   Consensus_bitcoin,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "btc_regtest",
	}

	// GoerliLocalnet is Ethereum local goerli (localnet)
	GoerliLocalnet = Chain{
		ChainName:   ChainName_goerli_localnet,
		ChainId:     1337,
		Network:     Network_eth,
		NetworkType: NetworkType_privnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_ethereum,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "goerli_localnet",
	}

	// SolanaLocalnet is Solana localnet
	// TODO: define final chain ID
	// https://github.com/zeta-chain/node/issues/2421
	SolanaLocalnet = Chain{
		ChainName:   ChainName_solana_localnet,
		ChainId:     902,
		Network:     Network_solana,
		NetworkType: NetworkType_privnet,
		Vm:          Vm_svm,
		Consensus:   Consensus_solana_consensus,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "solana_localnet",
	}

	/**
	* Deprecated chains
	 */

	// Goerli is Ethereum goerli testnet (deprecated for sepolia)
	Goerli = Chain{
		ChainName:   ChainName_goerli_testnet,
		ChainId:     5,
		Network:     Network_eth,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_ethereum,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "goerli_testnet",
	}

	// Mumbai is Polygon mumbai testnet (deprecated for amoy)
	Mumbai = Chain{
		ChainName:   ChainName_mumbai_testnet,
		ChainId:     80001,
		Network:     Network_polygon,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_ethereum,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "mumbai_testnet",
	}
)

// ErrNotZetaChain is the error for chain not being a ZetaChain chain
var ErrNotZetaChain = fmt.Errorf("chain is not a ZetaChain chain")

// BtcNonceMarkOffset is the offset satoshi amount to calculate the nonce mark output
func BtcNonceMarkOffset() int64 {
	return 2000
}

// DefaultChainsList returns a list of default chains
func DefaultChainsList() []Chain {
	return []Chain{
		BitcoinMainnet,
		BscMainnet,
		Ethereum,
		BitcoinTestnet,
		BitcoinSignetTestnet,
		Mumbai,
		Amoy,
		BscTestnet,
		Goerli,
		Sepolia,
		BitcoinRegtest,
		GoerliLocalnet,
		ZetaChainMainnet,
		ZetaChainTestnet,
		ZetaChainDevnet,
		ZetaChainPrivnet,
		Polygon,
		OptimismMainnet,
		OptimismSepolia,
		BaseMainnet,
		BaseSepolia,
		SolanaMainnet,
		SolanaDevnet,
		SolanaLocalnet,
	}
}

// ChainListByNetworkType returns a list of chains by network type
func ChainListByNetworkType(networkType NetworkType, additionalChains []Chain) []Chain {
	var chainList []Chain
	for _, chain := range CombineDefaultChainsList(additionalChains) {
		if chain.NetworkType == networkType {
			chainList = append(chainList, chain)
		}
	}
	return chainList
}

// ChainListByNetwork returns a list of chains by network
func ChainListByNetwork(network Network, additionalChains []Chain) []Chain {
	var chainList []Chain
	for _, chain := range CombineDefaultChainsList(additionalChains) {
		if chain.Network == network {
			chainList = append(chainList, chain)
		}
	}
	return chainList
}

// ExternalChainList returns a list chains that are not Zeta
func ExternalChainList(additionalChains []Chain) []Chain {
	var chainList []Chain
	for _, chain := range CombineDefaultChainsList(additionalChains) {
		if chain.IsExternal {
			chainList = append(chainList, chain)
		}
	}
	return chainList
}

// ChainListByConsensus returns a list of chains by consensus
func ChainListByConsensus(consensus Consensus, additionalChains []Chain) []Chain {
	var chainList []Chain
	for _, chain := range CombineDefaultChainsList(additionalChains) {
		if chain.Consensus == consensus {
			chainList = append(chainList, chain)
		}
	}
	return chainList
}

func ChainListByGateway(gateway CCTXGateway, additionalChains []Chain) []Chain {
	var chainList []Chain
	for _, chain := range CombineDefaultChainsList(additionalChains) {
		if chain.CctxGateway == gateway {
			chainList = append(chainList, chain)
		}
	}
	return chainList
}

// ChainListForHeaderSupport returns a list of chains that support headers
func ChainListForHeaderSupport(additionalChains []Chain) []Chain {
	var chainList []Chain
	for _, chain := range CombineDefaultChainsList(additionalChains) {
		if chain.Consensus == Consensus_ethereum || chain.Consensus == Consensus_bitcoin {
			chainList = append(chainList, chain)
		}
	}
	return chainList
}

// ZetaChainFromCosmosChainID returns a ZetaChain chain object from a Cosmos chain ID
func ZetaChainFromCosmosChainID(chainID string) (Chain, error) {
	ethChainID, err := CosmosToEthChainID(chainID)
	if err != nil {
		return Chain{}, err
	}

	return ZetaChainFromChainID(ethChainID)
}

// ZetaChainFromChainID returns a ZetaChain chain object from a chain ID
func ZetaChainFromChainID(chainID int64) (Chain, error) {
	switch chainID {
	case ZetaChainPrivnet.ChainId:
		return ZetaChainPrivnet, nil
	case ZetaChainMainnet.ChainId:
		return ZetaChainMainnet, nil
	case ZetaChainTestnet.ChainId:
		return ZetaChainTestnet, nil
	case ZetaChainDevnet.ChainId:
		return ZetaChainDevnet, nil
	default:
		return Chain{}, ErrNotZetaChain
	}
}

// CombineDefaultChainsList combines the default chains list with a list of chains
// duplicated chain ID are overwritten by the second list
func CombineDefaultChainsList(chains []Chain) []Chain {
	return CombineChainList(DefaultChainsList(), chains)
}

// CombineChainList combines a list of chains with a list of chains
// duplicated chain ID are overwritten by the second list
func CombineChainList(base []Chain, additional []Chain) []Chain {
	combined := make([]Chain, 0, len(base)+len(additional))
	combined = append(combined, base...)

	// map chain ID in combined to index in the list
	chainIDIndexMap := make(map[int64]int)
	for i, chain := range combined {
		chainIDIndexMap[chain.ChainId] = i
	}

	// add chains2 to combined
	// if chain ID already exists in chains1, overwrite it
	for _, chain := range additional {
		if index, ok := chainIDIndexMap[chain.ChainId]; ok {
			combined[index] = chain
		} else {
			combined = append(combined, chain)
		}
	}

	return combined
}

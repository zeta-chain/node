package chains

// ChainFilter is a function that filters chains based on some criteria
type ChainFilter func(c Chain) bool

// FilterExternalChains filters chains that are external
func FilterExternalChains(c Chain) bool {
	return c.IsExternal
}

// FilterGatewayObserver filters chains that are gateway observers
func FilterGatewayObserver(c Chain) bool {
	return c.CctxGateway == CCTXGateway_observers
}

// FilterConsensusEthereum filters chains that have the ethereum consensus
func FilterConsensusEthereum(c Chain) bool {
	return c.Consensus == Consensus_ethereum
}

// FilterConsensusBitcoin filters chains that have the bitcoin consensus
func FilterConsensusBitcoin(c Chain) bool { return c.Consensus == Consensus_bitcoin }

// FilterConsensusSolana filters chains that have the solana consensus
func FilterConsensusSolana(c Chain) bool { return c.Consensus == Consensus_solana_consensus }

// FilterChains applies a list of filters to a list of chains
func FilterChains(chainList []Chain, filters ...ChainFilter) []Chain {
	// Apply each filter to the list of supported chains
	for _, filter := range filters {
		var filteredChains []Chain
		for _, chain := range chainList {
			if filter(chain) {
				filteredChains = append(filteredChains, chain)
			}
		}
		chainList = filteredChains
	}

	// Return the filtered list of chains
	return chainList
}

func CombineFilterChains(chainLists ...[]Chain) []Chain {
	chainMap := make(map[Chain]bool)
	var combinedChains []Chain

	// Add chains from each slice to remove duplicates
	for _, chains := range chainLists {
		for _, chain := range chains {
			if !chainMap[chain] {
				chainMap[chain] = true
				combinedChains = append(combinedChains, chain)
			}
		}
	}

	return combinedChains
}

package chains

// ChainFilter is a function that filters chains based on some criteria
type ChainFilter func(c Chain) bool

// FilterExternalChains filters chains that are external
func FilterExternalChains(c Chain) bool {
	return c.IsExternal
}

// FilterByGateway filters chains by gateway
func FilterByGateway(gw CCTXGateway) ChainFilter {
	return func(chain Chain) bool { return chain.CctxGateway == gw }
}

// FilterByConsensus filters chains by consensus type
func FilterByConsensus(cs Consensus) ChainFilter {
	return func(chain Chain) bool { return chain.Consensus == cs }
}

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

// CombineFilterChains combines multiple lists of chains into a single list
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

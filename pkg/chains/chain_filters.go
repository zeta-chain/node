package chains

type ChainFilter func(c Chain) bool

func FilterExternalChains(c Chain) bool {
	return c.IsExternal
}

func FilterGatewayObserver(c Chain) bool {
	return c.CctxGateway == CCTXGateway_observers
}

func FilterConsensusEthereum(c Chain) bool {
	return c.Consensus == Consensus_ethereum
}

func FilterConsensusBitcoin(c Chain) bool { return c.Consensus == Consensus_bitcoin }

func FilterConsensusSolana(c Chain) bool { return c.Consensus == Consensus_solana_consensus }

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

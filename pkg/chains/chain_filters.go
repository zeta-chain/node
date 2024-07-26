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

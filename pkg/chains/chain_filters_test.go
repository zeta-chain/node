package chains_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
)

func TestFilterChains(t *testing.T) {
	tt := []struct {
		name     string
		filters  []chains.ChainFilter
		expected func() []chains.Chain
	}{
		{
			name:    "Filter external chains",
			filters: []chains.ChainFilter{chains.FilterExternalChains},
			expected: func() []chains.Chain {
				return chains.ExternalChainList([]chains.Chain{})
			},
		},
		{
			name:    "Filter gateway observer chains",
			filters: []chains.ChainFilter{chains.FilterGatewayObserver},
			expected: func() []chains.Chain {
				return chains.ChainListByGateway(chains.CCTXGateway_observers, []chains.Chain{})
			},
		},
		{
			name:    "Filter consensus ethereum chains",
			filters: []chains.ChainFilter{chains.FilterConsensusEthereum},
			expected: func() []chains.Chain {
				return chains.ChainListByConsensus(chains.Consensus_ethereum, []chains.Chain{})
			},
		},
		{
			name:    "Filter consensus bitcoin chains",
			filters: []chains.ChainFilter{chains.FilterConsensusBitcoin},
			expected: func() []chains.Chain {
				return chains.ChainListByConsensus(chains.Consensus_bitcoin, []chains.Chain{})
			},
		},
		{
			name:    "Filter consensus solana chains",
			filters: []chains.ChainFilter{chains.FilterConsensusSolana},
			expected: func() []chains.Chain {
				return chains.ChainListByConsensus(chains.Consensus_solana_consensus, []chains.Chain{})
			},
		},
		{
			name:    "Apply multiple filters external chains and gateway observer",
			filters: []chains.ChainFilter{chains.FilterExternalChains, chains.FilterGatewayObserver},
			expected: func() []chains.Chain {
				externalChains := chains.ExternalChainList([]chains.Chain{})
				var gatewayObserverChains []chains.Chain
				for _, chain := range externalChains {
					if chain.CctxGateway == chains.CCTXGateway_observers {
						gatewayObserverChains = append(gatewayObserverChains, chain)
					}
				}
				return gatewayObserverChains
			},
		},
		{
			name: "Apply multiple filters external chains with gateway observer and consensus ethereum",
			filters: []chains.ChainFilter{
				chains.FilterExternalChains,
				chains.FilterGatewayObserver,
				chains.FilterConsensusEthereum,
			},
			expected: func() []chains.Chain {
				externalChains := chains.ExternalChainList([]chains.Chain{})
				var filterMultipleChains []chains.Chain
				for _, chain := range externalChains {
					if chain.CctxGateway == chains.CCTXGateway_observers &&
						chain.Consensus == chains.Consensus_ethereum {
						filterMultipleChains = append(filterMultipleChains, chain)
					}
				}
				return filterMultipleChains
			},
		},
		{
			name: "Apply multiple filters external chains with gateway observer and consensus bitcoin",
			filters: []chains.ChainFilter{
				chains.FilterExternalChains,
				chains.FilterGatewayObserver,
				chains.FilterConsensusBitcoin,
			},
			expected: func() []chains.Chain {
				externalChains := chains.ExternalChainList([]chains.Chain{})
				var filterMultipleChains []chains.Chain
				for _, chain := range externalChains {
					if chain.CctxGateway == chains.CCTXGateway_observers &&
						chain.Consensus == chains.Consensus_bitcoin {
						filterMultipleChains = append(filterMultipleChains, chain)
					}
				}
				return filterMultipleChains
			},
		},
		{
			name: "Test multiple filters in random order",
			filters: []chains.ChainFilter{
				chains.FilterGatewayObserver,
				chains.FilterConsensusEthereum,
				chains.FilterExternalChains,
			},
			expected: func() []chains.Chain {
				externalChains := chains.ExternalChainList([]chains.Chain{})
				var filterMultipleChains []chains.Chain
				for _, chain := range externalChains {
					if chain.CctxGateway == chains.CCTXGateway_observers &&
						chain.Consensus == chains.Consensus_ethereum {
						filterMultipleChains = append(filterMultipleChains, chain)
					}
				}
				return filterMultipleChains
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			chainList := chains.ExternalChainList([]chains.Chain{})
			filteredChains := chains.FilterChains(chainList, tc.filters...)
			require.ElementsMatch(t, tc.expected(), filteredChains)
		})
	}
}

func TestCombineFilterChains(t *testing.T) {
	tt := []struct {
		name       string
		chainLists func() [][]chains.Chain
		expected   func() []chains.Chain
	}{
		{
			name: "test support TSS migration filter",
			chainLists: func() [][]chains.Chain {
				return [][]chains.Chain{
					chains.FilterChains(
						chains.ExternalChainList([]chains.Chain{}),
						[]chains.ChainFilter{
							chains.FilterExternalChains,
							chains.FilterGatewayObserver,
							chains.FilterConsensusEthereum,
						}...),
					chains.FilterChains(
						chains.ExternalChainList([]chains.Chain{}),
						[]chains.ChainFilter{
							chains.FilterExternalChains,
							chains.FilterGatewayObserver,
							chains.FilterConsensusBitcoin,
						}...),
				}
			},
			expected: func() []chains.Chain {
				chainList := chains.ExternalChainList([]chains.Chain{})
				var filterMultipleChains []chains.Chain
				for _, chain := range chainList {
					if chain.CctxGateway == chains.CCTXGateway_observers &&
						(chain.Consensus == chains.Consensus_ethereum || chain.Consensus == chains.Consensus_bitcoin) {
						filterMultipleChains = append(filterMultipleChains, chain)
					}
				}
				return filterMultipleChains
			},
		},
		{
			name: "test support TSS migration filter with solana",
			chainLists: func() [][]chains.Chain {
				return [][]chains.Chain{
					chains.FilterChains(
						chains.ExternalChainList([]chains.Chain{}),
						[]chains.ChainFilter{
							chains.FilterExternalChains,
							chains.FilterGatewayObserver,
							chains.FilterConsensusEthereum,
						}...),
					chains.FilterChains(
						chains.ExternalChainList([]chains.Chain{}),
						[]chains.ChainFilter{
							chains.FilterExternalChains,
							chains.FilterGatewayObserver,
							chains.FilterConsensusBitcoin,
						}...),
					chains.FilterChains(
						chains.ExternalChainList([]chains.Chain{}),
						[]chains.ChainFilter{
							chains.FilterExternalChains,
							chains.FilterGatewayObserver,
							chains.FilterConsensusSolana,
						}...),
				}
			},
			expected: func() []chains.Chain {
				chainList := chains.ExternalChainList([]chains.Chain{})
				var filterMultipleChains []chains.Chain
				for _, chain := range chainList {
					if chain.CctxGateway == chains.CCTXGateway_observers &&
						(chain.Consensus == chains.Consensus_ethereum || chain.Consensus == chains.Consensus_bitcoin || chain.Consensus == chains.Consensus_solana_consensus) {
						filterMultipleChains = append(filterMultipleChains, chain)
					}
				}
				return filterMultipleChains
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			chainLists := tc.chainLists()
			combinedChains := chains.CombineFilterChains(chainLists...)
			require.ElementsMatch(t, tc.expected(), combinedChains)
		})
	}
}

package types

import "github.com/zeta-chain/zetacore/pkg/chains"

func DefaultBlockHeaderVerification() BlockHeaderVerification {
	return BlockHeaderVerification{
		EnabledChains: DefaultVerificationFlags(),
	}
}

// DefaultVerificationFlags returns the default verification flags.
// By default, everything disabled.
func DefaultVerificationFlags() []EnabledChain {
	return []EnabledChain{
		{
			ChainId: chains.EthChain.ChainId,
			Enabled: false,
		},
		{
			ChainId: chains.BscMainnetChain.ChainId,
			Enabled: false,
		},
		{
			ChainId: chains.SepoliaChain.ChainId,
			Enabled: false,
		},
		{
			ChainId: chains.BscTestnetChain.ChainId,
			Enabled: false,
		},
		{
			ChainId: chains.GoerliLocalnetChain.ChainId,
			Enabled: false,
		},
		{
			ChainId: chains.GoerliChain.ChainId,
			Enabled: false,
		},
	}
}

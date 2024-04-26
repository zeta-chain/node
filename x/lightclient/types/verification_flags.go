package types

import "github.com/zeta-chain/zetacore/pkg/chains"

// DefaultVerificationFlags returns the default verification flags.
// By default, everything disabled.
func DefaultVerificationFlags() []VerificationFlags {
	return []VerificationFlags{
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

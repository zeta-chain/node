package types

import "github.com/zeta-chain/node/pkg/chains"

func DefaultBlockHeaderVerification() BlockHeaderVerification {
	return BlockHeaderVerification{
		HeaderSupportedChains: DefaultHeaderSupportedChains(),
	}
}

// DefaultHeaderSupportedChains returns the default verification flags.
// By default, everything disabled.
func DefaultHeaderSupportedChains() []HeaderSupportedChain {
	return []HeaderSupportedChain{
		{
			ChainId: chains.Ethereum.ChainId,
			Enabled: false,
		},
		{
			ChainId: chains.BscMainnet.ChainId,
			Enabled: false,
		},
		{
			ChainId: chains.Sepolia.ChainId,
			Enabled: false,
		},
		{
			ChainId: chains.BscTestnet.ChainId,
			Enabled: false,
		},
		{
			ChainId: chains.GoerliLocalnet.ChainId,
			Enabled: false,
		},
		{
			ChainId: chains.Goerli.ChainId,
			Enabled: false,
		},
	}
}

package address_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/address"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

func appendChains(chainLists ...[]chains.Chain) []chains.Chain {
	var combined []chains.Chain
	for _, chains := range chainLists {
		combined = append(combined, chains...)
	}
	return combined
}

func TestValidateAddressForChain(t *testing.T) {
	additionalChains := []chains.Chain{}
	evmChains := []chains.Chain{
		// Ethereum
		chains.Ethereum,
		chains.Sepolia,
		chains.Goerli,
		chains.GoerliLocalnet,
		// Polygon
		chains.Polygon,
		chains.Amoy,
		chains.Mumbai,
		// BSC
		chains.BscMainnet,
		chains.BscTestnet,
	}
	zetaChains := []chains.Chain{
		chains.ZetaChainMainnet,
		chains.ZetaChainTestnet,
		chains.ZetaChainDevnet,
		chains.ZetaChainPrivnet,
	}
	baseChains := []chains.Chain{
		chains.BaseMainnet,
		chains.BaseSepolia,
	}

	solanaChains := []chains.Chain{
		chains.SolanaMainnet,
		chains.SolanaDevnet,
		chains.SolanaLocalnet,
	}

	opChains := []chains.Chain{
		chains.OptimismMainnet,
		chains.OptimismSepolia,
	}

	bitcoinChains := []chains.Chain{
		chains.BitcoinMainnet,
		chains.BitcoinTestnet,
		chains.BitcoinRegtest,
	}

	// Append all chains
	allChains := appendChains(evmChains, zetaChains, baseChains, solanaChains, opChains, bitcoinChains)
	require.ElementsMatch(t, allChains, chains.DefaultChainsList())

	// test for evm chain chains
	for _, chain := range evmChains {
		require.Error(t, address.ValidateAddressForChain("0x123", chain.ChainId, additionalChains))
		require.Error(t, address.ValidateAddressForChain("", chain.ChainId, additionalChains))
		require.Error(t, address.ValidateAddressForChain("%%%%", chain.ChainId, additionalChains))
		require.NoError(
			t,
			address.ValidateAddressForChain(sample.EthAddress().Hex(), chain.ChainId, additionalChains),
		)
		sample.EthAddress()
	}

	// test for zeta chain
	for _, chain := range zetaChains {
		require.NoError(t, address.ValidateAddressForChain("bcrt1qs758ursh4q9z627kt3pp5yysm78ddny6txaqgw", chain.ChainId, additionalChains))
		require.NoError(t, address.ValidateAddressForChain("327z4GyFM8Y8DiYfasGKQWhRK4MvyMSEgE", chain.ChainId, additionalChains))
		require.NoError(t, address.ValidateAddressForChain("bc1p4scddlkkuw9486579autxumxmkvuphm5pz4jvf7f6pdh50p2uzqstawjt9", chain.ChainId, additionalChains))
		require.NoError(t, address.ValidateAddressForChain(sample.EthAddress().Hex(), chain.ChainId, additionalChains))
		require.NoError(t, address.ValidateAddressForChain(sample.SolanaAddress(t), chain.ChainId, additionalChains))
	}
	// test for base chain
	for _, chain := range baseChains {
		require.NoError(t, address.ValidateAddressForChain("bcrt1qs758ursh4q9z627kt3pp5yysm78ddny6txaqgw", chain.ChainId, additionalChains))
		require.NoError(t, address.ValidateAddressForChain("327z4GyFM8Y8DiYfasGKQWhRK4MvyMSEgE", chain.ChainId, additionalChains))
		require.NoError(t, address.ValidateAddressForChain("bc1p4scddlkkuw9486579autxumxmkvuphm5pz4jvf7f6pdh50p2uzqstawjt9", chain.ChainId, additionalChains))
		require.NoError(t, address.ValidateAddressForChain(sample.EthAddress().Hex(), chain.ChainId, additionalChains))
		require.NoError(t, address.ValidateAddressForChain(sample.SolanaAddress(t), chain.ChainId, additionalChains))
	}

	// test for solana chain
	for _, chain := range solanaChains {
		require.NoError(t, address.ValidateAddressForChain("bcrt1qs758ursh4q9z627kt3pp5yysm78ddny6txaqgw", chain.ChainId, additionalChains))
		require.NoError(t, address.ValidateAddressForChain("327z4GyFM8Y8DiYfasGKQWhRK4MvyMSEgE", chain.ChainId, additionalChains))
		require.NoError(t, address.ValidateAddressForChain("bc1p4scddlkkuw9486579autxumxmkvuphm5pz4jvf7f6pdh50p2uzqstawjt9", chain.ChainId, additionalChains))
		require.NoError(t, address.ValidateAddressForChain(sample.EthAddress().Hex(), chain.ChainId, additionalChains))
		require.NoError(t, address.ValidateAddressForChain(sample.SolanaAddress(t), chain.ChainId, additionalChains))
	}

	// test for optimism chain
	for _, chain := range opChains {
		require.NoError(t, address.ValidateAddressForChain("bcrt1qs758ursh4q9z627kt3pp5yysm78ddny6txaqgw", chain.ChainId, additionalChains))
		require.NoError(t, address.ValidateAddressForChain("327z4GyFM8Y8DiYfasGKQWhRK4MvyMSEgE", chain.ChainId, additionalChains))
		require.NoError(t, address.ValidateAddressForChain("bc1p4scddlkkuw9486579autxumxmkvuphm5pz4jvf7f6pdh50p2uzqstawjt9", chain.ChainId, additionalChains))
		require.NoError(t, address.ValidateAddressForChain(sample.EthAddress().Hex(), chain.ChainId, additionalChains))
		require.NoError(t, address.ValidateAddressForChain(sample.SolanaAddress(t), chain.ChainId, additionalChains))
	}

	// test for btc chain

	require.NoError(
		t,
		address.ValidateAddressForChain(
			"bc1p4scddlkkuw9486579autxumxmkvuphm5pz4jvf7f6pdh50p2uzqstawjt9",
			chains.BitcoinMainnet.ChainId,
			additionalChains,
		),
	)
	require.NoError(
		t,
		address.ValidateAddressForChain("327z4GyFM8Y8DiYfasGKQWhRK4MvyMSEgE", chains.BitcoinMainnet.ChainId, additionalChains),
	)
	require.NoError(
		t,
		address.ValidateAddressForChain("1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3", chains.BitcoinMainnet.ChainId, additionalChains),
	)
	require.NoError(
		t,
		address.ValidateAddressForChain("bc1qysd4sp9q8my59ul9wsf5rvs9p387hf8vfwatzu", chains.BitcoinMainnet.ChainId, additionalChains),
	)
	require.NoError(
		t,
		address.ValidateAddressForChain("bcrt1qs758ursh4q9z627kt3pp5yysm78ddny6txaqgw", chains.BitcoinRegtest.ChainId, additionalChains),
	)

	// Invalid if the chain params are incorrect
	require.Error(
		t,
		address.ValidateAddressForChain("bcrt1qs758ursh4q9z627kt3pp5yysm78ddny6txaqgw", chains.BitcoinMainnet.ChainId, additionalChains),
	)
	// Invalid if address string is invalid
	require.Error(t, address.ValidateAddressForChain("", chains.BitcoinRegtest.ChainId, additionalChains))
}

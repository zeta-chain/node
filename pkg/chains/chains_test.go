package chains_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
)

func TestChain_Name(t *testing.T) {
	t.Run("new Name field is compatible with ChainName enum", func(t *testing.T) {
		for _, chain := range chains.DefaultChainsList() {
			if chain.ChainName != chains.ChainName_empty {
				require.EqualValues(t, chain.Name, chain.ChainName.String())
			}
		}
	})
}

func TestChainListByNetworkType(t *testing.T) {
	listTests := []struct {
		name        string
		networkType chains.NetworkType
		expected    []chains.Chain
	}{
		{
			"mainnet chains",
			chains.NetworkType_mainnet,
			[]chains.Chain{
				chains.ZetaChainMainnet,
				chains.BitcoinMainnet,
				chains.BscMainnet,
				chains.Ethereum,
				chains.Polygon,
				chains.OptimismMainnet,
				chains.BaseMainnet,
				chains.SolanaMainnet,
				chains.TONMainnet,
				chains.AvalancheMainnet,
				chains.ArbitrumMainnet,
				chains.WorldMainnet,
				chains.SuiMainnet,
			},
		},
		{
			"testnet chains",
			chains.NetworkType_testnet,
			[]chains.Chain{
				chains.ZetaChainTestnet,
				chains.BitcoinTestnet,
				chains.BitcoinSignetTestnet,
				chains.BitcoinTestnet4,
				chains.Mumbai,
				chains.Amoy,
				chains.BscTestnet,
				chains.Goerli,
				chains.Sepolia,
				chains.OptimismSepolia,
				chains.BaseSepolia,
				chains.SolanaDevnet,
				chains.TONTestnet,
				chains.AvalancheTestnet,
				chains.ArbitrumSepolia,
				chains.WorldTestnet,
				chains.SuiTestnet,
			},
		},
		{
			"privnet chains",
			chains.NetworkType_privnet,
			[]chains.Chain{
				chains.ZetaChainPrivnet,
				chains.BitcoinRegtest,
				chains.GoerliLocalnet,
				chains.SolanaLocalnet,
				chains.TONLocalnet,
				chains.SuiLocalnet,
			},
		},
	}

	for _, lt := range listTests {
		t.Run(lt.name, func(t *testing.T) {
			require.ElementsMatch(t, lt.expected, chains.ChainListByNetworkType(lt.networkType, []chains.Chain{}))
		})
	}
}

func TestChainListByNetwork(t *testing.T) {
	listTests := []struct {
		name     string
		network  chains.Network
		expected []chains.Chain
	}{
		{
			"Zeta",
			chains.Network_zeta,
			[]chains.Chain{
				chains.ZetaChainMainnet,
				chains.ZetaChainDevnet,
				chains.ZetaChainPrivnet,
				chains.ZetaChainTestnet,
			},
		},
		{
			"Btc",
			chains.Network_btc,
			[]chains.Chain{
				chains.BitcoinMainnet,
				chains.BitcoinTestnet,
				chains.BitcoinSignetTestnet,
				chains.BitcoinTestnet4,
				chains.BitcoinRegtest,
			},
		},
		{
			"Eth",
			chains.Network_eth,
			[]chains.Chain{chains.Ethereum, chains.Goerli, chains.Sepolia, chains.GoerliLocalnet},
		},
		{
			"Bsc",
			chains.Network_bsc,
			[]chains.Chain{chains.BscMainnet, chains.BscTestnet},
		},
		{
			"Polygon",
			chains.Network_polygon,
			[]chains.Chain{chains.Polygon, chains.Mumbai, chains.Amoy},
		},
		{
			"Optimism",
			chains.Network_optimism,
			[]chains.Chain{chains.OptimismMainnet, chains.OptimismSepolia},
		},
		{
			"Base",
			chains.Network_base,
			[]chains.Chain{chains.BaseMainnet, chains.BaseSepolia},
		},
		{
			"Solana",
			chains.Network_solana,
			[]chains.Chain{chains.SolanaMainnet, chains.SolanaDevnet, chains.SolanaLocalnet},
		},
		{
			"TON",
			chains.Network_ton,
			[]chains.Chain{chains.TONMainnet, chains.TONTestnet, chains.TONLocalnet},
		},
		{
			"Sui",
			chains.Network_sui,
			[]chains.Chain{chains.SuiMainnet, chains.SuiTestnet, chains.SuiLocalnet},
		},
	}

	for _, lt := range listTests {
		t.Run(lt.name, func(t *testing.T) {
			require.ElementsMatch(t, lt.expected, chains.ChainListByNetwork(lt.network, []chains.Chain{}))
		})
	}
}

func TestDefaultChainList(t *testing.T) {
	require.ElementsMatch(t, []chains.Chain{
		chains.BitcoinMainnet,
		chains.BscMainnet,
		chains.Ethereum,
		chains.BitcoinTestnet,
		chains.BitcoinSignetTestnet,
		chains.BitcoinTestnet4,
		chains.Mumbai,
		chains.Amoy,
		chains.BscTestnet,
		chains.Goerli,
		chains.Sepolia,
		chains.BitcoinRegtest,
		chains.GoerliLocalnet,
		chains.ZetaChainMainnet,
		chains.ZetaChainTestnet,
		chains.ZetaChainDevnet,
		chains.ZetaChainPrivnet,
		chains.Polygon,
		chains.OptimismMainnet,
		chains.OptimismSepolia,
		chains.BaseMainnet,
		chains.BaseSepolia,
		chains.SolanaMainnet,
		chains.SolanaDevnet,
		chains.SolanaLocalnet,
		chains.TONMainnet,
		chains.TONTestnet,
		chains.TONLocalnet,
		chains.AvalancheMainnet,
		chains.AvalancheTestnet,
		chains.ArbitrumSepolia,
		chains.ArbitrumMainnet,
		chains.WorldMainnet,
		chains.WorldTestnet,
		chains.SuiMainnet,
		chains.SuiTestnet,
		chains.SuiLocalnet,
	}, chains.DefaultChainsList())
}

func TestChainListByGateway(t *testing.T) {
	listTests := []struct {
		name     string
		gateway  chains.CCTXGateway
		expected []chains.Chain
	}{
		{
			"observers",
			chains.CCTXGateway_observers,
			[]chains.Chain{
				chains.BitcoinMainnet,
				chains.BscMainnet,
				chains.Ethereum,
				chains.BitcoinTestnet,
				chains.BitcoinSignetTestnet,
				chains.BitcoinTestnet4,
				chains.Mumbai,
				chains.Amoy,
				chains.BscTestnet,
				chains.Goerli,
				chains.Sepolia,
				chains.BitcoinRegtest,
				chains.GoerliLocalnet,
				chains.Polygon,
				chains.OptimismMainnet,
				chains.OptimismSepolia,
				chains.BaseMainnet,
				chains.BaseSepolia,
				chains.SolanaMainnet,
				chains.SolanaDevnet,
				chains.SolanaLocalnet,
				chains.TONMainnet,
				chains.TONTestnet,
				chains.TONLocalnet,
				chains.AvalancheMainnet,
				chains.AvalancheTestnet,
				chains.ArbitrumSepolia,
				chains.ArbitrumMainnet,
				chains.WorldMainnet,
				chains.WorldTestnet,
				chains.SuiMainnet,
				chains.SuiTestnet,
				chains.SuiLocalnet,
			},
		},
		{
			"zevm",
			chains.CCTXGateway_zevm,
			[]chains.Chain{
				chains.ZetaChainMainnet,
				chains.ZetaChainTestnet,
				chains.ZetaChainDevnet,
				chains.ZetaChainPrivnet,
			},
		},
	}

	for _, lt := range listTests {
		t.Run(lt.name, func(t *testing.T) {
			require.ElementsMatch(t, lt.expected, chains.ChainListByGateway(lt.gateway, []chains.Chain{}))
		})
	}
}

func TestExternalChainList(t *testing.T) {
	require.ElementsMatch(t, []chains.Chain{
		chains.BitcoinMainnet,
		chains.BscMainnet,
		chains.Ethereum,
		chains.BitcoinTestnet,
		chains.BitcoinSignetTestnet,
		chains.BitcoinTestnet4,
		chains.Mumbai,
		chains.Amoy,
		chains.BscTestnet,
		chains.Goerli,
		chains.Sepolia,
		chains.BitcoinRegtest,
		chains.GoerliLocalnet,
		chains.Polygon,
		chains.OptimismMainnet,
		chains.OptimismSepolia,
		chains.BaseMainnet,
		chains.BaseSepolia,
		chains.SolanaMainnet,
		chains.SolanaDevnet,
		chains.SolanaLocalnet,
		chains.TONMainnet,
		chains.TONTestnet,
		chains.TONLocalnet,
		chains.AvalancheMainnet,
		chains.AvalancheTestnet,
		chains.ArbitrumSepolia,
		chains.ArbitrumMainnet,
		chains.WorldMainnet,
		chains.WorldTestnet,
		chains.SuiMainnet,
		chains.SuiTestnet,
		chains.SuiLocalnet,
	}, chains.ExternalChainList([]chains.Chain{}))
}

func TestZetaChainFromCosmosChainID(t *testing.T) {
	tests := []struct {
		name     string
		chainID  string
		expected chains.Chain
		wantErr  bool
	}{
		{
			name:     "ZetaChainMainnet",
			chainID:  "cosmoshub_7000-1",
			expected: chains.ZetaChainMainnet,
			wantErr:  false,
		},
		{
			name:     "ZetaChainTestnet",
			chainID:  "cosmoshub_7001-1",
			expected: chains.ZetaChainTestnet,
			wantErr:  false,
		},
		{
			name:     "ZetaChainDevnet",
			chainID:  "cosmoshub_70000-1",
			expected: chains.ZetaChainDevnet,
			wantErr:  false,
		},
		{
			name:     "ZetaChainPrivnet",
			chainID:  "cosmoshub_101-1",
			expected: chains.ZetaChainPrivnet,
			wantErr:  false,
		},
		{
			name:     "unknown chain",
			chainID:  "cosmoshub_1234-1",
			expected: chains.Chain{},
			wantErr:  true,
		},
		{
			name:     "invalid chain id",
			chainID:  "cosmoshub_abc-1",
			expected: chains.Chain{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := chains.ZetaChainFromCosmosChainID(tt.chainID)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, chains.Chain{}, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestZetaChainFromChainID(t *testing.T) {
	tests := []struct {
		name     string
		chainID  int64
		expected chains.Chain
		wantErr  bool
	}{
		{
			name:     "ZetaChainMainnet",
			chainID:  7000,
			expected: chains.ZetaChainMainnet,
			wantErr:  false,
		},
		{
			name:     "ZetaChainTestnet",
			chainID:  7001,
			expected: chains.ZetaChainTestnet,
			wantErr:  false,
		},
		{
			name:     "ZetaChainDevnet",
			chainID:  70000,
			expected: chains.ZetaChainDevnet,
			wantErr:  false,
		},
		{
			name:     "ZetaChainPrivnet",
			chainID:  101,
			expected: chains.ZetaChainPrivnet,
			wantErr:  false,
		},
		{
			name:     "unknown chain",
			chainID:  1234,
			expected: chains.Chain{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := chains.ZetaChainFromChainID(tt.chainID)
			if tt.wantErr {
				require.Error(t, err)
				require.ErrorIs(t, err, chains.ErrNotZetaChain)
				require.Equal(t, chains.Chain{}, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestCombineDefaultChainsList(t *testing.T) {
	// prepare array containing pre-defined chains
	// chain IDs are 11000 - 11009 to not conflict with the default chains
	var chainList = make([]chains.Chain, 0, 10)
	for i := int64(11000); i < 10; i++ {
		chainList = append(chainList, sample.Chain(i))
	}

	bitcoinMainnetChainID := chains.BitcoinMainnet.ChainId
	require.Equal(
		t,
		bitcoinMainnetChainID,
		chains.DefaultChainsList()[0].ChainId,
		"Bitcoin mainnet be the first in the default chain list for TestCombineDefaultChainsList tests",
	)
	alternativeBitcoinMainnet := sample.Chain(bitcoinMainnetChainID)

	tests := []struct {
		name     string
		list     []chains.Chain
		expected []chains.Chain
	}{
		{
			name:     "empty list",
			list:     []chains.Chain{},
			expected: chains.DefaultChainsList(),
		},
		{
			name:     "no duplicates",
			list:     chainList,
			expected: append(chains.DefaultChainsList(), chainList...),
		},
		{
			name:     "duplicates",
			list:     []chains.Chain{alternativeBitcoinMainnet},
			expected: append(chains.DefaultChainsList()[1:], alternativeBitcoinMainnet),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.ElementsMatch(t, tt.expected, chains.CombineDefaultChainsList(tt.list))
		})
	}
}

func TestCombineChainList(t *testing.T) {
	// prepare array containing pre-defined chains
	var chainList = make([]chains.Chain, 0, 10)
	for i := int64(0); i < 10; i++ {
		chainList = append(chainList, sample.Chain(i))
	}

	// prepare second array for duplicated chain IDs
	var duplicatedChainList = make([]chains.Chain, 0, 10)
	for i := int64(0); i < 10; i++ {
		duplicatedChainList = append(duplicatedChainList, sample.Chain(i))
	}

	tests := []struct {
		name     string
		list1    []chains.Chain
		list2    []chains.Chain
		expected []chains.Chain
	}{
		{
			name:     "empty lists",
			list1:    []chains.Chain{},
			list2:    []chains.Chain{},
			expected: []chains.Chain{},
		},
		{
			name:     "empty list 1",
			list1:    []chains.Chain{},
			list2:    chainList,
			expected: chainList,
		},
		{
			name:     "empty list 2",
			list1:    chainList,
			list2:    []chains.Chain{},
			expected: chainList,
		},
		{
			name:     "no duplicates",
			list1:    chainList[:5],
			list2:    chainList[5:],
			expected: chainList,
		},
		{
			name:     "all duplicates",
			list1:    chainList,
			list2:    duplicatedChainList,
			expected: duplicatedChainList,
		},
		{
			name:     "some duplicates",
			list1:    chainList[:5],
			list2:    duplicatedChainList[3:],
			expected: append(chainList[:3], duplicatedChainList[3:]...),
		},
		{
			name:     "one duplicate",
			list1:    chainList[:5],
			list2:    append(chainList[5:], duplicatedChainList[0]),
			expected: append(chainList[1:], duplicatedChainList[0]),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.ElementsMatch(t, tt.expected, chains.CombineChainList(tt.list1, tt.list2))
		})
	}
}

package chains_test

import (
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"testing"
)

func TestChainListByNetworkType(t *testing.T) {
	listTests := []struct {
		name        string
		networkType chains.NetworkType
		expected    []*chains.Chain
	}{
		{
			"mainnet chains",
			chains.NetworkType_mainnet,
			[]*chains.Chain{
				&chains.ZetaChainMainnet,
				&chains.BitcoinMainnet,
				&chains.BscMainnet,
				&chains.Ethereum,
				&chains.Polygon,
				&chains.OptimismMainnet,
				&chains.BaseMainnet,
			},
		},
		{
			"testnet chains",
			chains.NetworkType_testnet,
			[]*chains.Chain{
				&chains.ZetaChainTestnet,
				&chains.BitcoinTestnet,
				&chains.Mumbai,
				&chains.Amoy,
				&chains.BscTestnet,
				&chains.Goerli,
				&chains.Sepolia,
				&chains.OptimismSepolia,
				&chains.BaseSepolia,
			},
		},
		{
			"privnet chains",
			chains.NetworkType_privnet,
			[]*chains.Chain{
				&chains.ZetaChainPrivnet,
				&chains.BitcoinRegtest,
				&chains.GoerliLocalnet,
			},
		},
	}

	for _, lt := range listTests {
		t.Run(lt.name, func(t *testing.T) {
			require.ElementsMatch(t, lt.expected, chains.ChainListByNetworkType(lt.networkType))
		})
	}
}

func TestChainListByNetwork(t *testing.T) {
	listTests := []struct {
		name     string
		network  chains.Network
		expected []*chains.Chain
	}{
		{
			"Zeta",
			chains.Network_zeta,
			[]*chains.Chain{&chains.ZetaChainMainnet, &chains.ZetaChainDevnet, &chains.ZetaChainPrivnet, &chains.ZetaChainTestnet},
		},
		{
			"Btc",
			chains.Network_btc,
			[]*chains.Chain{&chains.BitcoinMainnet, &chains.BitcoinTestnet, &chains.BitcoinRegtest},
		},
		{
			"Eth",
			chains.Network_eth,
			[]*chains.Chain{&chains.Ethereum, &chains.Goerli, &chains.Sepolia, &chains.GoerliLocalnet},
		},
		{
			"Bsc",
			chains.Network_bsc,
			[]*chains.Chain{&chains.BscMainnet, &chains.BscTestnet},
		},
		{
			"Polygon",
			chains.Network_polygon,
			[]*chains.Chain{&chains.Polygon, &chains.Mumbai, &chains.Amoy},
		},
		{
			"Optimism",
			chains.Network_optimism,
			[]*chains.Chain{&chains.OptimismMainnet, &chains.OptimismSepolia},
		},
		{
			"Base",
			chains.Network_base,
			[]*chains.Chain{&chains.BaseMainnet, &chains.BaseSepolia},
		},
	}

	for _, lt := range listTests {
		t.Run(lt.name, func(t *testing.T) {
			require.ElementsMatch(t, lt.expected, chains.ChainListByNetwork(lt.network))
		})
	}
}
func TestChainListFunctions(t *testing.T) {
	listTests := []struct {
		name     string
		function func() []*chains.Chain
		expected []*chains.Chain
	}{
		{
			"DefaultChainsList",
			chains.DefaultChainsList,
			[]*chains.Chain{
				&chains.BitcoinMainnet,
				&chains.BscMainnet,
				&chains.Ethereum,
				&chains.BitcoinTestnet,
				&chains.Mumbai,
				&chains.Amoy,
				&chains.BscTestnet,
				&chains.Goerli,
				&chains.Sepolia,
				&chains.BitcoinRegtest,
				&chains.GoerliLocalnet,
				&chains.ZetaChainMainnet,
				&chains.ZetaChainTestnet,
				&chains.ZetaChainDevnet,
				&chains.ZetaChainPrivnet,
				&chains.Polygon,
				&chains.OptimismMainnet,
				&chains.OptimismSepolia,
				&chains.BaseMainnet,
				&chains.BaseSepolia,
			},
		},
		{
			"ExternalChainList",
			chains.ExternalChainList,
			[]*chains.Chain{
				&chains.BitcoinMainnet,
				&chains.BscMainnet,
				&chains.Ethereum,
				&chains.BitcoinTestnet,
				&chains.Mumbai,
				&chains.Amoy,
				&chains.BscTestnet,
				&chains.Goerli,
				&chains.Sepolia,
				&chains.BitcoinRegtest,
				&chains.GoerliLocalnet,
				&chains.Polygon,
				&chains.OptimismMainnet,
				&chains.OptimismSepolia,
				&chains.BaseMainnet,
				&chains.BaseSepolia,
			},
		},
	}

	for _, lt := range listTests {
		t.Run(lt.name, func(t *testing.T) {
			chains := lt.function()
			require.ElementsMatch(t, lt.expected, chains)
		})
	}
}

func TestZetaChainFromChainID(t *testing.T) {
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
			result, err := chains.ZetaChainFromChainID(tt.chainID)
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

func TestCombineChainList(t *testing.T) {
	// prepare array containing pre-defined chains
	var chainList = make([]*chains.Chain, 0, 10)
	for i := int64(0); i < 10; i++ {
		chainList = append(chainList, sample.Chain(i))
	}

	// prepare second array for duplicated chain IDs
	var duplicatedChainList = make([]*chains.Chain, 0, 10)
	for i := int64(0); i < 10; i++ {
		duplicatedChainList = append(duplicatedChainList, sample.Chain(i))
	}

	tests := []struct {
		name     string
		list1    []*chains.Chain
		list2    []*chains.Chain
		expected []*chains.Chain
	}{
		{
			name:     "empty lists",
			list1:    []*chains.Chain{},
			list2:    []*chains.Chain{},
			expected: []*chains.Chain{},
		},
		{
			name:     "empty list 1",
			list1:    []*chains.Chain{},
			list2:    chainList,
			expected: chainList,
		},
		{
			name:     "empty list 2",
			list1:    chainList,
			list2:    []*chains.Chain{},
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

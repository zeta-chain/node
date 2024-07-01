package chains

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChainListByNetworkType(t *testing.T) {
	listTests := []struct {
		name        string
		networkType NetworkType
		expected    []*Chain
	}{
		{
			"mainnet chains",
			NetworkType_mainnet,
			[]*Chain{
				&ZetaChainMainnet,
				&BitcoinMainnet,
				&BscMainnet,
				&Ethereum,
				&Polygon,
				&OptimismMainnet,
				&BaseMainnet,
				&SolanaMainnet,
			},
		},
		{
			"testnet chains",
			NetworkType_testnet,
			[]*Chain{
				&ZetaChainTestnet,
				&BitcoinTestnet,
				&Mumbai,
				&Amoy,
				&BscTestnet,
				&Goerli,
				&Sepolia,
				&OptimismSepolia,
				&BaseSepolia,
				&SolanaTestnet,
			},
		},
		{
			"privnet chains",
			NetworkType_privnet,
			[]*Chain{
				&ZetaChainPrivnet,
				&BitcoinRegtest,
				&GoerliLocalnet,
				&SolanaLocalnet,
			},
		},
	}

	for _, lt := range listTests {
		t.Run(lt.name, func(t *testing.T) {
			require.ElementsMatch(t, lt.expected, ChainListByNetworkType(lt.networkType))
		})
	}
}

func TestChainListByNetwork(t *testing.T) {
	listTests := []struct {
		name     string
		network  Network
		expected []*Chain
	}{
		{
			"Zeta",
			Network_zeta,
			[]*Chain{&ZetaChainMainnet, &ZetaChainDevnet, &ZetaChainPrivnet, &ZetaChainTestnet},
		},
		{
			"Btc",
			Network_btc,
			[]*Chain{&BitcoinMainnet, &BitcoinTestnet, &BitcoinRegtest},
		},
		{
			"Eth",
			Network_eth,
			[]*Chain{&Ethereum, &Goerli, &Sepolia, &GoerliLocalnet},
		},
		{
			"Bsc",
			Network_bsc,
			[]*Chain{&BscMainnet, &BscTestnet},
		},
		{
			"Polygon",
			Network_polygon,
			[]*Chain{&Polygon, &Mumbai, &Amoy},
		},
		{
			"Optimism",
			Network_optimism,
			[]*Chain{&OptimismMainnet, &OptimismSepolia},
		},
		{
			"Base",
			Network_base,
			[]*Chain{&BaseMainnet, &BaseSepolia},
		},
		{
			"Solana",
			Network_solana,
			[]*Chain{&SolanaMainnet, &SolanaTestnet, &SolanaLocalnet},
		},
	}

	for _, lt := range listTests {
		t.Run(lt.name, func(t *testing.T) {
			require.ElementsMatch(t, lt.expected, ChainListByNetwork(lt.network))
		})
	}
}
func TestChainListFunctions(t *testing.T) {
	listTests := []struct {
		name     string
		function func() []*Chain
		expected []*Chain
	}{
		{
			"DefaultChainsList",
			DefaultChainsList,
			[]*Chain{
				&BitcoinMainnet,
				&BscMainnet,
				&Ethereum,
				&BitcoinTestnet,
				&Mumbai,
				&Amoy,
				&BscTestnet,
				&Goerli,
				&Sepolia,
				&BitcoinRegtest,
				&GoerliLocalnet,
				&ZetaChainMainnet,
				&ZetaChainTestnet,
				&ZetaChainDevnet,
				&ZetaChainPrivnet,
				&Polygon,
				&OptimismMainnet,
				&OptimismSepolia,
				&BaseMainnet,
				&BaseSepolia,
				&SolanaMainnet,
				&SolanaTestnet,
				&SolanaLocalnet,
			},
		},
		{
			"ExternalChainList",
			ExternalChainList,
			[]*Chain{
				&BitcoinMainnet,
				&BscMainnet,
				&Ethereum,
				&BitcoinTestnet,
				&Mumbai,
				&Amoy,
				&BscTestnet,
				&Goerli,
				&Sepolia,
				&BitcoinRegtest,
				&GoerliLocalnet,
				&Polygon,
				&OptimismMainnet,
				&OptimismSepolia,
				&BaseMainnet,
				&BaseSepolia,
				&SolanaMainnet,
				&SolanaTestnet,
				&SolanaLocalnet,
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
		expected Chain
		wantErr  bool
	}{
		{
			name:     "ZetaChainMainnet",
			chainID:  "cosmoshub_7000-1",
			expected: ZetaChainMainnet,
			wantErr:  false,
		},
		{
			name:     "ZetaChainTestnet",
			chainID:  "cosmoshub_7001-1",
			expected: ZetaChainTestnet,
			wantErr:  false,
		},
		{
			name:     "ZetaChainDevnet",
			chainID:  "cosmoshub_70000-1",
			expected: ZetaChainDevnet,
			wantErr:  false,
		},
		{
			name:     "ZetaChainPrivnet",
			chainID:  "cosmoshub_101-1",
			expected: ZetaChainPrivnet,
			wantErr:  false,
		},
		{
			name:     "unknown chain",
			chainID:  "cosmoshub_1234-1",
			expected: Chain{},
			wantErr:  true,
		},
		{
			name:     "invalid chain id",
			chainID:  "cosmoshub_abc-1",
			expected: Chain{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ZetaChainFromChainID(tt.chainID)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, Chain{}, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
			}
		})
	}
}

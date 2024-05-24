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
			},
		},
		{
			"testnet chains",
			NetworkType_testnet,
			[]*Chain{
				&ZetaTestnetChain,
				&BtcTestNetChain,
				&MumbaiChain,
				&AmoyChain,
				&BscTestnetChain,
				&GoerliChain,
				&SepoliaChain,
				&OptimismSepolia,
				&BaseSepolia,
			},
		},
		{
			"privnet chains",
			NetworkType_privnet,
			[]*Chain{
				&ZetaPrivnetChain,
				&BtcRegtestChain,
				&GoerliLocalnetChain,
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
			[]*Chain{&ZetaChainMainnet, &ZetaMocknetChain, &ZetaPrivnetChain, &ZetaTestnetChain},
		},
		{
			"Btc",
			Network_btc,
			[]*Chain{&BitcoinMainnet, &BtcTestNetChain, &BtcRegtestChain},
		},
		{
			"Eth",
			Network_eth,
			[]*Chain{&Ethereum, &GoerliChain, &SepoliaChain, &GoerliLocalnetChain},
		},
		{
			"Bsc",
			Network_bsc,
			[]*Chain{&BscMainnet, &BscTestnetChain},
		},
		{
			"Polygon",
			Network_polygon,
			[]*Chain{&Polygon, &MumbaiChain, &AmoyChain},
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
				&BtcTestNetChain,
				&MumbaiChain,
				&AmoyChain,
				&BscTestnetChain,
				&GoerliChain,
				&SepoliaChain,
				&BtcRegtestChain,
				&GoerliLocalnetChain,
				&ZetaChainMainnet,
				&ZetaTestnetChain,
				&ZetaMocknetChain,
				&ZetaPrivnetChain,
				&Polygon,
				&OptimismMainnet,
				&OptimismSepolia,
				&BaseMainnet,
				&BaseSepolia,
			},
		},
		{
			"ExternalChainList",
			ExternalChainList,
			[]*Chain{
				&BitcoinMainnet,
				&BscMainnet,
				&Ethereum,
				&BtcTestNetChain,
				&MumbaiChain,
				&AmoyChain,
				&BscTestnetChain,
				&GoerliChain,
				&SepoliaChain,
				&BtcRegtestChain,
				&GoerliLocalnetChain,
				&Polygon,
				&OptimismMainnet,
				&OptimismSepolia,
				&BaseMainnet,
				&BaseSepolia,
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
			name:     "ZetaTestnetChain",
			chainID:  "cosmoshub_7001-1",
			expected: ZetaTestnetChain,
			wantErr:  false,
		},
		{
			name:     "ZetaMocknetChain",
			chainID:  "cosmoshub_70000-1",
			expected: ZetaMocknetChain,
			wantErr:  false,
		},
		{
			name:     "ZetaPrivnetChain",
			chainID:  "cosmoshub_101-1",
			expected: ZetaPrivnetChain,
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

package chains

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChainRetrievalFunctions(t *testing.T) {
	tests := []struct {
		name     string
		function func() Chain
		expected Chain
	}{
		{"ZetaChainMainnet", ZetaChainMainnet, Chain{ChainName: ChainName_zeta_mainnet, ChainId: 7000}},
		{"ZetaTestnetChain", ZetaTestnetChain, Chain{ChainName: ChainName_zeta_testnet, ChainId: 7001}},
		{"ZetaMocknetChain", ZetaMocknetChain, Chain{ChainName: ChainName_zeta_mainnet, ChainId: 70000}},
		{"ZetaPrivnetChain", ZetaPrivnetChain, Chain{ChainName: ChainName_zeta_mainnet, ChainId: 101}},
		{"EthChain", EthChain, Chain{ChainName: ChainName_eth_mainnet, ChainId: 1}},
		{"BscMainnetChain", BscMainnetChain, Chain{ChainName: ChainName_bsc_mainnet, ChainId: 56}},
		{"BtcMainnetChain", BtcMainnetChain, Chain{ChainName: ChainName_btc_mainnet, ChainId: 8332}},
		{"PolygonChain", PolygonChain, Chain{ChainName: ChainName_polygon_mainnet, ChainId: 137}},
		{"SepoliaChain", SepoliaChain, Chain{ChainName: ChainName_sepolia_testnet, ChainId: 11155111}},
		{"GoerliChain", GoerliChain, Chain{ChainName: ChainName_goerli_testnet, ChainId: 5}},
		{"BscTestnetChain", BscTestnetChain, Chain{ChainName: ChainName_bsc_testnet, ChainId: 97}},
		{"BtcTestNetChain", BtcTestNetChain, Chain{ChainName: ChainName_btc_testnet, ChainId: 18332}},
		{"MumbaiChain", MumbaiChain, Chain{ChainName: ChainName_mumbai_testnet, ChainId: 80001}},
		{"AmoyChain", AmoyChain, Chain{ChainName: ChainName_amoy_testnet, ChainId: 80002}},
		{"BtcRegtestChain", BtcRegtestChain, Chain{ChainName: ChainName_btc_regtest, ChainId: 18444}},
		{"GoerliLocalnetChain", GoerliLocalnetChain, Chain{ChainName: ChainName_goerli_localnet, ChainId: 1337}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			chain := tc.function()
			require.Equal(t, tc.expected, chain)
		})
	}
}

func TestChainListFunctions(t *testing.T) {
	listTests := []struct {
		name     string
		function func() []*Chain
		expected []Chain
	}{
		{"DefaultChainsList", DefaultChainsList, []Chain{BtcMainnetChain(), BscMainnetChain(), EthChain(), BtcTestNetChain(), MumbaiChain(), AmoyChain(), BscTestnetChain(), GoerliChain(), SepoliaChain(), BtcRegtestChain(), GoerliLocalnetChain(), ZetaChainMainnet(), ZetaTestnetChain(), ZetaMocknetChain(), ZetaPrivnetChain()}},
		{"MainnetChainList", MainnetChainList, []Chain{ZetaChainMainnet(), BtcMainnetChain(), BscMainnetChain(), EthChain()}},
		{"TestnetChainList", TestnetChainList, []Chain{ZetaTestnetChain(), BtcTestNetChain(), MumbaiChain(), AmoyChain(), BscTestnetChain(), GoerliChain(), SepoliaChain()}},
		{"PrivnetChainList", PrivnetChainList, []Chain{ZetaPrivnetChain(), BtcRegtestChain(), GoerliLocalnetChain()}},
		{"ExternalChainList", ExternalChainList, []Chain{BtcMainnetChain(), BscMainnetChain(), EthChain(), BtcTestNetChain(), MumbaiChain(), AmoyChain(), BscTestnetChain(), GoerliChain(), SepoliaChain(), BtcRegtestChain(), GoerliLocalnetChain()}},
		{"ZetaChainList", ZetaChainList, []Chain{ZetaChainMainnet(), ZetaTestnetChain(), ZetaMocknetChain(), ZetaPrivnetChain()}},
	}

	for _, lt := range listTests {
		t.Run(lt.name, func(t *testing.T) {
			chains := lt.function()
			require.Equal(t, len(lt.expected), len(chains))
			for i, expectedChain := range lt.expected {
				require.Equal(t, &expectedChain, chains[i])
			}
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
			expected: ZetaChainMainnet(),
			wantErr:  false,
		},
		{
			name:     "ZetaTestnetChain",
			chainID:  "cosmoshub_7001-1",
			expected: ZetaTestnetChain(),
			wantErr:  false,
		},
		{
			name:     "ZetaMocknetChain",
			chainID:  "cosmoshub_70000-1",
			expected: ZetaMocknetChain(),
			wantErr:  false,
		},
		{
			name:     "ZetaPrivnetChain",
			chainID:  "cosmoshub_101-1",
			expected: ZetaPrivnetChain(),
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

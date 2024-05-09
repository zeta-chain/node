package chains

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChainRetrievalFunctions(t *testing.T) {
	tests := []struct {
		name     string
		chain    Chain
		expected Chain
	}{
		{"ZetaChainMainnet", ZetaChainMainnet, Chain{
			ChainName:   ChainName_zeta_mainnet,
			ChainId:     7000,
			Network:     Network_zeta,
			NetworkType: NetworkType_mainnet,
			IsExternal:  false,
			Vm:          Vm_evm,
			Consensus:   Consensus_tendermint,
		},
		},
		{"ZetaTestnetChain", ZetaTestnetChain, Chain{
			ChainName:   ChainName_zeta_testnet,
			ChainId:     7001,
			Network:     Network_zeta,
			NetworkType: NetworkType_testnet,
			IsExternal:  false,
			Vm:          Vm_evm,
			Consensus:   Consensus_tendermint,
		},
		},
		{"ZetaMocknetChain", ZetaMocknetChain, Chain{
			ChainName:   ChainName_zeta_mainnet,
			ChainId:     70000,
			Network:     Network_zeta,
			NetworkType: NetworkType_devnet,
			IsExternal:  false,
			Vm:          Vm_evm,
			Consensus:   Consensus_tendermint,
		}},
		{"ZetaPrivnetChain", ZetaPrivnetChain, Chain{
			ChainName:   ChainName_zeta_mainnet,
			ChainId:     101,
			Network:     Network_zeta,
			NetworkType: NetworkType_privnet,
			IsExternal:  false,
			Vm:          Vm_evm,
			Consensus:   Consensus_tendermint,
		}},
		{"EthChain", EthChain, Chain{
			ChainName:   ChainName_eth_mainnet,
			ChainId:     1,
			Network:     Network_eth,
			NetworkType: NetworkType_mainnet,
			IsExternal:  true,
			Vm:          Vm_evm,
			Consensus:   Consensus_ethereum,
		}},
		{"BscMainnetChain", BscMainnetChain, Chain{
			ChainName:   ChainName_bsc_mainnet,
			ChainId:     56,
			Network:     Network_bsc,
			NetworkType: NetworkType_mainnet,
			IsExternal:  true,
			Vm:          Vm_evm,
			Consensus:   Consensus_ethereum,
		}},
		{"BtcMainnetChain", BtcMainnetChain, Chain{
			ChainName:   ChainName_btc_mainnet,
			ChainId:     8332,
			Network:     Network_btc,
			NetworkType: NetworkType_mainnet,
			IsExternal:  true,
			Vm:          Vm_no_vm,
			Consensus:   Consensus_bitcoin,
		}},
		{"PolygonChain", PolygonChain, Chain{
			ChainName:   ChainName_polygon_mainnet,
			ChainId:     137,
			Network:     Network_polygon,
			NetworkType: NetworkType_mainnet,
			IsExternal:  true,
			Vm:          Vm_evm,
			Consensus:   Consensus_ethereum,
		}},
		{"SepoliaChain", SepoliaChain, Chain{
			ChainName:   ChainName_sepolia_testnet,
			ChainId:     11155111,
			Network:     Network_eth,
			NetworkType: NetworkType_testnet,
			IsExternal:  true,
			Vm:          Vm_evm,
			Consensus:   Consensus_ethereum,
		}},
		{"GoerliChain", GoerliChain, Chain{
			ChainName:   ChainName_goerli_testnet,
			ChainId:     5,
			Network:     Network_eth,
			NetworkType: NetworkType_testnet,
			IsExternal:  true,
			Vm:          Vm_evm,
			Consensus:   Consensus_ethereum,
		}},
		{"AmoyChain", AmoyChain, Chain{
			ChainName:   ChainName_amoy_testnet,
			ChainId:     80002,
			Network:     Network_polygon,
			NetworkType: NetworkType_testnet,
			IsExternal:  true,
			Vm:          Vm_evm,
			Consensus:   Consensus_ethereum,
		}},
		{"BscTestnetChain", BscTestnetChain, Chain{
			ChainName:   ChainName_bsc_testnet,
			ChainId:     97,
			Network:     Network_bsc,
			NetworkType: NetworkType_testnet,
			IsExternal:  true,
			Vm:          Vm_evm,
			Consensus:   Consensus_ethereum,
		}},
		{"MumbaiChain", MumbaiChain, Chain{
			ChainName:   ChainName_mumbai_testnet,
			ChainId:     80001,
			Network:     Network_polygon,
			NetworkType: NetworkType_testnet,
			IsExternal:  true,
			Vm:          Vm_evm,

			Consensus: Consensus_ethereum,
		}},
		{"BtcTestNetChain", BtcTestNetChain, Chain{
			ChainName:   ChainName_btc_testnet,
			ChainId:     18332,
			Network:     Network_btc,
			NetworkType: NetworkType_testnet,
			IsExternal:  true,
			Vm:          Vm_no_vm,
			Consensus:   Consensus_bitcoin,
		}},
		{"BtcRegtestChain", BtcRegtestChain, Chain{
			ChainName:   ChainName_btc_regtest,
			ChainId:     18444,
			Network:     Network_btc,
			NetworkType: NetworkType_privnet,
			IsExternal:  true,
			Vm:          Vm_no_vm,

			Consensus: Consensus_bitcoin,
		}},
		{"GoerliLocalnetChain", GoerliLocalnetChain, Chain{
			ChainName:   ChainName_goerli_localnet,
			ChainId:     1337,
			Network:     Network_eth,
			NetworkType: NetworkType_privnet,
			IsExternal:  true,
			Vm:          Vm_evm,
			Consensus:   Consensus_ethereum,
		}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			chain := tc.chain
			require.Equal(t, tc.expected, chain)
		})
	}
}
func TestChainListByNetworkType(t *testing.T) {
	listTests := []struct {
		name        string
		networkType NetworkType
		expected    []Chain
	}{
		{"MainnetChainList", NetworkType_mainnet, []Chain{ZetaChainMainnet, BtcMainnetChain, BscMainnetChain, EthChain, PolygonChain}},
		{"TestnetChainList", NetworkType_testnet, []Chain{ZetaTestnetChain, BtcTestNetChain, MumbaiChain, AmoyChain, BscTestnetChain, GoerliChain, SepoliaChain}},
		{"PrivnetChainList", NetworkType_privnet, []Chain{ZetaPrivnetChain, BtcRegtestChain, GoerliLocalnetChain}},
	}

	for _, lt := range listTests {
		t.Run(lt.name, func(t *testing.T) {
			chains := ChainListByNetworkType(lt.networkType)
			require.Equal(t, len(lt.expected), len(chains))
			sort.Slice(chains, func(i, j int) bool {
				return chains[i].ChainId < chains[j].ChainId
			})
			sort.Slice(lt.expected, func(i, j int) bool {
				return lt.expected[i].ChainId < lt.expected[j].ChainId
			})
			for i, expectedChain := range lt.expected {
				require.Equal(t, &expectedChain, chains[i])
			}
		})
	}
}

func TestChainListByNetwork(t *testing.T) {
	listTests := []struct {
		name     string
		network  Network
		expected []Chain
	}{
		{"Zeta", Network_zeta, []Chain{ZetaChainMainnet, ZetaMocknetChain, ZetaPrivnetChain, ZetaTestnetChain}},
		{"Btc", Network_btc, []Chain{BtcMainnetChain, BtcTestNetChain, BtcRegtestChain}},
		{"Eth", Network_eth, []Chain{EthChain, GoerliChain, SepoliaChain, GoerliLocalnetChain}},
		{"Bsc", Network_bsc, []Chain{BscMainnetChain, BscTestnetChain}},
		{"Polygon", Network_polygon, []Chain{PolygonChain, MumbaiChain, AmoyChain}},
	}

	for _, lt := range listTests {
		t.Run(lt.name, func(t *testing.T) {
			chains := ChainListByNetwork(lt.network)
			require.Equal(t, len(lt.expected), len(chains))
			sort.Slice(chains, func(i, j int) bool {
				return chains[i].ChainId < chains[j].ChainId
			})
			sort.Slice(lt.expected, func(i, j int) bool {
				return lt.expected[i].ChainId < lt.expected[j].ChainId
			})
			for i, expectedChain := range lt.expected {
				require.Equal(t, &expectedChain, chains[i])
			}
		})
	}
}
func TestChainListFunctions(t *testing.T) {
	listTests := []struct {
		name     string
		function func() []*Chain
		expected []Chain
	}{
		{"DefaultChainsList", DefaultChainsList, []Chain{BtcMainnetChain, BscMainnetChain, EthChain, BtcTestNetChain, MumbaiChain, AmoyChain, BscTestnetChain, GoerliChain, SepoliaChain, BtcRegtestChain, GoerliLocalnetChain, ZetaChainMainnet, ZetaTestnetChain, ZetaMocknetChain, ZetaPrivnetChain, PolygonChain}},
		{"ExternalChainList", ExternalChainList, []Chain{BtcMainnetChain, BscMainnetChain, EthChain, BtcTestNetChain, MumbaiChain, AmoyChain, BscTestnetChain, GoerliChain, SepoliaChain, BtcRegtestChain, GoerliLocalnetChain, PolygonChain}},
	}

	for _, lt := range listTests {
		t.Run(lt.name, func(t *testing.T) {
			chains := lt.function()
			require.Equal(t, len(lt.expected), len(chains))
			sort.Slice(chains, func(i, j int) bool {
				return chains[i].ChainId < chains[j].ChainId
			})
			sort.Slice(lt.expected, func(i, j int) bool {
				return lt.expected[i].ChainId < lt.expected[j].ChainId
			})
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

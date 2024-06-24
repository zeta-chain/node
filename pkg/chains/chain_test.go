package chains_test

import (
	"github.com/zeta-chain/zetacore/pkg/chains"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestChain_Validate(t *testing.T) {
	tests := []struct {
		name   string
		chain  chains.Chain
		errStr string
	}{
		{
			name: "should pass if chain is valid",
			chain: chains.Chain{
				ChainId:     42,
				ChainName:   chains.ChainName_empty,
				Network:     chains.Network_optimism,
				NetworkType: chains.NetworkType_testnet,
				Vm:          chains.Vm_evm,
				Consensus:   chains.Consensus_op_stack,
				IsExternal:  true,
			},
		},
		{
			name: "should error if chain ID is zero",
			chain: chains.Chain{
				ChainId:     0,
				ChainName:   chains.ChainName_empty,
				Network:     chains.Network_optimism,
				NetworkType: chains.NetworkType_testnet,
				Vm:          chains.Vm_evm,
				Consensus:   chains.Consensus_op_stack,
				IsExternal:  true,
			},
			errStr: "chain ID must be positive",
		},
		{
			name: "should error if chain ID is negative",
			chain: chains.Chain{
				ChainId:     0,
				ChainName:   chains.ChainName_empty,
				Network:     chains.Network_optimism,
				NetworkType: chains.NetworkType_testnet,
				Vm:          chains.Vm_evm,
				Consensus:   chains.Consensus_op_stack,
				IsExternal:  true,
			},
			errStr: "chain ID must be positive",
		},
		{
			name: "should error if chain name invalid",
			chain: chains.Chain{
				ChainId:     42,
				ChainName:   chains.ChainName_base_sepolia + 1,
				Network:     chains.Network_optimism,
				NetworkType: chains.NetworkType_testnet,
				Vm:          chains.Vm_evm,
				Consensus:   chains.Consensus_op_stack,
				IsExternal:  true,
			},
			errStr: "invalid chain name",
		},
		{
			name: "should error if network invalid",
			chain: chains.Chain{
				ChainId:     42,
				ChainName:   chains.ChainName_empty,
				Network:     chains.Network_base + 1,
				NetworkType: chains.NetworkType_testnet,
				Vm:          chains.Vm_evm,
				Consensus:   chains.Consensus_op_stack,
				IsExternal:  true,
			},
			errStr: "invalid network",
		},
		{
			name: "should error if network type invalid",
			chain: chains.Chain{
				ChainId:     42,
				ChainName:   chains.ChainName_empty,
				Network:     chains.Network_base,
				NetworkType: chains.NetworkType_devnet + 1,
				Vm:          chains.Vm_evm,
				Consensus:   chains.Consensus_op_stack,
				IsExternal:  true,
			},
			errStr: "invalid network type",
		},
		{
			name: "should error if vm invalid",
			chain: chains.Chain{
				ChainId:     42,
				ChainName:   chains.ChainName_empty,
				Network:     chains.Network_base,
				NetworkType: chains.NetworkType_devnet,
				Vm:          chains.Vm_evm + 1,
				Consensus:   chains.Consensus_op_stack,
				IsExternal:  true,
			},
			errStr: "invalid vm",
		},
		{
			name: "should error if consensus invalid",
			chain: chains.Chain{
				ChainId:     42,
				ChainName:   chains.ChainName_empty,
				Network:     chains.Network_base,
				NetworkType: chains.NetworkType_devnet,
				Vm:          chains.Vm_evm,
				Consensus:   chains.Consensus_op_stack + 1,
				IsExternal:  true,
			},
			errStr: "invalid consensus",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.errStr != "" {
				require.ErrorContains(t, tt.chain.Validate(), tt.errStr)
			} else {
				require.NoError(t, tt.chain.Validate())
			}
		})
	}

	t.Run("all default chains are valid", func(t *testing.T) {
		for _, chain := range chains.DefaultChainsList() {
			require.NoError(t, chain.Validate())
		}
	})
}

func TestChain_EncodeAddress(t *testing.T) {
	tests := []struct {
		name    string
		chain   chains.Chain
		b       []byte
		want    string
		wantErr bool
	}{
		{
			name: "should error if b is not a valid address on the bitcoin network",
			chain: chains.Chain{
				ChainName: chains.ChainName_btc_testnet,
				ChainId:   18332,
			},
			b:       []byte("bc1qk0cc73p8m7hswn8y2q080xa4e5pxapnqgp7h9c"),
			want:    "",
			wantErr: true,
		},
		{
			name: "should pass if b is a valid address on the network",
			chain: chains.Chain{
				ChainName: chains.ChainName_btc_mainnet,
				ChainId:   8332,
			},
			b:       []byte("bc1qk0cc73p8m7hswn8y2q080xa4e5pxapnqgp7h9c"),
			want:    "bc1qk0cc73p8m7hswn8y2q080xa4e5pxapnqgp7h9c",
			wantErr: false,
		},
		{
			name: "should error if b is not a valid address on the evm network",
			chain: chains.Chain{
				ChainName: chains.ChainName_goerli_testnet,
				ChainId:   5,
			},
			b:       ethcommon.Hex2Bytes("0x321"),
			want:    "",
			wantErr: true,
		},
		{
			name: "should pass if b is a valid address on the evm network",
			chain: chains.Chain{
				ChainName: chains.ChainName_goerli_testnet,
				ChainId:   5,
			},
			b:       []byte("0x321"),
			want:    "0x0000000000000000000000000000003078333231",
			wantErr: false,
		},
		{
			name: "should error if chain not supported",
			chain: chains.Chain{
				ChainName: 999,
				ChainId:   999,
			},
			b:       ethcommon.Hex2Bytes("0x321"),
			want:    "",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			s, err := tc.chain.EncodeAddress(tc.b)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.Equal(t, tc.want, s)
		})
	}
}

func TestIsZetaChain(t *testing.T) {
	tests := []struct {
		name    string
		chainID int64
		want    bool
	}{
		{"Zeta Mainnet", chains.ZetaChainMainnet.ChainId, true},
		{"Zeta Testnet", chains.ZetaChainTestnet.ChainId, true},
		{"Zeta Mocknet", chains.ZetaChainDevnet.ChainId, true},
		{"Zeta Privnet", chains.ZetaChainPrivnet.ChainId, true},
		{"Non-Zeta", chains.Ethereum.ChainId, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, chains.IsZetaChain(tt.chainID))
		})
	}
}

func TestIsEVMChain(t *testing.T) {
	tests := []struct {
		name    string
		chainID int64
		want    bool
	}{
		{"Ethereum Mainnet", chains.Ethereum.ChainId, true},
		{"Goerli Testnet", chains.Goerli.ChainId, true},
		{"Sepolia Testnet", chains.Sepolia.ChainId, true},
		{"Non-EVM", chains.BitcoinMainnet.ChainId, false},
		{"Zeta Mainnet", chains.ZetaChainMainnet.ChainId, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, IsEVMChain(tt.chainID))
		})
	}
}

func TestIsHeaderSupportedChain(t *testing.T) {
	tests := []struct {
		name    string
		chainID int64
		want    bool
	}{
		{"Ethereum Mainnet", chains.Ethereum.ChainId, true},
		{"Goerli Testnet", chains.Goerli.ChainId, true},
		{"Goerli Localnet", chains.GoerliLocalnet.ChainId, true},
		{"Sepolia Testnet", chains.Sepolia.ChainId, true},
		{"BSC Testnet", chains.BscTestnet.ChainId, true},
		{"BSC Mainnet", chains.BscMainnet.ChainId, true},
		{"BTC", chains.BitcoinMainnet.ChainId, true},
		{"Zeta Mainnet", chains.ZetaChainMainnet.ChainId, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, chains.IsHeaderSupportedChain(tt.chainID))
		})
	}
}

func TestSupportMerkleProof(t *testing.T) {
	tests := []struct {
		name  string
		chain chains.Chain
		want  bool
	}{
		{"Ethereum Mainnet", chains.Ethereum, true},
		{"BSC Testnet", chains.BscTestnet, true},
		{"BSC Mainnet", chains.BscMainnet, true},
		{"Non-EVM", chains.BitcoinMainnet, true},
		{"Zeta Mainnet", chains.ZetaChainMainnet, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.chain.SupportMerkleProof())
		})
	}
}

func TestIsBitcoinChain(t *testing.T) {
	tests := []struct {
		name    string
		chainID int64
		want    bool
	}{
		{"Bitcoin Mainnet", chains.BitcoinMainnet.ChainId, true},
		{"Bitcoin Testnet", chains.BitcoinTestnet.ChainId, true},
		{"Bitcoin Regtest", chains.BitcoinRegtest.ChainId, true},
		{"Non-Bitcoin", chains.Ethereum.ChainId, false},
		{"Zeta Mainnet", chains.ZetaChainMainnet.ChainId, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, chains.IsBitcoinChain(tt.chainID))
		})
	}
}

func TestIsEthereumChain(t *testing.T) {
	tests := []struct {
		name    string
		chainID int64
		want    bool
	}{
		{"Ethereum Mainnet", chains.Ethereum.ChainId, true},
		{"Goerli Testnet", chains.Goerli.ChainId, true},
		{"Sepolia Testnet", chains.Sepolia.ChainId, true},
		{"Non-Ethereum", chains.BitcoinMainnet.ChainId, false},
		{"Zeta Mainnet", chains.ZetaChainMainnet.ChainId, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, chains.IsEthereumChain(tt.chainID))
		})
	}
}

func TestChain_IsExternalChain(t *testing.T) {
	require.False(t, chains.ZetaChainMainnet.IsExternalChain())
	require.True(t, chains.Ethereum.IsExternalChain())
}

func TestChain_IsZetaChain(t *testing.T) {
	require.True(t, chains.ZetaChainMainnet.IsZetaChain())
	require.False(t, chains.Ethereum.IsZetaChain())
}

func TestChain_IsEmpty(t *testing.T) {
	require.True(t, chains.Chain{}.IsEmpty())
	require.False(t, chains.ZetaChainMainnet.IsEmpty())
}

func TestGetChainFromChainID(t *testing.T) {
	chain := chains.GetChainFromChainID(chains.ZetaChainMainnet.ChainId)
	require.Equal(t, chains.ZetaChainMainnet, *chain)
	require.Nil(t, chains.GetChainFromChainID(9999))
}

func TestGetBTCChainParams(t *testing.T) {
	params, err := chains.GetBTCChainParams(chains.BitcoinMainnet.ChainId)
	require.NoError(t, err)
	require.Equal(t, &chaincfg.MainNetParams, params)

	_, err = chains.GetBTCChainParams(9999)
	require.Error(t, err)
}

func TestGetBTCChainIDFromChainParams(t *testing.T) {
	chainID, err := chains.GetBTCChainIDFromChainParams(&chaincfg.MainNetParams)
	require.NoError(t, err)
	require.Equal(t, int64(8332), chainID)

	chainID, err = chains.GetBTCChainIDFromChainParams(&chaincfg.RegressionNetParams)
	require.NoError(t, err)
	require.Equal(t, int64(18444), chainID)

	chainID, err = chains.GetBTCChainIDFromChainParams(&chaincfg.TestNet3Params)
	require.NoError(t, err)
	require.Equal(t, int64(18332), chainID)

	_, err = chains.GetBTCChainIDFromChainParams(&chaincfg.Params{Name: "unknown"})
	require.Error(t, err)
}

func TestChainIDInChainList(t *testing.T) {
	require.True(t, chains.ChainIDInChainList(chains.ZetaChainMainnet.ChainId, chains.ChainListByNetwork(chains.Network_zeta)))
	require.False(t, chains.ChainIDInChainList(chains.Ethereum.ChainId, chains.ChainListByNetwork(chains.Network_zeta)))
}

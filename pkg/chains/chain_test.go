package chains_test

import (
	"testing"

	"github.com/zeta-chain/zetacore/testutil/sample"

	"github.com/btcsuite/btcd/chaincfg"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
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
				Name:        "foo",
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
				Name:        "foo",
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
				Name:        "foo",
				Network:     chains.Network_optimism,
				NetworkType: chains.NetworkType_testnet,
				Vm:          chains.Vm_evm,
				Consensus:   chains.Consensus_op_stack,
				IsExternal:  true,
			},
			errStr: "chain ID must be positive",
		},
		{
			name: "should error if chain name empty",
			chain: chains.Chain{
				ChainId:     42,
				Name:        "",
				Network:     chains.Network_optimism,
				NetworkType: chains.NetworkType_testnet,
				Vm:          chains.Vm_evm,
				Consensus:   chains.Consensus_op_stack,
				IsExternal:  true,
			},
			errStr: "chain name cannot be empty",
		},
		{
			name: "should error if network invalid",
			chain: chains.Chain{
				ChainId:     42,
				Name:        "foo",
				Network:     chains.Network_solana + 1,
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
				Name:        "foo",
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
				Name:        "foo",
				Network:     chains.Network_base,
				NetworkType: chains.NetworkType_devnet,
				Vm:          chains.Vm_svm + 1,
				Consensus:   chains.Consensus_op_stack,
				IsExternal:  true,
			},
			errStr: "invalid vm",
		},
		{
			name: "should error if consensus invalid",
			chain: chains.Chain{
				ChainId:     42,
				Name:        "foo",
				Network:     chains.Network_base,
				NetworkType: chains.NetworkType_devnet,
				Vm:          chains.Vm_evm,
				Consensus:   chains.Consensus_solana_consensus + 1,
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
			name:    "should error if b is not a valid address on the bitcoin network",
			chain:   chains.BitcoinTestnet,
			b:       []byte("bc1qk0cc73p8m7hswn8y2q080xa4e5pxapnqgp7h9c"),
			want:    "",
			wantErr: true,
		},
		{
			name:    "should pass if b is a valid address on the network",
			chain:   chains.BitcoinMainnet,
			b:       []byte("bc1qk0cc73p8m7hswn8y2q080xa4e5pxapnqgp7h9c"),
			want:    "bc1qk0cc73p8m7hswn8y2q080xa4e5pxapnqgp7h9c",
			wantErr: false,
		},
		{
			name:    "should pass if b is a valid wallet address on the solana network",
			chain:   chains.SolanaMainnet,
			b:       []byte("DCAK36VfExkPdAkYUQg6ewgxyinvcEyPLyHjRbmveKFw"),
			want:    "DCAK36VfExkPdAkYUQg6ewgxyinvcEyPLyHjRbmveKFw",
			wantErr: false,
		},
		{
			name:    "should error if b is not a valid Base58 address",
			chain:   chains.SolanaMainnet,
			b:       []byte("9G0P8HkKqegZ7B6cE2hGvkZjHjSH14WZXDNZQmwYLokAc"), // contains invalid digit '0'
			want:    "",
			wantErr: true,
		},
		{
			name:    "should error if b is not a valid address on the evm network",
			chain:   chains.Ethereum,
			b:       ethcommon.Hex2Bytes("0x321"),
			want:    "",
			wantErr: true,
		},
		{
			name:    "should pass if b is a valid address on the evm network",
			chain:   chains.Ethereum,
			b:       []byte("0x321"),
			want:    "0x0000000000000000000000000000003078333231",
			wantErr: false,
		},
		{
			name: "should error if chain not supported",
			chain: chains.Chain{
				ChainId: 999,
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

func TestChain_IsEVMChain(t *testing.T) {
	tests := []struct {
		name  string
		chain chains.Chain
		want  bool
	}{
		{"Ethereum Mainnet", chains.Ethereum, true},
		{"Goerli Testnet", chains.Goerli, true},
		{"Sepolia Testnet", chains.Sepolia, true},
		{"Non-EVM", chains.BitcoinMainnet, false},
		{"Zeta Mainnet", chains.ZetaChainMainnet, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.chain.IsEVMChain())
		})
	}
}

func TestChain_IsBitcoinChain(t *testing.T) {
	tests := []struct {
		name  string
		chain chains.Chain
		want  bool
	}{
		{"Bitcoin Mainnet", chains.BitcoinMainnet, true},
		{"Bitcoin Testnet", chains.BitcoinTestnet, true},
		{"Bitcoin Regtest", chains.BitcoinRegtest, true},
		{"Non-Bitcoin", chains.Ethereum, false},
		{"Zeta Mainnet", chains.ZetaChainMainnet, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.chain.IsBitcoinChain())
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
			require.Equal(t, tt.want, chains.IsZetaChain(tt.chainID, []chains.Chain{}))
		})
	}
}

func TestDecodeAddressFromChainID(t *testing.T) {
	ethAddr := sample.EthAddress()

	tests := []struct {
		name    string
		chainID int64
		addr    string
		want    []byte
		wantErr bool
	}{
		{
			name:    "Ethereum",
			chainID: chains.Ethereum.ChainId,
			addr:    ethAddr.Hex(),
			want:    ethAddr.Bytes(),
		},
		{
			name:    "Bitcoin",
			chainID: chains.BitcoinMainnet.ChainId,
			addr:    "bc1qk0cc73p8m7hswn8y2q080xa4e5pxapnqgp7h9c",
			want:    []byte("bc1qk0cc73p8m7hswn8y2q080xa4e5pxapnqgp7h9c"),
		},
		{
			name:    "Solana",
			chainID: chains.SolanaMainnet.ChainId,
			addr:    "DCAK36VfExkPdAkYUQg6ewgxyinvcEyPLyHjRbmveKFw",
			want:    []byte("DCAK36VfExkPdAkYUQg6ewgxyinvcEyPLyHjRbmveKFw"),
		},
		{
			name:    "Non-supported chain",
			chainID: 9999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := chains.DecodeAddressFromChainID(tt.chainID, tt.addr, []chains.Chain{})
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
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
			require.Equal(t, tt.want, chains.IsEVMChain(tt.chainID, []chains.Chain{}))
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
			require.Equal(t, tt.want, chains.IsBitcoinChain(tt.chainID, []chains.Chain{}))
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
			require.Equal(t, tt.want, chains.IsEthereumChain(tt.chainID, []chains.Chain{}))
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
	chain, found := chains.GetChainFromChainID(chains.ZetaChainMainnet.ChainId, []chains.Chain{})
	require.EqualValues(t, chains.ZetaChainMainnet, chain)
	require.True(t, found)
	_, found = chains.GetChainFromChainID(9999, []chains.Chain{})
	require.False(t, found)
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
	require.True(
		t,
		chains.ChainIDInChainList(
			chains.ZetaChainMainnet.ChainId,
			chains.ChainListByNetwork(chains.Network_zeta, []chains.Chain{}),
		),
	)
	require.False(
		t,
		chains.ChainIDInChainList(
			chains.Ethereum.ChainId,
			chains.ChainListByNetwork(chains.Network_zeta, []chains.Chain{}),
		),
	)
}

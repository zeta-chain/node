package chains_test

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/contracts/sui"
	"github.com/zeta-chain/node/testutil/sample"
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
				Network:     chains.Network_sui + 1,
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
				Vm:          chains.Vm_mvm_sui + 1,
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
				Consensus:   chains.Consensus_sui_consensus + 1,
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
			name:    "should error if b is not a valid address on the bitcoin testnet network",
			chain:   chains.BitcoinTestnet,
			b:       []byte("bc1qk0cc73p8m7hswn8y2q080xa4e5pxapnqgp7h9c"),
			want:    "",
			wantErr: true,
		},
		{
			name:    "should error if b is not a valid address on the bitcoin signet network",
			chain:   chains.BitcoinSignetTestnet,
			b:       []byte("bc1qk0cc73p8m7hswn8y2q080xa4e5pxapnqgp7h9c"),
			want:    "",
			wantErr: true,
		},
		{
			name:    "should error if b is not a valid address on the bitcoin testnet4 network",
			chain:   chains.BitcoinTestnet4,
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
			name:    "valid bitcoin testnet address",
			chain:   chains.BitcoinTestnet,
			b:       []byte("tb1qy9pqmk2pd9sv63g27jt8r657wy0d9ueeh0nqur"),
			want:    "tb1qy9pqmk2pd9sv63g27jt8r657wy0d9ueeh0nqur",
			wantErr: false,
		},
		{
			name:    "valid bitcoin signet address",
			chain:   chains.BitcoinSignetTestnet,
			b:       []byte("tb1qy9pqmk2pd9sv63g27jt8r657wy0d9ueeh0nqur"),
			want:    "tb1qy9pqmk2pd9sv63g27jt8r657wy0d9ueeh0nqur",
			wantErr: false,
		},
		{
			name:    "valid bitcoin testnet4 address",
			chain:   chains.BitcoinTestnet4,
			b:       []byte("tb1qy9pqmk2pd9sv63g27jt8r657wy0d9ueeh0nqur"),
			want:    "tb1qy9pqmk2pd9sv63g27jt8r657wy0d9ueeh0nqur",
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
		t.Run(tc.name, func(t *testing.T) {
			s, err := tc.chain.EncodeAddress(tc.b)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
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
		{"Zeta Mainnet", chains.ZetaChainMainnet, true},
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
		{"Bitcoin Signet Testnet", chains.BitcoinSignetTestnet, true},
		{"Bitcoin Testnet4", chains.BitcoinTestnet4, true},
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

	suiSample := "0x2a4c5a97b561ac5b38edc4b4e9b2c183c57b56df5b1ea2f1c6f2e4a44b92d59f"
	suiExpected, err := sui.EncodeAddress(suiSample)
	require.NoError(t, err)

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
			name:    "TON",
			chainID: chains.TONMainnet.ChainId,
			addr:    "0:55798cb7b87168251a7c39f6806b8c202f6caa0f617a76f4070b3fdacfd056a1",
			want:    []byte("0:55798cb7b87168251a7c39f6806b8c202f6caa0f617a76f4070b3fdacfd056a1"),
		},
		{
			name:    "TON",
			chainID: chains.TONMainnet.ChainId,
			// human friendly address should be always represented in raw format
			addr: "EQB3ncyBUTjZUA5EnFKR5_EnOMI9V1tTEAAPaiU71gc4TiUt",
			want: []byte("0:779dcc815138d9500e449c5291e7f12738c23d575b5310000f6a253bd607384e"),
		},
		{
			name:    "Sui",
			chainID: chains.SuiMainnet.ChainId,
			addr:    suiSample,
			want:    suiExpected,
		},
		{
			name:    "Sui - invalid",
			chainID: chains.SuiMainnet.ChainId,
			addr:    suiSample + "aa",
			wantErr: true,
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
		{"Zeta Mainnet", chains.ZetaChainMainnet.ChainId, true},
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
		{"Bitcoin Signet Testnet", chains.BitcoinSignetTestnet.ChainId, true},
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
	tt := []struct {
		name           string
		chainID        int64
		expectedParams *chaincfg.Params
		expectedError  require.ErrorAssertionFunc
	}{
		{
			name:           "Bitcoin Mainnet",
			chainID:        chains.BitcoinMainnet.ChainId,
			expectedParams: &chaincfg.MainNetParams,
			expectedError:  require.NoError,
		},
		{
			name:           "Bitcoin Testnet",
			chainID:        chains.BitcoinTestnet.ChainId,
			expectedParams: &chaincfg.TestNet3Params,
			expectedError:  require.NoError,
		},
		{
			name:           "Bitcoin Regtest",
			chainID:        chains.BitcoinRegtest.ChainId,
			expectedParams: &chaincfg.RegressionNetParams,
			expectedError:  require.NoError,
		},
		{
			name:           "Bitcoin Signet Testnet",
			chainID:        chains.BitcoinSignetTestnet.ChainId,
			expectedParams: &chaincfg.SigNetParams,
			expectedError:  require.NoError,
		},
		{
			name:           "Unknown Chain",
			chainID:        9999,
			expectedParams: nil,
			expectedError: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorContains(t, err, "error chainID 9999 is not a bitcoin chain")
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			params, err := chains.GetBTCChainParams(tc.chainID)
			tc.expectedError(t, err)
			require.Equal(t, tc.expectedParams, params)
		})
	}

}

func TestGetBTCChainIDFromChainParams(t *testing.T) {
	tt := []struct {
		name            string
		params          *chaincfg.Params
		expectedChainID int64
		expectedError   require.ErrorAssertionFunc
	}{
		{
			name:            "Bitcoin Mainnet",
			params:          &chaincfg.MainNetParams,
			expectedChainID: chains.BitcoinMainnet.ChainId,
			expectedError:   require.NoError,
		},
		{
			name:            "Bitcoin Testnet",
			params:          &chaincfg.TestNet3Params,
			expectedChainID: chains.BitcoinTestnet.ChainId,
			expectedError:   require.NoError,
		},
		{
			name:            "Bitcoin Regtest",
			params:          &chaincfg.RegressionNetParams,
			expectedChainID: chains.BitcoinRegtest.ChainId,
			expectedError:   require.NoError,
		},
		{
			name:            "Bitcoin Signet Testnet",
			params:          &chaincfg.SigNetParams,
			expectedChainID: chains.BitcoinSignetTestnet.ChainId,
			expectedError:   require.NoError,
		},
		{
			name:            "Bitcoin Testnet4",
			params:          &chains.TestNet4Params,
			expectedChainID: chains.BitcoinTestnet4.ChainId,
			expectedError:   require.NoError,
		},
		{
			name:            "Unknown Chain",
			params:          &chaincfg.Params{Name: "unknown"},
			expectedChainID: 0,
			expectedError: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorContains(t, err, "error chain unknown is not a bitcoin chain")
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			chainID, err := chains.GetBTCChainIDFromChainParams(tc.params)
			tc.expectedError(t, err)
			require.Equal(t, tc.expectedChainID, chainID)
		})
	}
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

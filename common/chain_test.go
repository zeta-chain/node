package common

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
)

func TestChain_EncodeAddress(t *testing.T) {
	tests := []struct {
		name    string
		chain   Chain
		b       []byte
		want    string
		wantErr bool
	}{
		{
			name: "should error if b is not a valid address on the network",
			chain: Chain{
				ChainName: ChainName_btc_testnet,
				ChainId:   18332,
			},
			b:       []byte("bc1qk0cc73p8m7hswn8y2q080xa4e5pxapnqgp7h9c"),
			want:    "",
			wantErr: true,
		},
		{
			name: "should pass if b is a valid address on the network",
			chain: Chain{
				ChainName: ChainName_btc_mainnet,
				ChainId:   8332,
			},
			b:       []byte("bc1qk0cc73p8m7hswn8y2q080xa4e5pxapnqgp7h9c"),
			want:    "bc1qk0cc73p8m7hswn8y2q080xa4e5pxapnqgp7h9c",
			wantErr: false,
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

func TestChain_InChainList(t *testing.T) {
	require.True(t, ZetaChainMainnet().InChainList(ZetaChainList()))
	require.True(t, ZetaMocknetChain().InChainList(ZetaChainList()))
	require.True(t, ZetaPrivnetChain().InChainList(ZetaChainList()))
	require.True(t, ZetaTestnetChain().InChainList(ZetaChainList()))
	require.False(t, EthChain().InChainList(ZetaChainList()))
}

func TestIsZetaChain(t *testing.T) {
	tests := []struct {
		name    string
		chainID int64
		want    bool
	}{
		{"Zeta Mainnet", ZetaChainMainnet().ChainId, true},
		{"Zeta Testnet", ZetaTestnetChain().ChainId, true},
		{"Zeta Mocknet", ZetaMocknetChain().ChainId, true},
		{"Zeta Privnet", ZetaPrivnetChain().ChainId, true},
		{"Non-Zeta", EthChain().ChainId, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, IsZetaChain(tt.chainID))
		})
	}
}

func TestIsEVMChain(t *testing.T) {
	tests := []struct {
		name    string
		chainID int64
		want    bool
	}{
		{"Ethereum Mainnet", EthChain().ChainId, true},
		{"Goerli Testnet", GoerliChain().ChainId, true},
		{"Sepolia Testnet", SepoliaChain().ChainId, true},
		{"Non-EVM", BtcMainnetChain().ChainId, false},
		{"Zeta Mainnet", ZetaChainMainnet().ChainId, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, IsEVMChain(tt.chainID))
		})
	}
}

func TestIsBitcoinChain(t *testing.T) {
	tests := []struct {
		name    string
		chainID int64
		want    bool
	}{
		{"Bitcoin Mainnet", BtcMainnetChain().ChainId, true},
		{"Bitcoin Testnet", BtcTestNetChain().ChainId, true},
		{"Bitcoin Regtest", BtcRegtestChain().ChainId, true},
		{"Non-Bitcoin", EthChain().ChainId, false},
		{"Zeta Mainnet", ZetaChainMainnet().ChainId, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, IsBitcoinChain(tt.chainID))
		})
	}
}

func TestIsEthereumChain(t *testing.T) {
	tests := []struct {
		name    string
		chainID int64
		want    bool
	}{
		{"Ethereum Mainnet", EthChain().ChainId, true},
		{"Goerli Testnet", GoerliChain().ChainId, true},
		{"Sepolia Testnet", SepoliaChain().ChainId, true},
		{"Non-Ethereum", BtcMainnetChain().ChainId, false},
		{"Zeta Mainnet", ZetaChainMainnet().ChainId, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, IsEthereumChain(tt.chainID))
		})
	}
}

func TestChain_IsExternalChain(t *testing.T) {
	require.False(t, ZetaChainMainnet().IsExternalChain())
	require.True(t, EthChain().IsExternalChain())
}

func TestChain_IsEmpty(t *testing.T) {
	require.True(t, Chain{}.IsEmpty())
	require.False(t, ZetaChainMainnet().IsEmpty())
}

func TestChains_Has(t *testing.T) {
	chains := Chains{ZetaChainMainnet(), ZetaTestnetChain()}
	require.True(t, chains.Has(ZetaChainMainnet()))
	require.False(t, chains.Has(EthChain()))
}

func TestChains_Distinct(t *testing.T) {
	chains := Chains{ZetaChainMainnet(), ZetaChainMainnet(), ZetaTestnetChain()}
	distinctChains := chains.Distinct()
	require.Len(t, distinctChains, 2)
}

func TestChains_Strings(t *testing.T) {
	chains := Chains{ZetaChainMainnet(), ZetaTestnetChain()}
	strings := chains.Strings()
	expected := []string{chains[0].String(), chains[1].String()}
	require.Equal(t, expected, strings)
}

func TestGetChainFromChainID(t *testing.T) {
	chain := GetChainFromChainID(ZetaChainMainnet().ChainId)
	require.Equal(t, ZetaChainMainnet(), *chain)
	require.Nil(t, GetChainFromChainID(9999))
}

func TestGetBTCChainParams(t *testing.T) {
	params, err := GetBTCChainParams(BtcMainnetChain().ChainId)
	require.NoError(t, err)
	require.Equal(t, &chaincfg.MainNetParams, params)

	_, err = GetBTCChainParams(9999)
	require.Error(t, err)
}

func TestGetBTCChainIDFromChainParams(t *testing.T) {
	chainID, err := GetBTCChainIDFromChainParams(&chaincfg.MainNetParams)
	require.NoError(t, err)
	require.Equal(t, int64(8332), chainID)

	_, err = GetBTCChainIDFromChainParams(&chaincfg.Params{Name: "unknown"})
	require.Error(t, err)
}

func TestChainIDInChainList(t *testing.T) {
	require.True(t, ChainIDInChainList(ZetaChainMainnet().ChainId, ZetaChainList()))
	require.False(t, ChainIDInChainList(EthChain().ChainId, ZetaChainList()))
}

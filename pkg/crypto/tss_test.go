package crypto

import (
	"testing"

	zetasui "github.com/zeta-chain/node/pkg/contracts/sui"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/cosmos"
)

func TestGetTSSAddrEVM(t *testing.T) {
	_, pubKey, _ := testdata.KeyTestPubAddr()
	pk, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey)
	require.NoError(t, err)
	decompresspubkey, err := crypto.DecompressPubkey(pubKey.Bytes())
	require.NoError(t, err)
	testCases := []struct {
		name      string
		tssPubkey string
		wantAddr  ethcommon.Address
		wantErr   bool
	}{
		{
			name:      "Valid TSS pubkey",
			tssPubkey: pk,
			wantAddr:  crypto.PubkeyToAddress(*decompresspubkey),
			wantErr:   false,
		},
		{
			name:      "Invalid TSS pubkey",
			tssPubkey: "invalid",
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			addr, err := GetTSSAddrEVM(tc.tssPubkey)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.Equal(t, tc.wantAddr, addr)
				require.NoError(t, err)
				require.NotEmpty(t, addr)
			}
		})
	}
}

func TestGetTSSAddrSui(t *testing.T) {
	_, pubKey, _ := testdata.KeyTestPubAddr()
	pk, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey)
	require.NoError(t, err)
	decompresspubkey, err := crypto.DecompressPubkey(pubKey.Bytes())
	require.NoError(t, err)
	testCases := []struct {
		name      string
		tssPubkey string
		wantAddr  string
		wantErr   bool
	}{
		{
			name:      "Valid TSS pubkey",
			tssPubkey: pk,
			wantAddr:  zetasui.AddressFromPubKeyECDSA(decompresspubkey),
			wantErr:   false,
		},
		{
			name:      "Invalid TSS pubkey",
			tssPubkey: "invalid",
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			addr, err := GetTSSAddrSui(tc.tssPubkey)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.Equal(t, tc.wantAddr, addr)
				require.NoError(t, err)
				require.NotEmpty(t, addr)
			}
		})
	}
}

func TestGetTSSAddrBTC(t *testing.T) {
	_, pubKey, _ := testdata.KeyTestPubAddr()
	pk, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey)
	require.NoError(t, err)
	testCases := []struct {
		name          string
		tssPubkey     string
		bitcoinParams *chaincfg.Params
		wantErr       bool
	}{
		{
			name:          "Valid TSS pubkey testnet params",
			tssPubkey:     pk,
			bitcoinParams: &chaincfg.TestNet3Params,
			wantErr:       false,
		},
		{
			name:          "Valid TSS pubkey signet params",
			tssPubkey:     pk,
			bitcoinParams: &chaincfg.SigNetParams,
			wantErr:       false,
		},
		{
			name:          "Invalid TSS pubkey testnet params",
			tssPubkey:     "invalid",
			bitcoinParams: &chaincfg.TestNet3Params,
			wantErr:       true,
		},
		{
			name:          "Valid TSS pubkey mainnet params",
			tssPubkey:     pk,
			bitcoinParams: &chaincfg.MainNetParams,
			wantErr:       false,
		},
		{
			name:          "Invalid TSS pubkey mainnet params",
			tssPubkey:     "invalid",
			bitcoinParams: &chaincfg.MainNetParams,
			wantErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			addr, err := GetTSSAddrBTC(tc.tssPubkey, tc.bitcoinParams)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				expectedAddr, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pubKey.Bytes()), tc.bitcoinParams)
				require.NoError(t, err)
				require.NotEmpty(t, addr)
				require.Equal(t, expectedAddr.EncodeAddress(), addr)
			}
		})
	}
}

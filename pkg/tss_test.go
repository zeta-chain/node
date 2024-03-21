package common

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common/cosmos"
)

func TestGetTssAddrEVM(t *testing.T) {
	_, pubKey, _ := testdata.KeyTestPubAddr()
	pk, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey)
	require.NoError(t, err)
	testCases := []struct {
		name      string
		tssPubkey string
		wantErr   bool
	}{
		{
			name:      "Valid TSS pubkey",
			tssPubkey: pk,
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
			addr, err := GetTssAddrEVM(tc.tssPubkey)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, addr)
			}
		})
	}
}

func TestGetTssAddrBTC(t *testing.T) {
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
			name:          "Valid TSS pubkey",
			tssPubkey:     pk,
			bitcoinParams: &chaincfg.TestNet3Params,
			wantErr:       false,
		},
		{
			name:          "Invalid TSS pubkey",
			tssPubkey:     "invalid",
			bitcoinParams: &chaincfg.TestNet3Params,
			wantErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			addr, err := GetTssAddrBTC(tc.tssPubkey, tc.bitcoinParams)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, addr)
			}
		})
	}
}

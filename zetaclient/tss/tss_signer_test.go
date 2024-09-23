package tss

import (
	"fmt"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/cmd"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/cosmos"
	"github.com/zeta-chain/node/pkg/crypto"
	"github.com/zeta-chain/node/zetaclient/testutils"
)

func setupConfig() {
	testConfig := sdk.GetConfig()
	testConfig.SetBech32PrefixForAccount(cmd.Bech32PrefixAccAddr, cmd.Bech32PrefixAccPub)
	testConfig.SetBech32PrefixForValidator(cmd.Bech32PrefixValAddr, cmd.Bech32PrefixValPub)
	testConfig.SetBech32PrefixForConsensusNode(cmd.Bech32PrefixConsAddr, cmd.Bech32PrefixConsPub)
	testConfig.SetFullFundraiserPath(cmd.ZetaChainHDPath)
	sdk.SetCoinDenomRegex(func() string {
		return cmd.DenomRegex
	})
}

func Test_LoadTssFilesFromDirectory(t *testing.T) {

	tt := []struct {
		name string
		n    int
	}{
		{
			name: "2 keyshare files",
			n:    2,
		},
		{
			name: "10 keyshare files",
			n:    10,
		},
		{
			name: "No keyshare files",
			n:    0,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tempdir, err := os.MkdirTemp("", "test-tss")
			require.NoError(t, err)
			err = GenerateKeyshareFiles(tc.n, tempdir)
			require.NoError(t, err)
			tss := TSS{
				logger:        zerolog.New(os.Stdout),
				Keys:          map[string]*Key{},
				CurrentPubkey: "",
			}
			err = tss.LoadTssFilesFromDirectory(tempdir)
			require.NoError(t, err)
			require.Equal(t, tc.n, len(tss.Keys))
		})
	}
}

func GenerateKeyshareFiles(n int, dir string) error {
	setupConfig()
	err := os.Chdir(dir)
	if err != nil {
		return err
	}
	for i := 0; i < n; i++ {
		_, pubKey, _ := testdata.KeyTestPubAddr()
		spk, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey)
		if err != nil {
			return err
		}
		pk, err := crypto.NewPubKey(spk)
		if err != nil {
			return err
		}
		filename := fmt.Sprintf("localstate-%s", pk.String())
		b, err := pk.MarshalJSON()
		if err != nil {
			return err
		}
		err = os.WriteFile(filename, b, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func Test_EVMAddress(t *testing.T) {
	setupConfig()

	// test Athens3 tss pubkey
	tss := TSS{
		CurrentPubkey: testutils.TSSPubkeyAthens3,
	}
	evmAddr := tss.EVMAddress()
	require.Equal(t, testutils.TSSAddressEVMAthens3, evmAddr.String())
}

func Test_BTCAddress(t *testing.T) {
	setupConfig()

	tests := []struct {
		name       string
		tssPubkey  string
		btcChainID int64
		wantAddr   string
	}{
		{
			name:       "Athens3 tss pubkey",
			tssPubkey:  testutils.TSSPubkeyAthens3,
			btcChainID: chains.BitcoinTestnet.ChainId,
			wantAddr:   testutils.TSSAddressBTCAthens3,
		},
		{
			name:       "local network tss pubkey",
			tssPubkey:  "zetapub1addwnpepqdax2apf4qmqcaxzae7t4m9xz76mungtppsyw5shvznd52ldy6sjjsjfa3z",
			btcChainID: chains.BitcoinRegtest.ChainId,
			wantAddr:   "bcrt1q30ew8md3rd9fx6n4qx0a9tmz0mz44lzjwxppnu",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tss := TSS{
				CurrentPubkey: tc.tssPubkey,
			}
			btcAddr := tss.BTCAddress(tc.btcChainID)
			require.Equal(t, tc.wantAddr, btcAddr)
		})
	}

}

func Test_ValidateAddresses(t *testing.T) {
	setupConfig()

	tests := []struct {
		name        string
		tssPubkey   string
		btcChainIDs []int64
		errMsg      string
	}{
		{
			name:      "Validation success",
			tssPubkey: testutils.TSSPubkeyAthens3,
			btcChainIDs: []int64{
				chains.BitcoinTestnet.ChainId,
				chains.BitcoinSignetTestnet.ChainId,
			},
			errMsg: "",
		},
		{
			name:        "Validation failed on EVM address",
			tssPubkey:   "invalidpubkey", // to make EVMAddress() failed
			btcChainIDs: []int64{},
			errMsg:      "blank tss evm address",
		},
		{
			name:      "Validation failed on BTC address",
			tssPubkey: testutils.TSSPubkeyAthens3,
			btcChainIDs: []int64{
				chains.BitcoinTestnet.ChainId,
				100, // unknown BTC chain ID
			},
			errMsg: "cannot derive btc address",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tss := TSS{
				CurrentPubkey: tc.tssPubkey,
			}
			err := tss.ValidateAddresses(tc.btcChainIDs)
			if tc.errMsg != "" {
				require.Contains(t, err.Error(), tc.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

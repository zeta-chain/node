package tss

import (
	"fmt"
	"os"
	"testing"

	"github.com/zeta-chain/zetacore/zetaclient/keys"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/cosmos"
)

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
			assert.NoError(t, err)
			err = GenerateKeyshareFiles(tc.n, tempdir)
			assert.NoError(t, err)
			tss := TSS{
				logger:        zerolog.New(os.Stdout),
				Keys:          map[string]*Key{},
				CurrentPubkey: "",
			}
			err = tss.LoadTssFilesFromDirectory(tempdir)
			assert.NoError(t, err)
			assert.Equal(t, tc.n, len(tss.Keys))
		})
	}
}

func GenerateKeyshareFiles(n int, dir string) error {
	keys.SetupConfigForTest()
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
		pk, err := common.NewPubKey(spk)
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

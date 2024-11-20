package tss

import (
	"fmt"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/cosmos"
	"github.com/zeta-chain/node/pkg/crypto"
)

func Test_ParsePubKeysFromPath(t *testing.T) {
	for _, tt := range []struct {
		name string
		n    int
	}{
		{name: "2 keyshare files", n: 2},
		{name: "10 keyshare files", n: 10},
		{name: "No keyshare files", n: 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			logger := zerolog.New(zerolog.NewTestWriter(t))

			dir, err := os.MkdirTemp("", "test-tss")
			require.NoError(t, err)

			generateKeyShareFiles(t, tt.n, dir)

			// ACT
			keys, err := ParsePubKeysFromPath(dir, logger)

			// ASSERT
			require.NoError(t, err)
			require.Equal(t, tt.n, len(keys))
		})
	}
}

func generateKeyShareFiles(t *testing.T, n int, dir string) {
	err := os.Chdir(dir)
	require.NoError(t, err)
	for i := 0; i < n; i++ {
		_, pubKey, _ := testdata.KeyTestPubAddr()

		spk, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey)
		require.NoError(t, err)

		pk, err := crypto.NewPubKey(spk)
		require.NoError(t, err)

		b, err := pk.MarshalJSON()
		require.NoError(t, err)

		filename := fmt.Sprintf("localstate-%s.json", pk.String())

		err = os.WriteFile(filename, b, 0644)
		require.NoError(t, err)
	}
}

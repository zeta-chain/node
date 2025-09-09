package tss

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
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

			dir := t.TempDir()

			generateKeyShareFiles(t, tt.n, dir)

			// ACT
			keys, err := ParsePubKeysFromPath(logger, dir)

			// ASSERT
			require.NoError(t, err)
			require.Equal(t, tt.n, len(keys))
		})
	}
}

func Test_ResolvePreParamsFromPath(t *testing.T) {
	t.Run("file not found", func(t *testing.T) {
		// ARRANGE
		path := filepath.Join(os.TempDir(), "hello-123.json")

		// ACT
		_, err := ResolvePreParamsFromPath(path)

		// ASSERT
		require.Error(t, err)
		require.Contains(t, err.Error(), "unable to read pre-params")
	})

	t.Run("invalid file", func(t *testing.T) {
		// ARRANGE
		tmpFile, err := os.CreateTemp(os.TempDir(), "pre-params-*.json")
		require.NoError(t, err)
		t.Cleanup(func() {
			require.NoError(t, os.Remove(tmpFile.Name()))
		})

		_, err = tmpFile.WriteString(`invalid-json`)
		require.NoError(t, err)
		tmpFile.Close()

		// ACT
		_, err = ResolvePreParamsFromPath(tmpFile.Name())

		// ASSERT
		require.Error(t, err)
		require.Contains(t, err.Error(), "unable to decode pre-params")
	})

	t.Run("AllGood", func(t *testing.T) {
		// ARRANGE
		tmpFile, err := os.CreateTemp(os.TempDir(), "pre-params-*.json")
		require.NoError(t, err)
		t.Cleanup(func() {
			require.NoError(t, os.Remove(tmpFile.Name()))
		})

		createPreParams(t, tmpFile.Name())

		// ACT
		resolvedPreParams, err := ResolvePreParamsFromPath(tmpFile.Name())

		// Assert
		require.NoError(t, err)
		require.NotNil(t, resolvedPreParams)
	})
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

//go:embed testdata/pre-params.json
var preParamsFixture []byte

// createPreParams creates a pre-params file at the given path.
// uses fixture to skip long setup.
func createPreParams(t *testing.T, filePath string) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0600)
	require.NoError(t, err)

	_, err = file.Write(preParamsFixture)
	require.NoError(t, err)
	require.NoError(t, file.Close())
}

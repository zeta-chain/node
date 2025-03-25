package ton

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/wallet"
)

func TestWalletConstruction(t *testing.T) {
	// ARRANGE
	seed := wallet.RandomSeed()

	pk, err := wallet.SeedToPrivateKey(seed)
	require.NoError(t, err)

	t.Logf("seed[ %s ] ==> privateKey(0x%s)", seed, hex.EncodeToString(pk.Seed()))

	// ACT
	accInit, w, err := ConstructWalletFromPrivateKey(pk, nil)

	// ASSERT
	require.NoError(t, err)
	require.NotNil(t, accInit)
	require.NotNil(t, w)
}

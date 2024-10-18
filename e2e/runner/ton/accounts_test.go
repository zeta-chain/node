package ton

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/wallet"
)

func TestWalletConstruction(t *testing.T) {
	// ARRANGE
	seed := wallet.RandomSeed()

	// ACT
	accInit, w, err := ConstructWalletFromSeed(seed, nil)

	// ASSERT
	require.NoError(t, err)
	require.NotNil(t, accInit)
	require.NotNil(t, w)
}

package ton

import (
	"crypto/ecdsa"
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/wallet"
)

func TestGateway(t *testing.T) {
	// ARRANGE
	// Given TSS address
	const sampleTSSPrivateKey = "0xb984cd65727cfd03081fc7bf33bf5c208bca697ce16139b5ded275887e81395a"

	pkBytes, err := hex.DecodeString(sampleTSSPrivateKey[2:])
	require.NoError(t, err)

	privateKey, err := crypto.ToECDSA(pkBytes)
	require.NoError(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)

	tss := crypto.PubkeyToAddress(*publicKeyECDSA)

	t.Run("Build gateway deployment payload", func(t *testing.T) {
		// ACT
		codeCell, errCode := getGatewayCode()
		stateCell, errState := buildGatewayData(tss)

		// ASSERT
		require.NoError(t, errCode)
		require.NoError(t, errState)

		codeString, err := codeCell.ToBocStringCustom(false, true, false, 0)
		require.NoError(t, err)

		stateString, err := stateCell.ToBocStringCustom(false, true, false, 0)
		require.NoError(t, err)

		t.Logf("Gateway code: %s", codeString)
		t.Logf("Gateway state: %s", stateString)

		// Taken from jest tests in protocol-contracts-ton (using the same TSS address private key)
		const expectedState = "b5ee9c7241010101001c0000338000000000124d38a790fdf1d9311fae87d4b21aeffd77bc26c0776433f3"

		require.Equal(t, expectedState, stateString)
	})
}

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

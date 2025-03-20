package ton

import (
	"context"
	"encoding/hex"
	"os"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/wallet"
	"golang.org/x/crypto/ed25519"
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

func TestWIP(t *testing.T) {
	// ARRANGE
	// secret from GH (solana_private_key: "${solana_private_key}")
	b58 := os.Getenv("SOLA_PRIVATE_KEY")
	require.NotEmpty(t, b58, "SOLA_PRIVATE_KEY env var not set")

	// Convert to ed25519 private key
	decoded := base58.Decode(b58)
	pk := ed25519.PrivateKey(decoded)

	// Create a client
	c, err := liteapi.NewClientWithDefaultTestnet()
	require.NoError(t, err)

	// Create a wallet
	w, err := wallet.New(pk, wallet.V5R1, c)
	require.NoError(t, err)

	// Ensure address matches testnet user
	require.Equal(t, "0:48097c7d2cab0ff3bb7bdda7689f25774431e48c73b3ec4761ad1248cc778c29", w.GetAddress().ToRaw())

	// ACT
	// Send 0.1 TON to someone.
	// This will internally trigger the logic for contract deployment & activation
	_, err = w.SendV2(context.Background(), 10*time.Second, wallet.SimpleTransfer{
		// 0.1 TON
		Amount: 100_000_000,
		// Dmitry's test wallet
		Address: ton.MustParseAccountID("0:74a36900b786949a60c95ee20a56e583f908f2e957f3ffcb1e9770cc9edd408d"),
		Comment: "hey",
	})

	// ASSERT
	require.NoError(t, err)
}

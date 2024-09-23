package ton

import (
	"embed"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

func TestParsing(t *testing.T) {
	t.Run("Donate", func(t *testing.T) {
		// ARRANGE
		// Given a tx
		tx, fx := getFixtureTX(t, "00-donation")

		// Given a gateway contract
		gw := NewGateway(ton.MustParseAccountID(fx.Account))

		// ACT
		parsedTX, err := gw.ParseTransaction(tx)

		// ASSERT
		require.NoError(t, err)

		assert.Equal(t, int(OpDonate), int(parsedTX.Operation))
		assert.Equal(t, true, parsedTX.IsInbound())

		const (
			expectedSender   = "0:9594c719ec4c95f66683b2fb1ca0b09de4a41f6fb087ba4c8d265b96a4cce50f"
			expectedDonation = 1_499_432_947 // 1.49... TON
		)

		donation, err := parsedTX.Donation()
		assert.NoError(t, err)
		assert.Equal(t, expectedSender, donation.Sender.ToRaw())
		assert.Equal(t, expectedDonation, int(donation.Amount.Uint64()))
	})

	t.Run("Deposit", func(t *testing.T) {
		// ARRANGE
		// Given a tx
		tx, fx := getFixtureTX(t, "01-deposit")

		// Given a gateway contract
		gw := NewGateway(ton.MustParseAccountID(fx.Account))

		// ACT
		parsedTX, err := gw.ParseTransaction(tx)

		// ASSERT
		require.NoError(t, err)

		// Check tx props
		assert.Equal(t, int(OpDeposit), int(parsedTX.Operation))

		// Check
		deposit, err := parsedTX.Deposit()
		assert.NoError(t, err)

		const (
			expectedSender  = "0:9594c719ec4c95f66683b2fb1ca0b09de4a41f6fb087ba4c8d265b96a4cce50f"
			vitalikDotETH   = "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"
			expectedDeposit = 990_000_000 // 0.99 TON
		)

		assert.Equal(t, expectedSender, deposit.Sender.ToRaw())
		assert.Equal(t, expectedDeposit, int(deposit.Amount.Uint64()))
		assert.Equal(t, vitalikDotETH, deposit.Recipient.Hex())

		// Check that other casting fails
		_, err = parsedTX.Donation()
		assert.ErrorIs(t, err, ErrCast)
	})

	t.Run("Deposit and call", func(t *testing.T) {
		// todo
	})

	t.Run("Irrelevant tx", func(t *testing.T) {
		// todo
	})

	// todo
}

func TestFixtures(t *testing.T) {
	// ACT
	tx, _ := getFixtureTX(t, "01-deposit")

	// ASSERT
	require.Equal(t, uint64(26023788000003), tx.Lt)
	require.Equal(t, "cbd6e2261334d08120e2fef428ecbb4e7773606ced878d0e6da204f2b4bf42bf", tx.Hash().Hex())
}

//go:embed testdata
var fixtures embed.FS

type fixture struct {
	Account     string `json:"account"`
	BOC         string `json:"boc"`
	Description string `json:"description"`
	Hash        string `json:"hash"`
	LogicalTime uint64 `json:"logicalTime"`
	Test        bool   `json:"test"`
}

// testdata/$name.json tx
func getFixtureTX(t *testing.T, name string) (ton.Transaction, fixture) {
	t.Helper()

	filename := fmt.Sprintf("testdata/%s.json", name)

	b, err := fixtures.ReadFile(filename)
	require.NoError(t, err, filename)

	// bag of cells
	var fx fixture

	require.NoError(t, json.Unmarshal(b, &fx))

	cells, err := boc.DeserializeBocHex(fx.BOC)
	require.NoError(t, err)
	require.Len(t, cells, 1)

	cell := cells[0]

	var tx ton.Transaction

	require.NoError(t, tx.UnmarshalTLB(cell, &tlb.Decoder{}))

	t.Logf("Loaded fixture %s\n%s", filename, fx.Description)

	return tx, fx
}

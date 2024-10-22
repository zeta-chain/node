package ton

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

func TestParsing(t *testing.T) {
	swapBodyAndParse := func(gw *Gateway, tx ton.Transaction, body *boc.Cell) *Transaction {
		tx.Msgs.InMsg.Value.Value.Body.Value = tlb.Any(*body)

		parsed, err := gw.ParseTransaction(tx)
		require.NoError(t, err)

		return parsed
	}

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

		// Check that AsBody works
		var (
			parsedTX2 = swapBodyAndParse(gw, tx, lo.Must(donation.AsBody()))
			donation2 = lo.Must(parsedTX2.Donation())
		)

		assert.Equal(t, donation, donation2)
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

		// Check deposit
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

		// Check that AsBody works
		var (
			parsedTX2 = swapBodyAndParse(gw, tx, lo.Must(deposit.AsBody()))
			deposit2  = lo.Must(parsedTX2.Deposit())
		)

		assert.Equal(t, deposit, deposit2)
	})

	t.Run("Deposit and call", func(t *testing.T) {
		// ARRANGE
		// Given a tx
		tx, fx := getFixtureTX(t, "02-deposit-and-call")

		// Given a gateway contract
		gw := NewGateway(ton.MustParseAccountID(fx.Account))

		// ACT
		parsedTX, err := gw.ParseTransaction(tx)

		// ASSERT
		require.NoError(t, err)

		// Check tx props
		assert.Equal(t, int(OpDepositAndCall), int(parsedTX.Operation))

		// Check deposit and call
		depositAndCall, err := parsedTX.DepositAndCall()
		assert.NoError(t, err)

		const (
			expectedSender  = "0:9594c719ec4c95f66683b2fb1ca0b09de4a41f6fb087ba4c8d265b96a4cce50f"
			vitalikDotETH   = "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"
			expectedDeposit = 490_000_000 // 0.49 TON
		)

		expectedCallData := readFixtureFile(t, "testdata/long-call-data.txt")

		assert.Equal(t, expectedSender, depositAndCall.Sender.ToRaw())
		assert.Equal(t, expectedDeposit, int(depositAndCall.Amount.Uint64()))
		assert.Equal(t, vitalikDotETH, depositAndCall.Recipient.Hex())
		assert.Equal(t, expectedCallData, depositAndCall.CallData)

		// Check that AsBody works
		var (
			parsedTX2       = swapBodyAndParse(gw, tx, lo.Must(depositAndCall.AsBody()))
			depositAndCall2 = lo.Must(parsedTX2.DepositAndCall())
		)

		assert.Equal(t, depositAndCall, depositAndCall2)
	})

	t.Run("Irrelevant tx", func(t *testing.T) {
		t.Run("Failed tx", func(t *testing.T) {
			// ARRANGE
			// Given a tx
			tx, fx := getFixtureTX(t, "03-failed-tx")

			// Given a gateway contract
			gw := NewGateway(ton.MustParseAccountID(fx.Account))

			// ACT
			_, err := gw.ParseTransaction(tx)

			assert.ErrorIs(t, err, ErrParse)

			// 102 is 'unknown op'
			// https://github.com/zeta-chain/protocol-contracts-ton/blob/main/contracts/common/errors.fc
			assert.ErrorContains(t, err, "is not successful (exit code 102)")
		})

		t.Run("not a deposit nor withdrawal", func(t *testing.T) {
			// actually, it's a bounce of the previous tx

			// ARRANGE
			// Given a tx
			tx, fx := getFixtureTX(t, "04-bounced-msg")

			// Given a gateway contract
			gw := NewGateway(ton.MustParseAccountID(fx.Account))

			// ACT
			_, err := gw.ParseTransaction(tx)
			assert.Error(t, err)
		})
	})
}

func TestFiltering(t *testing.T) {
	t.Run("Inbound", func(t *testing.T) {
		for _, tt := range []struct {
			name  string
			skip  bool
			error bool
		}{
			// Should be parsed and filtered
			{"00-donation", false, false},
			{"01-deposit", false, false},
			{"02-deposit-and-call", false, false},

			// Should be skipped
			{"03-failed-tx", true, false},
			{"04-bounced-msg", true, false},
		} {
			t.Run(tt.name, func(t *testing.T) {
				// ARRANGE
				// Given a tx
				tx, fx := getFixtureTX(t, tt.name)

				// Given a gateway
				gw := NewGateway(ton.MustParseAccountID(fx.Account))

				// ACT
				parsedTX, skip, err := gw.ParseAndFilter(tx, FilterInbounds)

				if tt.error {
					require.Error(t, err)
					assert.False(t, skip)
					assert.Nil(t, parsedTX)
					return
				}

				require.NoError(t, err)
				assert.Equal(t, tt.skip, skip)

				if tt.skip {
					assert.Nil(t, parsedTX)
					return
				}

				assert.NotNil(t, parsedTX)
			})
		}
	})
}

func TestFixtures(t *testing.T) {
	// ACT
	tx, _ := getFixtureTX(t, "01-deposit")

	// ASSERT
	require.Equal(t, uint64(26023788000003), tx.Lt)
	require.Equal(t, "cbd6e2261334d08120e2fef428ecbb4e7773606ced878d0e6da204f2b4bf42bf", tx.Hash().Hex())
}

func TestSnakeData(t *testing.T) {
	for _, tt := range []string{
		"Hello world",
		"123",
		strings.Repeat(`ZetaChain `, 300),
		string(readFixtureFile(t, "testdata/long-call-data.txt")),
	} {
		a := []byte(tt)

		cell, err := MarshalSnakeCell(a)
		require.NoError(t, err)

		b, err := UnmarshalSnakeCell(cell)
		require.NoError(t, err)

		t.Logf(string(b))

		assert.Equal(t, a, b, tt)
	}
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

	var (
		filename = fmt.Sprintf("testdata/%s.json", name)
		b        = readFixtureFile(t, filename)
	)

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

func readFixtureFile(t *testing.T, filename string) []byte {
	t.Helper()

	b, err := fixtures.ReadFile(filename)
	require.NoError(t, err, filename)

	return b
}

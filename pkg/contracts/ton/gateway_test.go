package ton

import (
	"crypto/ecdsa"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

func TestParsing(t *testing.T) {
	// small helpers that allows to alter inMsg body and parse it again
	alterBodyAndParse := func(gw *Gateway, tx ton.Transaction, body *boc.Cell) *Transaction {
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
			expectedDonation = 599_509_877 // ~0.6 TON
		)

		donation, err := parsedTX.Donation()
		assert.NoError(t, err)
		assert.Equal(t, expectedSender, donation.Sender.ToRaw())
		assert.Equal(t, expectedDonation, int(donation.Amount.Uint64()))

		// Check that AsBody works
		var (
			parsedTX2 = alterBodyAndParse(gw, tx, lo.Must(donation.AsBody()))
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
			zevmRecipient   = "0xA1eb8D65b765D259E7520B791bc4783AdeFDd998"
			expectedDeposit = 996_000_000 // 0.996 TON
		)

		assert.Equal(t, expectedSender, deposit.Sender.ToRaw())
		assert.Equal(t, expectedDeposit, int(deposit.Amount.Uint64()))
		assert.Equal(t, zevmRecipient, deposit.Recipient.Hex())

		// Check that other casting fails
		_, err = parsedTX.Donation()
		assert.ErrorIs(t, err, ErrCast)

		// Check that AsBody works
		var (
			parsedTX2 = alterBodyAndParse(gw, tx, lo.Must(deposit.AsBody()))
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
			expectedSender   = "0:9594c719ec4c95f66683b2fb1ca0b09de4a41f6fb087ba4c8d265b96a4cce50f"
			zevmRecipient    = "0xA1eb8D65b765D259E7520B791bc4783AdeFDd998"
			expectedDeposit  = 394_800_000 // 0.4 TON - tx fee
			expectedCallData = `These "NFTs" - are they in the room with us right now?`
		)

		assert.Equal(t, expectedSender, depositAndCall.Sender.ToRaw())
		assert.Equal(t, expectedDeposit, int(depositAndCall.Amount.Uint64()))
		assert.Equal(t, zevmRecipient, depositAndCall.Recipient.Hex())
		assert.Equal(t, []byte(expectedCallData), depositAndCall.CallData)

		t.Run("Long call data", func(t *testing.T) {
			// ARRANGE
			longCallData := readFixtureFile(t, "testdata/long-call-data.txt")
			depositAndCall.CallData = longCallData

			// ACT
			parsedTX = alterBodyAndParse(gw, tx, lo.Must(depositAndCall.AsBody()))

			depositAndCall2, err := parsedTX.DepositAndCall()

			// ASSERT
			require.NoError(t, err)

			assert.Equal(t, longCallData, depositAndCall2.CallData)
		})

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

		t.Run("not a deposit nor a withdrawal", func(t *testing.T) {
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
	require.Equal(t, uint64(27040013000003), tx.Lt)
	require.Equal(t, "f893d7ed7fc3d73aedb44ca7c350026a5d27e679cf85c0c8df9e69db28387b06", tx.Hash().Hex())
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

func TestDeployment(t *testing.T) {
	// ARRANGE
	// Given TSS address & Authority address
	const (
		sampleTSSPrivateKey = "0xb984cd65727cfd03081fc7bf33bf5c208bca697ce16139b5ded275887e81395a"
		sampleAuthority     = "0:4686a2c066c784a915f3e01c853d3195ed254c948e21adbb3e4a9b3f5f3c74d7"
	)

	pkBytes, err := hex.DecodeString(sampleTSSPrivateKey[2:])
	require.NoError(t, err)

	privateKey, err := crypto.ToECDSA(pkBytes)
	require.NoError(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)

	tss := crypto.PubkeyToAddress(*publicKeyECDSA)

	// ACT
	code := GatewayCode()
	stateInit := GatewayStateInit(ton.MustParseAccountID(sampleAuthority), tss, true)

	// ASSERT
	codeString, err := code.ToBocStringCustom(false, true, false, 0)
	require.NoError(t, err)

	stateString, err := stateInit.ToBocStringCustom(false, true, false, 0)
	require.NoError(t, err)

	t.Logf("Gateway code: %s", codeString)
	t.Logf("Gateway state: %s", stateString)

	// Taken from jest tests in protocol-contracts-ton (using the same vars for initState config)
	const expectedState = "b5ee9c7241010101003c000074800000000124d38a790fdf1d9311fae87d4b21aeffd77bc26c004686a2c066c784a915f3e01c853d3195ed254c948e21adbb3e4a9b3f5f3c74d746f17671"

	require.Equal(t, expectedState, stateString)
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

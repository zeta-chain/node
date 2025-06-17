package ton

import (
	"crypto/ecdsa"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi"
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
		assert.Zero(t, parsedTX.ExitCode)

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

	t.Run("Deposit and call ABI", func(t *testing.T) {
		// ARRANGE
		// Given a tx
		tx, fx := getFixtureTX(t, "07-deposit-and-call-abi")

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
			expectedSender  = "0:74a36900b786949a60c95ee20a56e583f908f2e957f3ffcb1e9770cc9edd408d"
			zevmRecipient   = "0x13ad4f89050E83e8F485BB8349b40d3b89833790"
			expectedDeposit = 94_800_000 // 0.1 TON - tx fee

			// cast abi-encode "fn(string)" "ZETA ON TON YEAH"
			callData = "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000105a455441204f4e20544f4e205945414800000000000000000000000000000000"
		)

		assert.Equal(t, expectedSender, depositAndCall.Sender.ToRaw())
		assert.Equal(t, expectedDeposit, int(depositAndCall.Amount.Uint64()))
		assert.Equal(t, zevmRecipient, depositAndCall.Recipient.Hex())
		assert.Equal(t, callData, "0x"+hex.EncodeToString(depositAndCall.CallData))
	})

	t.Run("Call ABI", func(t *testing.T) {
		// ARRANGE
		// Given a tx
		tx, fx := getFixtureTX(t, "08-call-abi")

		// Given a gateway contract
		gw := NewGateway(ton.MustParseAccountID(fx.Account))

		// ACT
		parsedTX, err := gw.ParseTransaction(tx)

		// ASSERT
		require.NoError(t, err)

		// Check tx props
		assert.Equal(t, int(OpCall), int(parsedTX.Operation))

		// Check call
		call, err := parsedTX.Call()
		assert.NoError(t, err)

		// Check that call data is valid and also is ABI encoded
		// cast abi-encode "fn(string)" "ping-pong"
		const abiDef = `[{"name":"fn", "type":"function", "inputs":[{"type":"string"}]}]`

		parsedABI, err := abi.JSON(strings.NewReader(abiDef))
		require.NoError(t, err)

		args, err := parsedABI.Methods["fn"].Inputs.Unpack(call.CallData)
		require.NoError(t, err)
		require.Len(t, args, 1)

		assert.Equal(t, "ping-pong", args[0].(string))
		assert.Equal(t, "0:74a36900b786949a60c95ee20a56e583f908f2e957f3ffcb1e9770cc9edd408d", call.Sender.ToRaw())
		assert.Equal(t, "0xA1eb8D65b765D259E7520B791bc4783AdeFDd998", call.Recipient.Hex())
	})

	t.Run("Withdrawal", func(t *testing.T) {
		// ARRANGE
		// Given a tx
		tx, fx := getFixtureTX(t, "06-withdrawal")

		// Given a gateway contract
		gw := NewGateway(ton.MustParseAccountID(fx.Account))

		// ACT
		parsedTX, err := gw.ParseTransaction(tx)

		// ASSERT
		require.NoError(t, err)
		assert.Equal(t, OpWithdraw, parsedTX.Operation)
		assert.Zero(t, parsedTX.ExitCode)

		// 0.01 TON
		const expectedAmount = 10_000_000

		actualAmount := parsedTX.Msgs.OutMsgs.Values()[0].Value.Info.IntMsgInfo.Value.Grams
		assert.Equal(t, tlb.Coins(expectedAmount), actualAmount)

		// Check withdrawal
		w, err := parsedTX.Withdrawal()
		require.NoError(t, err)

		expectedRecipient := ton.MustParseAccountID("0QDQ51yWafHKgOjMF4ZwOfGsxnQ2yx_InZQG1NGwMfs2y62E")

		assert.Equal(t, expectedRecipient, w.Recipient)
		assert.Equal(t, uint32(1), w.Seqno)
		assert.Equal(t, math.NewUint(expectedAmount), w.Amount)

		// Expect TSS signer address
		signer, err := w.Signer()
		require.NoError(t, err)

		const expectedTSS = "0xFA033cebd2EB4A800F74d70C10dfc8710fF0d148"
		assert.Equal(t, expectedTSS, signer.Hex())
	})

	t.Run("IncreaseSeqno", func(t *testing.T) {
		// ARRANGE
		// Given a tx
		tx, fx := getFixtureTX(t, "09-increase-seqno")

		// Given a gateway contract
		gw := NewGateway(ton.MustParseAccountID(fx.Account))

		// ACT
		parsedTX, err := gw.ParseTransaction(tx)

		// ASSERT
		require.NoError(t, err)
		assert.Equal(t, OpIncreaseSeqno, parsedTX.Operation)
		assert.Zero(t, parsedTX.ExitCode)

		// Check withdrawal
		is, err := parsedTX.IncreaseSeqno()
		require.NoError(t, err)

		// Expect reason code
		require.Equal(t, uint32(1337), is.ReasonCode)

		// Expect TSS signer address
		signer, err := is.Signer()
		require.NoError(t, err)

		const expectedTSS = "0xFA033cebd2EB4A800F74d70C10dfc8710fF0d148"
		assert.Equal(t, expectedTSS, signer.Hex())
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

			assert.ErrorIs(t, err, ErrUnknownOp)

			// 102 is 'unknown op'
			// https://github.com/zeta-chain/protocol-contracts-ton/blob/main/contracts/common/errors.fc
			assert.ErrorContains(t, err, "unknown op")
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

func TestWithdrawal(t *testing.T) {
	// FYI: asserts are checks with protocol-contracts-ton tests (Typescript lib)

	// ARRANGE
	// Given a withdrawal msg
	withdrawal := &Withdrawal{
		Recipient: ton.MustParseAccountID("0:552f6db5da0cae7f0b3ab4ab58d85927f6beb962cda426a6a6ee751c82cead1f"),
		Amount:    Coins(5),
		Seqno:     2,
	}

	const expectedHash = "e8eddf26276c747bd14d5161d18bc235c5c1a050187ab468996572d34e2f8f30"

	// ACT
	hash, err := withdrawal.Hash()

	// ASSERT
	require.NoError(t, err)
	require.Equal(t, expectedHash, fmt.Sprintf("%x", hash[:]))

	t.Run("Sign", func(t *testing.T) {
		// ARRANGE
		// Given a sample EVM wallet (simulates TSS)
		privateKey := evmWallet(t, "0xb984cd65727cfd03081fc7bf33bf5c208bca697ce16139b5ded275887e81395a")

		// Given ECDSA signature ...
		sig, err := crypto.Sign(hash[:], privateKey)
		require.NoError(t, err)

		var sigArray [65]byte
		copy(sigArray[:], sig)

		// That is set in withdrawal
		withdrawal.SetSignature(sigArray)

		// ACT
		// Convert withdrawal to external-message's body
		body, err := withdrawal.AsBody()
		require.NoError(t, err)

		// ASSERT
		// Check that sig is not empty
		require.False(t, withdrawal.emptySig())

		// Ensure that signature has the right format when decoding
		var (
			r = big.NewInt(0).SetBytes(sig[:32])
			s = big.NewInt(0).SetBytes(sig[32:64])
			v = sig[64]
		)

		actualV, err := body.ReadUint(8)
		require.NoError(t, err)

		actualR, err := body.ReadBigUint(256)
		require.NoError(t, err)

		actualS, err := body.ReadBigUint(256)
		require.NoError(t, err)

		assert.Equal(t, 0, r.Cmp(actualR))
		assert.Equal(t, 0, s.Cmp(actualS))
		assert.Equal(t, uint64(v), actualV)
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
	h2b := func(raw string) []byte {
		b, err := hex.DecodeString(strings.TrimPrefix(raw, "0x"))
		require.NoError(t, err)

		return b
	}

	for _, tt := range []struct {
		name string
		data []byte
	}{
		{name: "simple", data: []byte("123")},
		{name: "hello world", data: []byte("Hello world")},
		{name: "zeta repeated 10 times", data: []byte(strings.Trim(strings.Repeat("ZetaChain ", 10), " "))},
		{name: "zeta repeated 50 times", data: []byte(strings.Trim(strings.Repeat("ZetaChain ", 50), " "))},
		{name: "long call data text", data: readFixtureFile(t, "testdata/long-call-data.txt")},
		{
			// cast abi-encode "fn(string)" "ping-pong"
			name: "abi call data 1",
			data: h2b("0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000970696e672d706f6e670000000000000000000000000000000000000000000000"),
		},
		{
			// cast abi-encode "swapExactTokensForTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline)" 123 456 "[0x0E0E08C73CD5019d6CAc98311Fa18edf98b70428,0x0E0E08C73CD5019d6CAc98311Fa18edf98b70428]" 0x0E0E08C73CD5019d6CAc98311Fa18edf98b70428 999
			name: "abi call data 2",
			data: h2b("0x000000000000000000000000000000000000000000000000000000000000007b00000000000000000000000000000000000000000000000000000000000001c800000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000e0e08c73cd5019d6cac98311fa18edf98b7042800000000000000000000000000000000000000000000000000000000000003e700000000000000000000000000000000000000000000000000000000000000020000000000000000000000000e0e08c73cd5019d6cac98311fa18edf98b704280000000000000000000000000e0e08c73cd5019d6cac98311fa18edf98b70428"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.data

			cell, err := MarshalSnakeCell(a)
			require.NoError(t, err)

			b, err := UnmarshalSnakeCell(cell)
			require.NoError(t, err)

			if len(a) != len(b) {
				t.Logf("Lengths are not equal. Want %d, got %d.\nA:\n```%s```\n\nB:\n```%s```", len(a), len(b), a, b)
				t.Logf("Last 8 bytes: A: %v; B: %v", a[len(a)-8:], b[len(b)-8:])
				t.FailNow()
			}

			require.Equal(t, a, b)
		})
	}
}

func TestDeployment(t *testing.T) {
	// ARRANGE
	// Given TSS address & Authority address
	const (
		sampleTSSPrivateKey = "0xb984cd65727cfd03081fc7bf33bf5c208bca697ce16139b5ded275887e81395a"
		sampleAuthority     = "0:4686a2c066c784a915f3e01c853d3195ed254c948e21adbb3e4a9b3f5f3c74d7"
	)

	privateKey := evmWallet(t, sampleTSSPrivateKey)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)

	tss := crypto.PubkeyToAddress(*publicKeyECDSA)

	// ACT
	code := GatewayCode()
	stateInit := GatewayStateInit(ton.MustParseAccountID(sampleAuthority), tss, true)

	// ASSERT
	_, err := code.ToBocStringCustom(false, true, false, 0)
	require.NoError(t, err)

	stateString, err := stateInit.ToBocStringCustom(false, true, false, 0)
	require.NoError(t, err)

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

func evmWallet(t *testing.T, privateKey string) *ecdsa.PrivateKey {
	pkBytes, err := hex.DecodeString(privateKey[2:])
	require.NoError(t, err)

	pk, err := crypto.ToECDSA(pkBytes)
	require.NoError(t, err)

	return pk
}

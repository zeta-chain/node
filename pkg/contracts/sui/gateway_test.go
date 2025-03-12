package sui

import (
	"encoding/base64"
	"fmt"
	"testing"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseEvent(t *testing.T) {
	// stubs
	const (
		packageID = "0x3e9fb7c01ef0d97911ccfec79306d9de2d58daa996bd3469da0f6d640cc443cf"
		gatewayID = "0x444fb7c01ef0d97911ccfec79306d9de2d58daa996bd3469da0f6d640cc443aa"
		sender    = "0x70386a9a912d9f7a603263abfbd8faae861df0ee5f8e2dbdf731fbd159f10e52"
		txHash    = "HjxLMxMXNz8YfUc2qT4e4CrogKvGeHRbDW7Arr6ntzqq"
	)

	gw := NewGateway(packageID, gatewayID)

	eventType := func(t string) string {
		return fmt.Sprintf("%s::%s::%s", packageID, gw.Module(), t)
	}

	receiverAlice := ethcommon.HexToAddress("0xa64AeD687591CfCAB52F2C1DF79a2424BbC5fEA1")
	receiverBob := ethcommon.HexToAddress("0xd4bED9bf67143d3B4A012B868E6A9566922cFDf7")

	var payload []any
	payloadBytes := []byte(base64.StdEncoding.EncodeToString([]byte{0, 1, 2}))
	for _, p := range payloadBytes {
		payload = append(payload, float64(p))
	}

	for _, tt := range []struct {
		name        string
		event       models.SuiEventResponse
		errContains string
		assert      func(t *testing.T, raw models.SuiEventResponse, out Event)
	}{
		{
			name: "deposit",
			event: models.SuiEventResponse{
				Id:        models.EventId{TxDigest: txHash, EventSeq: "0"},
				PackageId: packageID,
				Sender:    sender,
				Type:      eventType("DepositEvent"),
				ParsedJson: map[string]any{
					"coin_type": string(SUI),
					"amount":    "100",
					"sender":    sender,
					"receiver":  receiverAlice.String(),
				},
			},
			assert: func(t *testing.T, raw models.SuiEventResponse, out Event) {
				assert.Equal(t, txHash, out.TxHash)
				assert.Equal(t, DepositEvent, out.EventType)
				assert.Equal(t, uint64(0), out.EventIndex)

				deposit, err := out.Deposit()
				require.NoError(t, err)

				assert.Equal(t, SUI, deposit.CoinType)
				assert.True(t, math.NewUint(100).Equal(deposit.Amount))
				assert.Equal(t, sender, deposit.Sender)
				assert.Equal(t, receiverAlice, deposit.Receiver)
				assert.False(t, deposit.IsCrossChainCall)
				assert.True(t, deposit.IsGas())
			},
		},
		{
			name: "depositAndCall",
			event: models.SuiEventResponse{
				Id:        models.EventId{TxDigest: txHash, EventSeq: "1"},
				PackageId: packageID,
				Sender:    sender,
				Type:      eventType("DepositAndCallEvent"),
				ParsedJson: map[string]any{
					"coin_type": string(SUI),
					"amount":    "200",
					"sender":    sender,
					"receiver":  receiverBob.String(),
					"payload":   payload,
				},
			},
			assert: func(t *testing.T, raw models.SuiEventResponse, out Event) {
				assert.Equal(t, txHash, out.TxHash)
				assert.Equal(t, DepositAndCallEvent, out.EventType)
				assert.Equal(t, uint64(1), out.EventIndex)

				deposit, err := out.Deposit()
				require.NoError(t, err)

				assert.Equal(t, SUI, deposit.CoinType)
				assert.True(t, math.NewUint(200).Equal(deposit.Amount))
				assert.Equal(t, sender, deposit.Sender)
				assert.Equal(t, receiverBob, deposit.Receiver)
				assert.True(t, deposit.IsCrossChainCall)
				assert.Equal(t, []byte{0, 1, 2}, deposit.Payload)
			},
		},
		{
			name: "depositAndCall_empty_payload",
			event: models.SuiEventResponse{
				Id:        models.EventId{TxDigest: txHash, EventSeq: "1"},
				PackageId: packageID,
				Sender:    sender,
				Type:      eventType("DepositAndCallEvent"),
				ParsedJson: map[string]any{
					"coin_type": string(SUI),
					"amount":    "200",
					"sender":    sender,
					"receiver":  receiverBob.String(),
					"payload":   []any{},
				},
			},
			assert: func(t *testing.T, raw models.SuiEventResponse, out Event) {
				assert.Equal(t, txHash, out.TxHash)
				assert.Equal(t, DepositAndCallEvent, out.EventType)
				assert.Equal(t, uint64(1), out.EventIndex)

				deposit, err := out.Deposit()
				require.NoError(t, err)

				assert.Equal(t, SUI, deposit.CoinType)
				assert.True(t, math.NewUint(200).Equal(deposit.Amount))
				assert.Equal(t, sender, deposit.Sender)
				assert.Equal(t, receiverBob, deposit.Receiver)
				assert.True(t, deposit.IsCrossChainCall)
				assert.Equal(t, []byte{}, deposit.Payload)
			},
		},
		{
			name: "withdraw",
			event: models.SuiEventResponse{
				Id:        models.EventId{TxDigest: txHash, EventSeq: "1"},
				PackageId: packageID,
				Sender:    sender,
				Type:      eventType("WithdrawEvent"),
				ParsedJson: map[string]any{
					"coin_type": string(SUI),
					"amount":    "200",
					"sender":    sender,
					"receiver":  receiverBob.String(),
					"nonce":     "123",
				},
			},
			assert: func(t *testing.T, raw models.SuiEventResponse, out Event) {
				assert.Equal(t, txHash, out.TxHash)
				assert.Equal(t, WithdrawEvent, out.EventType)

				wd, err := out.Withdrawal()
				require.NoError(t, err)

				assert.Equal(t, SUI, wd.CoinType)
				assert.True(t, math.NewUint(200).Equal(wd.Amount))
				assert.Equal(t, sender, wd.Sender)
				assert.Equal(t, receiverBob.String(), wd.Receiver)
				assert.True(t, wd.IsGas())
			},
		},
		// ERRORS
		{
			name: "empty tx hash",
			event: models.SuiEventResponse{
				Id:        models.EventId{TxDigest: "", EventSeq: "0"},
				PackageId: "0x123",
			},
			errContains: "empty tx hash",
		},
		{
			name: "empty event id",
			event: models.SuiEventResponse{
				Id:        models.EventId{TxDigest: txHash, EventSeq: ""},
				PackageId: "0x123",
			},
			errContains: "empty event id",
		},
		{
			name: "invalid event id",
			event: models.SuiEventResponse{
				Id:        models.EventId{TxDigest: txHash, EventSeq: "hey"},
				PackageId: packageID,
			},
			errContains: `failed to parse event id "hey"`,
		},
		{
			name: "invalid package",
			event: models.SuiEventResponse{
				Id:        models.EventId{TxDigest: txHash, EventSeq: "0"},
				PackageId: "0x123",
			},
			errContains: "package id mismatch",
		},
		{
			name: "invalid module",
			event: models.SuiEventResponse{
				Id:                models.EventId{TxDigest: txHash, EventSeq: "0"},
				PackageId:         packageID,
				Type:              fmt.Sprintf("%s::%s::%s", packageID, "not_a_gateway", DepositEvent),
				TransactionModule: "foo",
			},
			errContains: "module mismatch",
		},
		{
			name: "invalid event type",
			event: models.SuiEventResponse{
				Id:        models.EventId{TxDigest: txHash, EventSeq: "0"},
				PackageId: packageID,
				Type:      eventType("bar"),
			},
			errContains: `unknown event "bar"`,
		},
		{
			name: "invalid coin type",
			event: models.SuiEventResponse{
				Id:        models.EventId{TxDigest: txHash, EventSeq: "0"},
				PackageId: packageID,
				Sender:    sender,
				Type:      eventType("DepositEvent"),
				ParsedJson: map[string]any{
					"coin_type": 123,
				},
			},
			errContains: "invalid coin_type",
		},
		{
			name: "invalid amount",
			event: models.SuiEventResponse{
				Id:        models.EventId{TxDigest: txHash, EventSeq: "0"},
				PackageId: packageID,
				Sender:    sender,
				Type:      eventType("DepositEvent"),
				ParsedJson: map[string]any{
					"coin_type": string(SUI),
					"amount":    "-1",
				},
			},
			errContains: "unable to parse amount",
		},
		{
			name: "invalid sender",
			event: models.SuiEventResponse{
				Id:        models.EventId{TxDigest: txHash, EventSeq: "0"},
				PackageId: packageID,
				Sender:    sender,
				Type:      eventType("DepositEvent"),
				ParsedJson: map[string]any{
					"coin_type": string(SUI),
					"amount":    "300",
					"sender":    0,
				},
			},
			errContains: "invalid sender",
		},
		{
			name: "invalid receiver can still be parsed",
			event: models.SuiEventResponse{
				Id:        models.EventId{TxDigest: txHash, EventSeq: "0"},
				PackageId: packageID,
				Sender:    sender,
				Type:      eventType("DepositEvent"),
				ParsedJson: map[string]any{
					"coin_type": string(SUI),
					"amount":    "300",
					"sender":    sender,
					"receiver":  "hello",
				},
			},
		},
		{
			name: "invalid payload",
			event: models.SuiEventResponse{
				Id:        models.EventId{TxDigest: txHash, EventSeq: "1"},
				PackageId: packageID,
				Sender:    sender,
				Type:      eventType("DepositAndCallEvent"),
				ParsedJson: map[string]any{
					"coin_type": string(SUI),
					"amount":    "200",
					"sender":    sender,
					"receiver":  receiverBob.String(),
					"payload":   []any{"boom"},
				},
			},
			errContains: "unable to convert payload: not a float64",
		},
		{
			name: "invalid payload float64",
			event: models.SuiEventResponse{
				Id:        models.EventId{TxDigest: txHash, EventSeq: "1"},
				PackageId: packageID,
				Sender:    sender,
				Type:      eventType("DepositAndCallEvent"),
				ParsedJson: map[string]any{
					"coin_type": string(SUI),
					"amount":    "200",
					"sender":    sender,
					"receiver":  receiverBob.String(),
					"payload":   []any{float64(1000)},
				},
			},
			errContains: "unable to convert payload: not a byte",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			out, err := gw.ParseEvent(tt.event)

			if tt.errContains != "" {
				require.ErrorIs(t, err, ErrParseEvent)
				require.ErrorContains(t, err, tt.errContains)
				return
			}

			require.NoError(t, err)
			tt.assert(t, tt.event, out)
		})
	}
}

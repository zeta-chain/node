package sui

import (
	"fmt"
	"testing"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/testutil/sample"
)

func TestParseEvent(t *testing.T) {
	// stubs
	const (
		packageID = "0x3e9fb7c01ef0d97911ccfec79306d9de2d58daa996bd3469da0f6d640cc443cf"
		gatewayID = "0x444fb7c01ef0d97911ccfec79306d9de2d58daa996bd3469da0f6d640cc443aa"
		sender    = "0x70386a9a912d9f7a603263abfbd8faae861df0ee5f8e2dbdf731fbd159f10e52"
		txHash    = "HjxLMxMXNz8YfUc2qT4e4CrogKvGeHRbDW7Arr6ntzqq"
	)

	eventType := func(t string) string {
		return fmt.Sprintf("%s::%s::%s", packageID, moduleName, t)
	}

	gw := NewGateway(packageID, gatewayID)

	receiverAlice := sample.EthAddress()
	receiverBob := sample.EthAddress()

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
				assert.Equal(t, Deposit, out.EventType)
				assert.Equal(t, uint64(0), out.EventIndex)

				inbound, err := out.Inbound()
				require.NoError(t, err)

				assert.Equal(t, SUI, inbound.CoinType)
				assert.True(t, math.NewUint(100).Equal(inbound.Amount))
				assert.Equal(t, sender, inbound.Sender)
				assert.Equal(t, receiverAlice, inbound.Receiver)
				assert.False(t, inbound.IsCrossChainCall)
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
					"payload":   []any{float64(0), float64(1), float64(2)},
				},
			},
			assert: func(t *testing.T, raw models.SuiEventResponse, out Event) {
				assert.Equal(t, txHash, out.TxHash)
				assert.Equal(t, DepositAndCall, out.EventType)
				assert.Equal(t, uint64(1), out.EventIndex)

				inbound, err := out.Inbound()
				require.NoError(t, err)

				assert.Equal(t, SUI, inbound.CoinType)
				assert.True(t, math.NewUint(200).Equal(inbound.Amount))
				assert.Equal(t, sender, inbound.Sender)
				assert.Equal(t, receiverBob, inbound.Receiver)
				assert.True(t, inbound.IsCrossChainCall)
				assert.Equal(t, []byte{0, 1, 2}, inbound.Payload)
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
				assert.Equal(t, DepositAndCall, out.EventType)
				assert.Equal(t, uint64(1), out.EventIndex)

				inbound, err := out.Inbound()
				require.NoError(t, err)

				assert.Equal(t, SUI, inbound.CoinType)
				assert.True(t, math.NewUint(200).Equal(inbound.Amount))
				assert.Equal(t, sender, inbound.Sender)
				assert.Equal(t, receiverBob, inbound.Receiver)
				assert.True(t, inbound.IsCrossChainCall)
				assert.Equal(t, []byte{}, inbound.Payload)
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
				Type:              fmt.Sprintf("%s::%s::%s", packageID, "not_a_gateway", Deposit),
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
			name: "invalid receiver",
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
			errContains: `invalid receiver address "hello"`,
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

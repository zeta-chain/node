package sui

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"slices"
	"testing"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ActiveMessageContextDynamicFieldName(t *testing.T) {
	got, err := ActiveMessageContextDynamicFieldName()
	require.NoError(t, err)

	expectedJSON := json.RawMessage(`[97,99,116,105,118,101,95,109,101,115,115,97,103,101,95,99,111,110,116,101,120,116]`)
	require.Equal(t, expectedJSON, got)
}

func TestNewGatewayFromAddress(t *testing.T) {
	// stubs
	const (
		packageID         = "0x3e9fb7c01ef0d97911ccfec79306d9de2d58daa996bd3469da0f6d640cc443cf"
		gatewayID         = "0x444fb7c01ef0d97911ccfec79306d9de2d58daa996bd3469da0f6d640cc443aa"
		withdrawCapID     = "0x84d96419097f3cd66c7dd732cd28c8df58c1183768ae617c0705a6261a60a870"
		originalPackageID = "0x9a6e7366064fb27ac1daeca6f7d4c13af2f86d26433b5e70bea9b6214e6253e4"
	)

	tests := []struct {
		name         string
		address      string
		wantErr      string
		wantOriginal bool
	}{
		{
			name:    "valid legacy gateway address",
			address: MakeAddress(packageID, gatewayID, "", ""),
		},
		{
			name:         "valid new gateway address with original package id and withdraw cap id",
			address:      MakeAddress(packageID, gatewayID, withdrawCapID, originalPackageID),
			wantOriginal: true,
		},
		{
			name:    "invalid gateway address, empty string",
			address: "",
			wantErr: "invalid gateway address",
		},
		{
			name:    "invalid gateway address, contains 1 part",
			address: "0x123",
			wantErr: "invalid gateway address",
		},
		{
			name:    "invalid gateway address, contains 3 parts",
			address: fmt.Sprintf("%s,%s,%s", packageID, gatewayID, originalPackageID),
			wantErr: "invalid gateway address",
		},
		{
			name:    "invalid Sui address",
			address: fmt.Sprintf("%s,0xabc", packageID),
			wantErr: "invalid Sui address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gw, err := NewGatewayFromAddress(tt.address)
			if tt.wantErr != "" {
				require.Nil(t, gw)
				require.ErrorContains(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, packageID, gw.PackageID())
			assert.Equal(t, gatewayID, gw.ObjectID())

			if !tt.wantOriginal {
				assert.Equal(t, []string{packageID}, gw.PackageIDs())
				return
			}

			assert.Equal(t, withdrawCapID, gw.WithdrawCapID())
			assert.Equal(t, originalPackageID, gw.Original().PackageID())
			assert.True(t, slices.Equal([]string{packageID, originalPackageID}, gw.PackageIDs()))
		})
	}
}

func Test_MakeAddress(t *testing.T) {
	const (
		packageID         = "0x3e9fb7c01ef0d97911ccfec79306d9de2d58daa996bd3469da0f6d640cc443cf"
		gatewayID         = "0x444fb7c01ef0d97911ccfec79306d9de2d58daa996bd3469da0f6d640cc443aa"
		withdrawCapID     = "0x84d96419097f3cd66c7dd732cd28c8df58c1183768ae617c0705a6261a60a870"
		originalPackageID = "0x9a6e7366064fb27ac1daeca6f7d4c13af2f86d26433b5e70bea9b6214e6253e4"
	)

	t.Run("original package id is empty", func(t *testing.T) {
		gatewayAddress := MakeAddress(packageID, gatewayID, "", "")
		assert.Equal(t, fmt.Sprintf("%s,%s", packageID, gatewayID), gatewayAddress)
	})

	t.Run("original package id is not empty", func(t *testing.T) {
		gatewayAddress := MakeAddress(packageID, gatewayID, withdrawCapID, originalPackageID)
		assert.Equal(t, fmt.Sprintf("%s,%s,%s,%s", packageID, gatewayID, withdrawCapID, originalPackageID), gatewayAddress)
	})
}

func Test_ToAddress(t *testing.T) {
	const (
		packageID         = "0x3e9fb7c01ef0d97911ccfec79306d9de2d58daa996bd3469da0f6d640cc443cf"
		gatewayID         = "0x444fb7c01ef0d97911ccfec79306d9de2d58daa996bd3469da0f6d640cc443aa"
		withdrawCapID     = "0x84d96419097f3cd66c7dd732cd28c8df58c1183768ae617c0705a6261a60a870"
		originalPackageID = "0x9a6e7366064fb27ac1daeca6f7d4c13af2f86d26433b5e70bea9b6214e6253e4"
	)

	t.Run("original package id is empty", func(t *testing.T) {
		gw := NewGateway(packageID, gatewayID)
		assert.Equal(t, MakeAddress(packageID, gatewayID, withdrawCapID, ""), gw.ToAddress())
	})

	t.Run("original package id is not empty", func(t *testing.T) {
		gatewayAddress := MakeAddress(packageID, gatewayID, withdrawCapID, originalPackageID)
		gw, err := NewGatewayFromAddress(gatewayAddress)
		require.NoError(t, err)
		assert.Equal(t, gatewayAddress, gw.ToAddress())
	})
}

func Test_Original(t *testing.T) {
	const (
		packageID         = "0x3e9fb7c01ef0d97911ccfec79306d9de2d58daa996bd3469da0f6d640cc443cf"
		gatewayID         = "0x444fb7c01ef0d97911ccfec79306d9de2d58daa996bd3469da0f6d640cc443aa"
		withdrawCapID     = "0x84d96419097f3cd66c7dd732cd28c8df58c1183768ae617c0705a6261a60a870"
		originalPackageID = "0x9a6e7366064fb27ac1daeca6f7d4c13af2f86d26433b5e70bea9b6214e6253e4"
	)

	t.Run("original package id is not empty", func(t *testing.T) {
		gw, err := NewGatewayFromAddress(MakeAddress(packageID, gatewayID, withdrawCapID, originalPackageID))
		require.NoError(t, err)

		gwOriginal := gw.Original()
		assert.Equal(t, originalPackageID, gwOriginal.PackageID())
		assert.Equal(t, gatewayID, gwOriginal.ObjectID())
		assert.Equal(t, withdrawCapID, gwOriginal.WithdrawCapID())
		assert.Empty(t, gwOriginal.originalPackageID)
	})

	t.Run("original package id is empty", func(t *testing.T) {
		gw, err := NewGatewayFromAddress(MakeAddress(packageID, gatewayID, withdrawCapID, ""))
		require.NoError(t, err)

		gwOriginal := gw.Original()
		assert.Equal(t, packageID, gwOriginal.PackageID())
		assert.Equal(t, gatewayID, gwOriginal.ObjectID())
		assert.Empty(t, gwOriginal.WithdrawCapID())
		assert.Empty(t, gwOriginal.originalPackageID)
	})
}

func Test_UpdateIDs(t *testing.T) {
	const (
		packageID = "0x3e9fb7c01ef0d97911ccfec79306d9de2d58daa996bd3469da0f6d640cc443cf"
		gatewayID = "0x444fb7c01ef0d97911ccfec79306d9de2d58daa996bd3469da0f6d640cc443aa"

		packageID2        = "0x9a6e7366064fb27ac1daeca6f7d4c13af2f86d26433b5e70bea9b6214e6253e4"
		gatewayID2        = "0xaf52affd195806d9aa9d967462cbda411bfed9a6efc4a032bf8e34a391469878"
		withdrawCapID     = "0x84d96419097f3cd66c7dd732cd28c8df58c1183768ae617c0705a6261a60a870"
		originalPackageID = "0x9a6e7366064fb27ac1daeca6f7d4c13af2f86d26433b5e70bea9b6214e6253e4"
	)

	// before update
	gw := NewGateway(packageID, gatewayID)
	assert.Equal(t, packageID, gw.PackageID())
	assert.Equal(t, gatewayID, gw.ObjectID())
	assert.Empty(t, gw.WithdrawCapID())
	assert.Empty(t, gw.originalPackageID)

	// after update
	require.NoError(t, gw.UpdateIDs(MakeAddress(packageID2, gatewayID2, withdrawCapID, originalPackageID)))
	assert.Equal(t, packageID2, gw.PackageID())
	assert.Equal(t, gatewayID2, gw.ObjectID())
	assert.Equal(t, withdrawCapID, gw.WithdrawCapID())
	assert.Equal(t, originalPackageID, gw.Original().PackageID())
}

func TestParseEvent(t *testing.T) {
	// stubs
	const (
		packageID         = "0x3e9fb7c01ef0d97911ccfec79306d9de2d58daa996bd3469da0f6d640cc443cf"
		gatewayID         = "0x444fb7c01ef0d97911ccfec79306d9de2d58daa996bd3469da0f6d640cc443aa"
		withdrawCapID     = "0x84d96419097f3cd66c7dd732cd28c8df58c1183768ae617c0705a6261a60a870"
		originalPackageID = "0x9a6e7366064fb27ac1daeca6f7d4c13af2f86d26433b5e70bea9b6214e6253e4"
		sender            = "0x70386a9a912d9f7a603263abfbd8faae861df0ee5f8e2dbdf731fbd159f10e52"
		txHash            = "HjxLMxMXNz8YfUc2qT4e4CrogKvGeHRbDW7Arr6ntzqq"
	)

	gw, err := NewGatewayFromAddress(MakeAddress(packageID, gatewayID, withdrawCapID, originalPackageID))
	require.NoError(t, err)

	eventType := func(t string) string {
		return fmt.Sprintf("%s::%s::%s", originalPackageID, GatewayModule, t)
	}

	receiverAlice := ethcommon.HexToAddress("0xa64AeD687591CfCAB52F2C1DF79a2424BbC5fEA1")
	receiverBob := ethcommon.HexToAddress("0xd4bED9bf67143d3B4A012B868E6A9566922cFDf7")

	payload := []any{float64(0), float64(1), float64(2)}

	var payloadBase64 []any
	payloadBytes := []byte(base64.StdEncoding.EncodeToString([]byte{3, 4, 5}))
	for _, p := range payloadBytes {
		payloadBase64 = append(payloadBase64, float64(p))
	}

	for _, tt := range []struct {
		name        string
		event       models.SuiEventResponse
		errContains string
		assert      func(t *testing.T, raw models.SuiEventResponse, out Event)
	}{
		{
			name: "deposit from non-original gateway",
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
			name: "depositAndCall with bytes payload",
			event: models.SuiEventResponse{
				Id:        models.EventId{TxDigest: txHash, EventSeq: "1"},
				PackageId: originalPackageID,
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
			name: "depositAndCall with Base64 formatted payload",
			event: models.SuiEventResponse{
				Id:        models.EventId{TxDigest: txHash, EventSeq: "1"},
				PackageId: originalPackageID,
				Sender:    sender,
				Type:      eventType("DepositAndCallEvent"),
				ParsedJson: map[string]any{
					"coin_type": string(SUI),
					"amount":    "200",
					"sender":    sender,
					"receiver":  receiverBob.String(),
					"payload":   payloadBase64,
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
				assert.Equal(t, []byte{3, 4, 5}, deposit.Payload)
			},
		},
		{
			name: "depositAndCall_empty_payload",
			event: models.SuiEventResponse{
				Id:        models.EventId{TxDigest: txHash, EventSeq: "1"},
				PackageId: originalPackageID,
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
				PackageId: originalPackageID,
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
				PackageId: originalPackageID,
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
				PackageId:         originalPackageID,
				Type:              fmt.Sprintf("%s::%s::%s", originalPackageID, "not_a_gateway", DepositEvent),
				TransactionModule: "foo",
			},
			errContains: "module mismatch",
		},
		{
			name: "invalid event type",
			event: models.SuiEventResponse{
				Id:        models.EventId{TxDigest: txHash, EventSeq: "0"},
				PackageId: originalPackageID,
				Type:      eventType("bar"),
			},
			errContains: `unknown event "bar"`,
		},
		{
			name: "invalid coin type",
			event: models.SuiEventResponse{
				Id:        models.EventId{TxDigest: txHash, EventSeq: "0"},
				PackageId: originalPackageID,
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
				PackageId: originalPackageID,
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
				PackageId: originalPackageID,
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
				PackageId: originalPackageID,
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
				PackageId: originalPackageID,
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
				PackageId: originalPackageID,
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

func Test_ParseOutboundEvent(t *testing.T) {
	// stubs
	const (
		packageID = "0x3e9fb7c01ef0d97911ccfec79306d9de2d58daa996bd3469da0f6d640cc443cf"
		gatewayID = "0x444fb7c01ef0d97911ccfec79306d9de2d58daa996bd3469da0f6d640cc443aa"
		sender    = "0x70386a9a912d9f7a603263abfbd8faae861df0ee5f8e2dbdf731fbd159f10e52"
		txHash    = "HjxLMxMXNz8YfUc2qT4e4CrogKvGeHRbDW7Arr6ntzqq"
		receiver  = "0xd4bED9bf67143d3B4A012B868E6A9566922cFDf7"
	)

	gw := NewGateway(packageID, gatewayID)

	eventType := func(t string) string {
		return fmt.Sprintf("%s::%s::%s", packageID, GatewayModule, t)
	}

	for _, tt := range []struct {
		name      string
		response  models.SuiTransactionBlockResponse
		wantEvent Event
		errMsg    string
	}{
		{
			name: "withdraw",
			response: models.SuiTransactionBlockResponse{
				Events: []models.SuiEventResponse{
					{
						Id:        models.EventId{TxDigest: txHash, EventSeq: "1"},
						PackageId: packageID,
						Sender:    sender,
						Type:      eventType("WithdrawEvent"),
						ParsedJson: map[string]any{
							"coin_type": string(SUI),
							"amount":    "200",
							"sender":    sender,
							"receiver":  receiver,
							"nonce":     "123",
						},
					},
				},
			},
			wantEvent: Event{
				TxHash:     txHash,
				EventIndex: 1,
				EventType:  WithdrawEvent,
				content: Withdrawal{
					CoinType: SUI,
					Amount:   math.NewUint(200),
					Sender:   sender,
					Receiver: receiver,
					Nonce:    123,
				},
			},
		},
		{
			name:     "withdrawAndCall with PTB",
			response: createPTBResponse(txHash, packageID, "200", "123"),
			wantEvent: Event{
				TxHash:     txHash,
				EventIndex: 0,
				EventType:  WithdrawAndCallEvent,
				content: WithdrawAndCallPTB{
					MoveCall: MoveCall{
						PackageID:  packageID,
						Module:     GatewayModule,
						Function:   FuncWithdrawImpl,
						ArgIndexes: ptbWithdrawImplArgIndexes,
					},
					Amount: math.NewUint(200),
					Nonce:  123,
				},
			},
		},
		{
			name: "cancelTx",
			response: models.SuiTransactionBlockResponse{
				Events: []models.SuiEventResponse{
					{
						Id:        models.EventId{TxDigest: txHash, EventSeq: "1"},
						PackageId: packageID,
						Sender:    sender,
						Type:      eventType("NonceIncreaseEvent"),
						ParsedJson: map[string]any{
							"nonce":  "123",
							"sender": sender,
						},
					},
				},
			},
			wantEvent: Event{
				TxHash:     txHash,
				EventIndex: 1,
				EventType:  CancelTxEvent,
				content: CancelTx{
					Nonce:  123,
					Sender: sender,
				},
			},
		},
		{
			name: "no event",
			response: models.SuiTransactionBlockResponse{
				Events: []models.SuiEventResponse{},
			},
			errMsg: "missing events",
		},
		{
			name: "unable to parse event",
			response: models.SuiTransactionBlockResponse{
				Events: []models.SuiEventResponse{
					{
						Id: models.EventId{TxDigest: "", EventSeq: ""}, // invalid EventId
					},
				},
			},
			errMsg: "unable to parse event",
		},
		{
			name: "not an outbound event",
			response: models.SuiTransactionBlockResponse{
				Events: []models.SuiEventResponse{
					{
						Id:        models.EventId{TxDigest: txHash, EventSeq: "0"},
						PackageId: packageID,
						Sender:    sender,
						Type:      eventType("DepositEvent"),
						ParsedJson: map[string]any{
							"coin_type": string(SUI),
							"amount":    "100",
							"sender":    sender,
							"receiver":  receiver,
						},
					},
				},
			},
			errMsg: "unsupported outbound event type",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			event, content, err := gw.ParseOutboundEvent(tt.response)

			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantEvent, event)
			require.Equal(t, tt.wantEvent.content, content)
		})
	}
}

func Test_ParseDynamicFieldValueStr(t *testing.T) {
	tests := []struct {
		name   string
		data   models.SuiParsedData
		want   string
		errMsg string
	}{
		{
			name: "valid dynamic field value",
			data: models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: map[string]any{
						"value": "0x123",
					},
				},
			},
			want: "0x123",
		},
		{
			name: "missing value field",
			data: models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: map[string]any{},
				},
			},
			errMsg: "missing value field",
		},
		{
			name: "value field type mismatch",
			data: models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: map[string]any{"value": 123},
				},
			},
			errMsg: "want string, got int for dynamic field value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDynamicFieldValueStr(tt.data)
			if tt.errMsg != "" {
				require.Empty(t, got)
				require.ErrorContains(t, err, tt.errMsg)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_ParseGatewayNonce(t *testing.T) {
	tests := []struct {
		name   string
		data   models.SuiParsedData
		nonce  uint64
		errMsg string
	}{
		{
			name: "valid nonce",
			data: models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: map[string]any{
						"nonce": "123",
					},
				},
			},
			nonce: 123,
		},
		{
			name: "missing nonce field",
			data: models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: map[string]any{},
				},
			},
			errMsg: "missing nonce field",
		},
		{
			name: "invalid nonce field",
			data: models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: map[string]any{"nonce": 123},
				},
			},
			errMsg: "want string, got int for nonce",
		},
		{
			name: "invalid nonce value",
			data: models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: map[string]any{"nonce": "not a number"},
				},
			},
			errMsg: "unable to parse nonce",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nonce, err := ParseGatewayNonce(tt.data)
			if tt.errMsg != "" {
				require.Zero(t, nonce)
				require.ErrorContains(t, err, tt.errMsg)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.nonce, nonce)
		})
	}
}

func Test_extractInteger(t *testing.T) {
	tests := []struct {
		name   string
		kv     map[string]any
		key    string
		want   any
		errMsg string
	}{
		{
			name: "valid int8",
			kv:   map[string]any{"key": float64(42)},
			key:  "key",
			want: int8(42),
		},
		{
			name: "valid int16",
			kv:   map[string]any{"key": float64(1000)},
			key:  "key",
			want: int16(1000),
		},
		{
			name: "valid int32",
			kv:   map[string]any{"key": float64(100000)},
			key:  "key",
			want: int32(100000),
		},
		{
			name: "valid int64",
			kv:   map[string]any{"key": float64(1000000000)},
			key:  "key",
			want: int64(1000000000),
		},
		{
			name: "valid uint8",
			kv:   map[string]any{"key": float64(42)},
			key:  "key",
			want: uint8(42),
		},
		{
			name: "valid uint16",
			kv:   map[string]any{"key": float64(1000)},
			key:  "key",
			want: uint16(1000),
		},
		{
			name: "valid uint32",
			kv:   map[string]any{"key": float64(100000)},
			key:  "key",
			want: uint32(100000),
		},
		{
			name: "valid uint64",
			kv:   map[string]any{"key": float64(1000000000)},
			key:  "key",
			want: uint64(1000000000),
		},
		{
			name:   "missing key",
			kv:     map[string]any{},
			key:    "key",
			errMsg: "missing key",
		},
		{
			name:   "invalid value type",
			kv:     map[string]any{"key": "not a number"},
			key:    "key",
			errMsg: "want float64, got string for key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "valid int8":
				got, err := extractInteger[int8](tt.kv, tt.key)
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			case "valid int16":
				got, err := extractInteger[int16](tt.kv, tt.key)
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			case "valid int32":
				got, err := extractInteger[int32](tt.kv, tt.key)
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			case "valid int64":
				got, err := extractInteger[int64](tt.kv, tt.key)
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			case "valid uint8":
				got, err := extractInteger[uint8](tt.kv, tt.key)
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			case "valid uint16":
				got, err := extractInteger[uint16](tt.kv, tt.key)
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			case "valid uint32":
				got, err := extractInteger[uint32](tt.kv, tt.key)
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			case "valid uint64":
				got, err := extractInteger[uint64](tt.kv, tt.key)
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			default:
				// Test error cases for all types
				if tt.errMsg != "" {
					// Test with int64 as an example
					_, err := extractInteger[int64](tt.kv, tt.key)
					require.ErrorContains(t, err, tt.errMsg)
				}
			}
		})
	}
}

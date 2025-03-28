package sui

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/stretchr/testify/require"
)

func Test_parseWithdrawal(t *testing.T) {
	const (
		sender   = "0x70386a9a912d9f7a603263abfbd8faae861df0ee5f8e2dbdf731fbd159f10e52"
		receiver = "0xd4bED9bf67143d3B4A012B868E6A9566922cFDf7"
	)

	sampleEventResponse := func() models.SuiEventResponse {
		return models.SuiEventResponse{
			ParsedJson: map[string]any{
				"coin_type": string(SUI),
				"amount":    "100",
				"sender":    sender,
				"receiver":  receiver,
				"nonce":     "1",
			},
		}
	}

	tests := []struct {
		name      string
		eventType EventType
		event     models.SuiEventResponse
		want      Withdrawal
		wantErr   bool
	}{
		{
			name:      "should be able to parse valid event",
			eventType: WithdrawEvent,
			event:     sampleEventResponse(),
			want: Withdrawal{
				CoinType: SUI,
				Amount:   math.NewUint(100),
				Sender:   sender,
				Receiver: receiver,
				Nonce:    1,
			},
		},
		{
			name:      "should return error on invalid event type",
			eventType: "invalid event type",
			event:     sampleEventResponse(),
			wantErr:   true,
		},
		{
			name:      "should return error if coin type is missing",
			eventType: WithdrawEvent,
			event: func() models.SuiEventResponse {
				resp := sampleEventResponse()
				delete(resp.ParsedJson, "coin_type")
				return resp
			}(),
			wantErr: true,
		},
		{
			name:      "should return error if amount is missing",
			eventType: WithdrawEvent,
			event: func() models.SuiEventResponse {
				resp := sampleEventResponse()
				delete(resp.ParsedJson, "amount")
				return resp
			}(),
			wantErr: true,
		},
		{
			name:      "should return error on invalid amount",
			eventType: WithdrawEvent,
			event: func() models.SuiEventResponse {
				resp := sampleEventResponse()
				resp.ParsedJson["amount"] = "not a number"
				return resp
			}(),
			wantErr: true,
		},
		{
			name:      "should return error if sender is missing",
			eventType: WithdrawEvent,
			event: func() models.SuiEventResponse {
				resp := sampleEventResponse()
				delete(resp.ParsedJson, "sender")
				return resp
			}(),
			wantErr: true,
		},
		{
			name:      "should return error if receiver is missing",
			eventType: WithdrawEvent,
			event: func() models.SuiEventResponse {
				resp := sampleEventResponse()
				delete(resp.ParsedJson, "receiver")
				return resp
			}(),
			wantErr: true,
		},
		{
			name:      "should return error if nonce is missing",
			eventType: WithdrawEvent,
			event: func() models.SuiEventResponse {
				resp := sampleEventResponse()
				delete(resp.ParsedJson, "nonce")
				return resp
			}(),
			wantErr: true,
		},
		{
			name:      "should return error on invalid nonce",
			eventType: WithdrawEvent,
			event: func() models.SuiEventResponse {
				resp := sampleEventResponse()
				resp.ParsedJson["nonce"] = "not a number"
				return resp
			}(),
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseWithdrawal(test.event, test.eventType)
			if test.wantErr {
				require.Error(t, err)
				require.Empty(t, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}

func Test_Withdrawal_TokenAmount(t *testing.T) {
	event := Withdrawal{
		Amount: math.NewUint(100),
		Nonce:  1,
	}
	require.Equal(t, math.NewUint(100), event.TokenAmount())
}

func Test_Withdrawal_TxNonce(t *testing.T) {
	event := Withdrawal{
		Amount: math.NewUint(100),
		Nonce:  1,
	}
	require.Equal(t, uint64(1), event.TxNonce())
}

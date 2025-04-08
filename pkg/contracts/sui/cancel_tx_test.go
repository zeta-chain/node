package sui

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/stretchr/testify/require"
)

func Test_parseCancelTx(t *testing.T) {
	tests := []struct {
		name      string
		eventType EventType
		event     models.SuiEventResponse
		want      CancelTx
		errMsg    string
	}{
		{
			name:      "should be able to parse valid event",
			eventType: CancelTxEvent,
			event: models.SuiEventResponse{
				ParsedJson: map[string]any{
					"nonce":  "1",
					"sender": "0x8e5016551584818c2fbd10ba63a359e816f31d576ac1ec06a8b9efd1c4768a26",
				},
			},
			want: CancelTx{
				Nonce:  1,
				Sender: "0x8e5016551584818c2fbd10ba63a359e816f31d576ac1ec06a8b9efd1c4768a26",
			},
		},
		{
			name:      "should return error on invalid event type",
			eventType: "invalid event type",
			event: models.SuiEventResponse{
				ParsedJson: map[string]any{
					"nonce":  "1",
					"sender": "0x8e5016551584818c2fbd10ba63a359e816f31d576ac1ec06a8b9efd1c4768a26",
				},
			},
			errMsg: "invalid event type",
		},
		{
			name:      "should return error if nonce is missing",
			eventType: CancelTxEvent,
			event: models.SuiEventResponse{
				ParsedJson: map[string]any{
					"sender": "0x8e5016551584818c2fbd10ba63a359e816f31d576ac1ec06a8b9efd1c4768a26",
				},
			},
			errMsg: "unable to extract nonce",
		},
		{
			name:      "should return error if sender is missing",
			eventType: CancelTxEvent,
			event: models.SuiEventResponse{
				ParsedJson: map[string]any{
					"nonce": "1",
				},
			},
			errMsg: "unable to extract sender",
		},
		{
			name:      "should return error if nonce is invalid",
			eventType: CancelTxEvent,
			event: models.SuiEventResponse{
				ParsedJson: map[string]any{
					"nonce":  "not a number",
					"sender": "0x8e5016551584818c2fbd10ba63a359e816f31d576ac1ec06a8b9efd1c4768a26",
				},
			},
			errMsg: "unable to parse nonce",
		},
		{
			name:      "should return error if nonce is not positive",
			eventType: CancelTxEvent,
			event: models.SuiEventResponse{
				ParsedJson: map[string]any{
					"nonce":  "0",
					"sender": "0x8e5016551584818c2fbd10ba63a359e816f31d576ac1ec06a8b9efd1c4768a26",
				},
			},
			errMsg: "nonce must be positive",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseCancelTx(test.event, test.eventType)
			if test.errMsg != "" {
				require.ErrorContains(t, err, test.errMsg)
				require.Empty(t, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}

func Test_TokenAmount(t *testing.T) {
	event := CancelTx{
		Nonce:  1,
		Sender: "0x8e5016551584818c2fbd10ba63a359e816f31d576ac1ec06a8b9efd1c4768a26",
	}
	require.Equal(t, math.NewUint(0), event.TokenAmount())
}

func Test_TxNonce(t *testing.T) {
	event := CancelTx{
		Nonce:  1,
		Sender: "0x8e5016551584818c2fbd10ba63a359e816f31d576ac1ec06a8b9efd1c4768a26",
	}
	require.Equal(t, uint64(0), event.TxNonce())
}

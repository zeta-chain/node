package sui_test

import (
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/contracts/sui"
	"github.com/zeta-chain/node/testutil/sample"
	"testing"
)

func TestInbound_IsGasDeposit(t *testing.T) {
	tests := []struct {
		name string
		d    *sui.Inbound
		want bool
	}{
		{
			name: "gas deposit",
			d: &sui.Inbound{
				TxHash:     "0x123",
				EventIndex: 1,
				CoinType:   sui.SUI,
				Amount:     100,
				Sender:     "0x456",
				Receiver:   sample.EthAddress(),
				Payload:    nil,
			},
			want: true,
		},
		{
			name: "not gas deposit",
			d: &sui.Inbound{
				TxHash:     "0x123",
				EventIndex: 1,
				CoinType:   "not_sui",
				Amount:     100,
				Sender:     "0x456",
				Receiver:   sample.EthAddress(),
				Payload:    nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want {
				require.True(t, tt.d.IsGasDeposit())
			} else {
				require.False(t, tt.d.IsGasDeposit())
			}
		})
	}
}

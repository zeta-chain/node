package types_test

import (
	"testing"

	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestOutboundTracker_IsMaxed(t *testing.T) {
	tests := []struct {
		name    string
		tracker types.OutboundTracker
		want    bool
	}{
		{"Not maxed", types.OutboundTracker{HashList: []*types.TxHash{
			{TxHash: "hash1", TxSigner: "signer1"},
			{TxHash: "hash2", TxSigner: "signer2"},
			{TxHash: "hash3", TxSigner: "signer3"},
		}},
			false},

		{"Maxed", types.OutboundTracker{HashList: []*types.TxHash{
			{TxHash: "hash1", TxSigner: "signer1"},
			{TxHash: "hash2", TxSigner: "signer2"},
			{TxHash: "hash3", TxSigner: "signer3"},
			{TxHash: "hash4", TxSigner: "signer4"},
			{TxHash: "hash5", TxSigner: "signer5"},
		}},
			true},
		{"More than Maxed", types.OutboundTracker{HashList: []*types.TxHash{
			{TxHash: "hash1", TxSigner: "signer1"},
			{TxHash: "hash2", TxSigner: "signer2"},
			{TxHash: "hash3", TxSigner: "signer3"},
			{TxHash: "hash4", TxSigner: "signer4"},
			{TxHash: "hash5", TxSigner: "signer5"},
			{TxHash: "hash6", TxSigner: "signer6"},
			{TxHash: "hash7", TxSigner: "signer7"},
		}},
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tracker.MaxReached(); got != tt.want {
				t.Errorf("OutboundTracker.MaxReached() = %v, want %v", got, tt.want)
			}
		})
	}
}

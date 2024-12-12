package coin_test

import (
	"testing"

	"github.com/zeta-chain/node/pkg/coin"
)

func TestCoinType_SupportsRefund(t *testing.T) {
	tests := []struct {
		name string
		c    coin.CoinType
		want bool
	}{
		{"ERC20", coin.CoinType_ERC20, true},
		{"Gas", coin.CoinType_Gas, true},
		{"Zeta", coin.CoinType_Zeta, true},
		{"Cmd", coin.CoinType_Cmd, false},
		{"Unknown", coin.CoinType(100), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.SupportsRefund(); got != tt.want {
				t.Errorf("CoinType.SupportsRefund() = %v, want %v", got, tt.want)
			}
		})
	}
}

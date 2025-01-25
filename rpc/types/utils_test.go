package types_test

import (
	"math/big"
	"testing"

	"github.com/zeta-chain/node/rpc/types"
)

func TestCheckTxFee(t *testing.T) {
	tests := []struct {
		name     string
		gasPrice *big.Int
		gas      uint64
		feeCap   float64
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid transaction under cap",
			gasPrice: big.NewInt(100000000000), // 100 Gwei
			gas:      21000,                    // Standard ETH transfer
			feeCap:   1.0,                      // 1 ETH cap
			wantErr:  false,
			errMsg:   "",
		},
		{
			name:     "transaction exceeds cap",
			gasPrice: big.NewInt(500000000000), // 500 Gwei
			gas:      1000000,                  // Complex contract interaction
			feeCap:   0.1,                      // 0.1 ETH cap
			wantErr:  true,
			errMsg:   "tx fee (0.50 ether) exceeds the configured cap (0.10 ether)",
		},
		{
			name:     "nil gas price",
			gasPrice: nil,
			gas:      21000,
			feeCap:   1.0,
			wantErr:  true,
			errMsg:   "gasprice is nil",
		},
		{
			name:     "zero fee cap",
			gasPrice: big.NewInt(100000000000),
			gas:      21000,
			feeCap:   0,
			wantErr:  false,
			errMsg:   "",
		},
		{
			name:     "very low gas price",
			gasPrice: big.NewInt(1),
			gas:      21000,
			feeCap:   1.0,
			wantErr:  false,
			errMsg:   "",
		},
		{
			name:     "very high gas price",
			gasPrice: new(big.Int).Mul(big.NewInt(1e18), big.NewInt(100)), // 100 ETH per gas unit
			gas:      21000,
			feeCap:   1000.0,
			wantErr:  true,
			errMsg:   "tx fee (2100000.00 ether) exceeds the configured cap (1000.00 ether)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := types.CheckTxFee(tt.gasPrice, tt.gas, tt.feeCap)

			// Check if we expected an error
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckTxFee() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expected an error, verify the error message
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("CheckTxFee() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

package errors_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/coin"
	pkgerrors "github.com/zeta-chain/node/pkg/errors"
	"github.com/zeta-chain/node/testutil/sample"
)

func Test_ErrTxMonitor(t *testing.T) {
	err := errors.New("test error")
	msg := sample.InboundVote(coin.CoinType_Gas, 1, 2)

	tests := []struct {
		name         string
		ErrTxMonitor pkgerrors.ErrTxMonitor
		zetaTxHash   string
		expectError  string
	}{
		{
			name: "nil error returns monitoring completed message",
			ErrTxMonitor: pkgerrors.ErrTxMonitor{
				Err:        nil,
				ZetaTxHash: "test-hash-1",
				Msg:        msg,
			},
			expectError: "monitoring completed without error",
		},
		{
			name: "actual error returns formatted error message",
			ErrTxMonitor: pkgerrors.ErrTxMonitor{
				Err:        err,
				ZetaTxHash: "test-hash-2",
				Msg:        msg,
			},
			expectError: fmt.Sprintf("monitoring error: %v, inbound hash: %s, zeta tx hash: %s, ballot index: %s", err, msg.InboundHash, "test-hash-2", msg.Digest()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errMsg := tt.ErrTxMonitor.Error()
			require.Equal(t, errMsg, tt.expectError)
		})
	}
}

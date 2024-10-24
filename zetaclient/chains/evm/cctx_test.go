package evm_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/evm"
)

func TestParseOutboundTypeFromCCTX(t *testing.T) {
	tests := []struct {
		name     string
		cctx     types.CrossChainTx
		expected evm.OutboundType
	}{
		{
			name: "Gas withdraw and call",
			cctx: types.CrossChainTx{
				InboundParams: &types.InboundParams{
					CoinType:         coin.CoinType_Gas,
					IsCrossChainCall: true,
				},
				CctxStatus: &types.Status{
					Status: types.CctxStatus_PendingOutbound,
				},
			},
			expected: evm.OutboundTypeGasWithdrawAndCall,
		},
		{
			name: "Gas withdraw",
			cctx: types.CrossChainTx{
				InboundParams: &types.InboundParams{
					CoinType:         coin.CoinType_Gas,
					IsCrossChainCall: false,
				},
				CctxStatus: &types.Status{
					Status: types.CctxStatus_PendingOutbound,
				},
			},
			expected: evm.OutboundTypeGasWithdraw,
		},
		{
			name: "ERC20 withdraw and call",
			cctx: types.CrossChainTx{
				InboundParams: &types.InboundParams{
					CoinType:         coin.CoinType_ERC20,
					IsCrossChainCall: true,
				},
				CctxStatus: &types.Status{
					Status: types.CctxStatus_PendingOutbound,
				},
			},
			expected: evm.OutboundTypeERC20WithdrawAndCall,
		},
		{
			name: "ERC20 withdraw revert and call on revert",
			cctx: types.CrossChainTx{
				InboundParams: &types.InboundParams{
					CoinType: coin.CoinType_ERC20,
				},
				CctxStatus: &types.Status{
					Status: types.CctxStatus_PendingRevert,
				},
				RevertOptions: types.RevertOptions{
					CallOnRevert: true,
				},
			},
			expected: evm.OutboundTypeERC20WithdrawRevertAndCallOnRevert,
		},
		{
			name: "No asset call",
			cctx: types.CrossChainTx{
				InboundParams: &types.InboundParams{
					CoinType: coin.CoinType_NoAssetCall,
				},
				CctxStatus: &types.Status{
					Status: types.CctxStatus_PendingOutbound,
				},
			},
			expected: evm.OutboundTypeCall,
		},
		{
			name: "ZETA gives Uuknown outbound type",
			cctx: types.CrossChainTx{
				InboundParams: &types.InboundParams{
					CoinType: coin.CoinType_Zeta,
				},
				CctxStatus: &types.Status{
					Status: types.CctxStatus_PendingOutbound,
				},
			},
			expected: evm.OutboundTypeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evm.ParseOutboundTypeFromCCTX(tt.cctx)
			require.Equal(t, tt.expected, result)
		})
	}
}

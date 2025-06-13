package common

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestParseOutboundTypeFromCCTX(t *testing.T) {
	tests := []struct {
		name     string
		cctx     types.CrossChainTx
		expected OutboundType
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
			expected: OutboundTypeGasWithdrawAndCall,
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
			expected: OutboundTypeGasWithdraw,
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
			expected: OutboundTypeERC20WithdrawAndCall,
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
			expected: OutboundTypeERC20WithdrawRevertAndCallOnRevert,
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
			expected: OutboundTypeCall,
		},
		{
			name: "ZETA withdraw",
			cctx: types.CrossChainTx{
				InboundParams: &types.InboundParams{
					CoinType: coin.CoinType_Zeta,
				},
				CctxStatus: &types.Status{
					Status: types.CctxStatus_PendingOutbound,
				},
			},
			expected: OutboundTypeZetaWithdraw,
		},
		{
			name: "ZETA withdraw and call",
			cctx: types.CrossChainTx{
				InboundParams: &types.InboundParams{
					CoinType:         coin.CoinType_Zeta,
					IsCrossChainCall: true,
				},
				CctxStatus: &types.Status{
					Status: types.CctxStatus_PendingOutbound,
				},
			},
			expected: OutboundTypeZetaWithdrawAndCall,
		},
		{
			name: "Zeta withdraw revert and call on revert",
			cctx: types.CrossChainTx{
				InboundParams: &types.InboundParams{
					CoinType: coin.CoinType_Zeta,
				},
				CctxStatus: &types.Status{
					Status: types.CctxStatus_PendingRevert,
				},
				RevertOptions: types.RevertOptions{
					CallOnRevert: true,
				},
			},
			expected: OutboundTypeZetaWithdrawRevertAndCallOnRevert,
		},
		{
			name: " Zeta withdraw revert",
			cctx: types.CrossChainTx{
				InboundParams: &types.InboundParams{
					CoinType: coin.CoinType_Zeta,
				},
				CctxStatus: &types.Status{
					Status: types.CctxStatus_PendingRevert,
				},
			},
			expected: OutboundTypeZetaWithdrawRevert,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseOutboundTypeFromCCTX(tt.cctx)
			require.Equal(t, tt.expected, result)
		})
	}
}
